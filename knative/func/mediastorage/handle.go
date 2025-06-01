package function

import (
	"encoding/json"
	"fmt"
	"function/model"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/scylladb/gocqlx/v3/qb"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/scylladb/gocqlx/v3/table"

	"github.com/gocql/gocql"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/scylladb/gocqlx/v3"
)

const ResourcePathRegex = "^(/mediastorage)?/([a-zA-Z0-9-]{36})$"
const ResourceMetadataPathRegex = "^(/mediastorage)?/([a-zA-Z0-9-]{36})/metadata$"

type ResponseMessage struct {
	Message string `json:"message"`
}

// metadata specifies table name and columns it must be in sync with schema.
var mediaMetadata = table.Metadata{
	Name:    "uoc_animals.media",
	Columns: []string{"id", "name", "contenttype", "location", "created_at", "status", "size"},
	PartKey: []string{"id"},
	SortKey: []string{"id"},
}

// Handle an HTTP Request.
func Handle(w http.ResponseWriter, r *http.Request) {

	godotenv.Load()
	storageFolder := os.Getenv("STORAGE_FOLDER")
	dbIp := os.Getenv("SCYLLADB_IP")
	dbUser := os.Getenv("SCYLLA_APPUSER")
	dbPwd := os.Getenv("SCYLLA_APPPWD")
	appVersion := os.Getenv("MEDIASTORE_VERSION")

	slog.Info("Info", "v.", appVersion, "Method", r.Method, "Storage folder", storageFolder)

	slog.Info("Try connection to ", "ip", dbIp)
	cluster := gocql.NewCluster(dbIp)
	//
	cluster.Authenticator = gocql.PasswordAuthenticator{Username: dbUser, Password: dbPwd}
	//
	session, err := gocqlx.WrapSession(cluster.CreateSession())

	if err != nil {
		slog.Error("Fail to connect to database" + err.Error())
		panic("Database connection failed ")
	}

	switch r.Method {
	case "GET":
		if strings.HasSuffix(r.URL.Path, "/metadata") {
			getMediaMetadata(w, r, session)
		} else {
			getMedia(w, r, session)
		}
	case "POST":
		w.Header().Set("Content-Type", "application/json")
		addMedia(w, r, session)
	case "DELETE":
		w.Header().Set("Content-Type", "application/json")
		deleteMedia(w, r, session)
	default:
		w.WriteHeader(http.StatusBadRequest)
	}

}

func addMedia(w http.ResponseWriter, r *http.Request, session gocqlx.Session) {
	f, fileHandler, err := r.FormFile("file")
	if err != nil {
		slog.Error("Fail to get file from form" + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer f.Close()

	if (r.URL.Path != "/") && (r.URL.Path != "/mediastorage") {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse := ResponseMessage{Message: "Bad url path for adding a file"}
		json.NewEncoder(w).Encode(errorResponse)
		return
	}
	id := uuid.New()

	slog.Info("file", "name", fileHandler.Filename, "size", fileHandler.Size, "mime", fileHandler.Header.Get("Content-Type"))
	newDirPath := filepath.Join(os.Getenv("STORAGE_FOLDER"), id.String()[:3], id.String())
	dbLocation := id.String()[:3] + "/" + id.String()
	dirError := os.MkdirAll(newDirPath, os.ModePerm)
	if dirError != nil {
		slog.Error("Fail to create directory" + dirError.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	newFilePath := filepath.Join(newDirPath, fileHandler.Filename)
	newFile, err := os.Create(newFilePath)
	if err != nil {
		slog.Error("Fail to create file" + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer newFile.Close()

	_, err = io.Copy(newFile, f)

	slog.Info(newFilePath)
	if err != nil {
		slog.Error("Fail to save the file" + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	p := model.Media{Id: id.String(), Name: fileHandler.Filename,
		ContentType: fileHandler.Header.Get("Content-Type"),
		Location:    dbLocation, CreatedAt: time.Now(), Status: 0, Size: fileHandler.Size}

	json.NewEncoder(w).Encode(p)

	var mediaTable = table.New(mediaMetadata)
	q := session.Query(mediaTable.Insert()).BindStruct(p)
	if err := q.ExecRelease(); err != nil {
		panic(fmt.Errorf("error in exec query to insert media %w", err))
	}

	event, err := publishEvent("123456")
	slog.Info("Publishing event", "event", event)
}

func publishEvent(mediaId string) (*cloudevents.Event, cloudevents.Result) {
	newEvent := cloudevents.NewEvent()
	// Setting the ID here is not necessary. When using NewDefaultClient the ID is set
	// automatically. We set the ID anyway so it appears in the log.
	newEvent.SetID(uuid.New().String())
	newEvent.SetSource("dev.jgnoguer.knative.uoc/mediastorage-service")
	newEvent.SetType("dev.jgnoguer.knative.uoc.imageadded")
	slog.Info("Responding with event\n%s\n", newEvent)
	if err := newEvent.SetData(cloudevents.ApplicationJSON, ImageAdded{MediaId: mediaId}); err != nil {
		return nil, cloudevents.NewHTTPResult(500, "failed to set response data: %s", err)
	}
	return &newEvent, cloudevents.ResultACK
}

func getMediaMetadata(w http.ResponseWriter, r *http.Request, session gocqlx.Session) {
	slog.Info("Finding media metadata", "id", r.URL.Path)

	mediaId := resolveMediaId(r.URL.Path, ResourceMetadataPathRegex)

	if mediaId == "" {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse := ResponseMessage{Message: "Please specify a valid media id."}
		json.NewEncoder(w).Encode(errorResponse)
	} else {
		var medias = findMediaById(mediaId, session)
		if len(medias) == 0 {
			w.WriteHeader(http.StatusNotFound)
		} else {
			json.NewEncoder(w).Encode(medias[0])
		}
	}

}

func getMedia(w http.ResponseWriter, r *http.Request, session gocqlx.Session) {
	slog.Info("Finding media", "id", r.URL.Path)

	mediaId := resolveMediaId(r.URL.Path, ResourcePathRegex)

	if mediaId == "" {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse := ResponseMessage{Message: "Please specify a valid media id."}
		json.NewEncoder(w).Encode(errorResponse)
	} else {
		var medias = findMediaById(mediaId, session)
		if len(medias) == 0 {
			w.WriteHeader(http.StatusNotFound)
		} else {
			foundMedia := medias[0]
			var location = resolveStorageFile(foundMedia)
			f, err := os.Open(location)
			if err != nil {
				panic(fmt.Errorf("error while reading file %w", err))
			}
			w.Header().Set("Content-Type", foundMedia.ContentType)
			http.ServeContent(w, r, foundMedia.Name, time.Now(), f)
		}
	}

}

func findMediaById(mediaId string, session gocqlx.Session) []model.Media {
	var medias []model.Media
	var mediaTable = table.New(mediaMetadata)
	q := session.Query(mediaTable.Select()).BindMap(qb.M{"id": mediaId})
	if err := q.SelectRelease(&medias); err != nil {
		panic(fmt.Errorf("error in exec query to get media: %w", err))
	}
	return medias
}

func deleteMedia(w http.ResponseWriter, r *http.Request, session gocqlx.Session) {
	slog.Info("Deleting media", "id", r.URL.Path)
	mediaId := resolveMediaId(r.URL.Path, ResourcePathRegex)

	if mediaId == "" {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse := ResponseMessage{Message: "Bad url path for adding a file"}
		json.NewEncoder(w).Encode(errorResponse)
	} else {
		var medias = findMediaById(mediaId, session)
		if len(medias) == 0 {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		q := session.Query(`DELETE FROM uoc_animals.media WHERE id = ?`,
			[]string{":id"}).
			BindMap(map[string]interface{}{
				":id": mediaId,
			})

		error := q.ExecRelease()
		if error != nil {
			w.WriteHeader(http.StatusBadRequest)
			errorResponse := ResponseMessage{Message: fmt.Errorf("error to exec delete query: %w", error).Error()}
			json.NewEncoder(w).Encode(errorResponse)
			return
		}
		slog.Info("Db deleted media", "id", mediaId)

		folder := resolveStorageFolder(medias[0])
		os.RemoveAll(folder)

		slog.Info("Storage deleted media ", "folder", folder)

		w.WriteHeader(http.StatusOK)
		errorResponse := ResponseMessage{Message: mediaId + " deleted"}
		json.NewEncoder(w).Encode(errorResponse)
	}
}

func resolveMediaId(urlPath string, pattern string) string {
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(urlPath)
	if len(matches) <= 0 {
		return ""
	} else {
		return matches[2]
	}
}

func resolveStorageFolder(media model.Media) string {
	storageFolder := os.Getenv("STORAGE_FOLDER")
	return filepath.Join(storageFolder, media.Location)
}

func resolveStorageFile(media model.Media) string {
	storageFolder := os.Getenv("STORAGE_FOLDER")
	return filepath.Join(storageFolder, media.Location, media.Name)
}
