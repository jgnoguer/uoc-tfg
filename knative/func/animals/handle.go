package function

import (
	"encoding/json"
	"errors"
	"fmt"
	"function/model"
	"log"
	"log/slog"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/scylladb/gocqlx/v3/qb"

	"github.com/scylladb/gocqlx/v3/table"

	"github.com/gocql/gocql"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/scylladb/gocqlx/v3"
)

const ResourcePathRegex = "^(/animals)?/([a-zA-Z0-9-]{36})$"
const ResourcePathMediaRegex = "^(/animals)?/([a-zA-Z0-9-]{36})/media/([a-zA-Z0-9-]{36})$"
const ResourcePathMediaQueryRegex = "^(/animals)?/([a-zA-Z0-9-]{36})/media"

type ResponseMessage struct {
	Message string `json:"message"`
}

// metadata specifies table name and columns it must be in sync with schema.
var animalMetadata = table.Metadata{
	Name:    "uoc_animals.animal",
	Columns: []string{"id", "name", "description", "breed", "type", "status", "created_at", "updated_at"},
	PartKey: []string{"id"},
	SortKey: []string{"id"},
}

// metadata specifies table name and columns it must be in sync with schema.
var animalMediaMetadata = table.Metadata{
	Name:    "uoc_animals.animal_media",
	Columns: []string{"id", "media_id", "description", "type", "created_at"},
	PartKey: []string{"id", "media_id"},
	SortKey: []string{"id", "media_id"},
}

func resolveAnimalId(urlPath string, pattern string) string {
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(urlPath)
	if len(matches) <= 0 {
		return ""
	} else {
		return matches[2]
	}
}

func resolveMediaId(urlPath string, pattern string) string {
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(urlPath)
	if len(matches) <= 0 {
		return ""
	} else {
		return matches[3]
	}
}

// Handle an HTTP Request.
func Handle(w http.ResponseWriter, r *http.Request) {

	godotenv.Load()
	dbIp := os.Getenv("SCYLLADB_IP")
	dbUser := os.Getenv("SCYLLA_APPUSER")
	dbPwd := os.Getenv("SCYLLA_APPPWD")
	appVersion := os.Getenv("ANIMALS_VERSION")

	slog.Info("Info", "v.", appVersion, "Method", r.Method)

	slog.Info("🐕 Try connection to ", "ip", dbIp)
	cluster := gocql.NewCluster(dbIp)
	//
	cluster.Authenticator = gocql.PasswordAuthenticator{Username: dbUser, Password: dbPwd}
	//
	session, err := gocqlx.WrapSession(cluster.CreateSession())

	if err != nil {
		panic("Database connection failed ")
	}

	w.Header().Set("Content-Type", "application/json")

	if strings.Contains(r.URL.Path, "/media") {
		switch r.Method {
		case "GET":
			getAnimalMedias(w, r, session)
		case "POST":
			addAnimalMedia(w, r, session)
		case "DELETE":
			deleteAnimalMedia(w, r, session)
		default:
			w.WriteHeader(http.StatusBadRequest)
		}
	} else {
		switch r.Method {
		case "GET":
			getAnimal(w, r, session)
		case "POST":
			addAnimal(w, r, session)
		case "PUT":
			modifyAnimal(w, r, session)
		case "DELETE":
			deleteAnimal(w, r, session)
		default:
			w.WriteHeader(http.StatusBadRequest)
		}
	}

}

func addAnimal(w http.ResponseWriter, r *http.Request, session gocqlx.Session) {
	if (r.URL.Path != "/") && (r.URL.Path != "/animals") {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse := ResponseMessage{Message: "Bad url path for adding an animal"}
		json.NewEncoder(w).Encode(errorResponse)
		return
	}
	var ag model.Animal

	err := decodeJSONBody(w, r, &ag)
	if err != nil {
		var mr *malformedRequest
		if errors.As(err, &mr) {
			http.Error(w, mr.msg, mr.status)
		} else {
			log.Print(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	slog.Info("Animal received", "data", ag)

	w.WriteHeader(http.StatusOK)
	animalId := uuid.New()
	ag.Id = animalId.String()
	ag.CreatedAt = time.Now()

	var animalTable = table.New(animalMetadata)
	q := session.Query(animalTable.Insert()).BindStruct(ag)
	if err := q.ExecRelease(); err != nil {
		panic(fmt.Errorf("error in exec query to insert animal %w", err))
	}

	okResponse := ResponseMessage{Message: animalId.String() + " added"}
	json.NewEncoder(w).Encode(okResponse)

}

func modifyAnimal(w http.ResponseWriter, r *http.Request, session gocqlx.Session) {
	var animal model.Animal
	animalId := resolveAnimalId(r.URL.Path, ResourcePathRegex)

	if animalId == "" {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse := ResponseMessage{Message: "Please specify a valid animal id."}
		json.NewEncoder(w).Encode(errorResponse)
	} else {
		var animals = findAnimalById(animalId, session)
		if len(animals) == 0 {
			w.WriteHeader(http.StatusNotFound)
			return
		}
	}

	err := decodeJSONBody(w, r, &animal)
	if err != nil {
		var mr *malformedRequest
		if errors.As(err, &mr) {
			http.Error(w, mr.msg, mr.status)
		} else {
			log.Print(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	slog.Info("Animal data received", "data", animal)

	animal.Id = animalId
	animal.CreatedAt = time.Now()

	q := qb.Update("uoc_animals.animal").
		Set("name", "description", "type", "breed", "status", "updated_at").
		Where(qb.Eq("id")).
		Query(session).
		BindStruct(animal)
	if err := q.ExecRelease(); err != nil {
		panic(fmt.Errorf("error in exec query to update animal %w", err))
	}

	w.WriteHeader(http.StatusOK)
	okResponse := ResponseMessage{Message: animalId + " updated"}
	json.NewEncoder(w).Encode(okResponse)

}

func getAnimal(w http.ResponseWriter, r *http.Request, session gocqlx.Session) {
	slog.Info("Finding animal metadata", "id", r.URL.Path)

	animalId := resolveAnimalId(r.URL.Path, ResourcePathRegex)

	if animalId == "" {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse := ResponseMessage{Message: "Please specify a valid animal id."}
		json.NewEncoder(w).Encode(errorResponse)
	} else {
		var animals = findAnimalById(animalId, session)
		if len(animals) == 0 {
			w.WriteHeader(http.StatusNotFound)
		} else {
			json.NewEncoder(w).Encode(animals[0])
		}
	}
}

func findAnimalById(animalId string, session gocqlx.Session) []model.Animal {
	var animals []model.Animal
	var animalTable = table.New(animalMetadata)
	q := session.Query(animalTable.Select()).BindMap(qb.M{"id": animalId})
	if err := q.SelectRelease(&animals); err != nil {
		panic(fmt.Errorf("error in exec query to get animal: %w", err))
	}
	return animals
}

func deleteAnimal(w http.ResponseWriter, r *http.Request, session gocqlx.Session) {
	slog.Info("Deleting animal", "id", r.URL.Path)
	animalId := resolveAnimalId(r.URL.Path, ResourcePathRegex)

	if animalId == "" {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse := ResponseMessage{Message: "Bad url path for adding an animal"}
		json.NewEncoder(w).Encode(errorResponse)
	} else {
		var animals = findAnimalById(animalId, session)
		if len(animals) == 0 {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		q := session.Query(`DELETE FROM uoc_animals.animal WHERE id = ?`,
			[]string{":id"}).
			BindMap(map[string]interface{}{
				":id": animalId,
			})

		error := q.ExecRelease()
		if error != nil {
			w.WriteHeader(http.StatusBadRequest)
			errorResponse := ResponseMessage{Message: fmt.Errorf("error to exec delete query: %w", error).Error()}
			json.NewEncoder(w).Encode(errorResponse)
			return
		}
		slog.Info("Db deleted animal", "id", animalId)

		w.WriteHeader(http.StatusOK)
		errorResponse := ResponseMessage{Message: animalId + " deleted"}
		json.NewEncoder(w).Encode(errorResponse)
	}
}

func addAnimalMedia(w http.ResponseWriter, r *http.Request, session gocqlx.Session) {

	animalId := resolveAnimalId(r.URL.Path, ResourcePathMediaRegex)
	mediaId := resolveMediaId(r.URL.Path, ResourcePathMediaRegex)

	if (animalId == "") || (mediaId == "") {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse := ResponseMessage{Message: "Bad url path for adding a media for an animal"}
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	var animalMedia model.AnimalMedia

	err := decodeJSONBody(w, r, &animalMedia)
	if err != nil {
		var mr *malformedRequest
		if errors.As(err, &mr) {
			http.Error(w, mr.msg, mr.status)
		} else {
			log.Print(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	slog.Info("Animal media received", "data", animalMedia)

	w.WriteHeader(http.StatusOK)

	animalMedia.CreatedAt = time.Now()
	animalMedia.MediaId = mediaId
	animalMedia.Id = animalId

	var animalTable = table.New(animalMediaMetadata)
	q := session.Query(animalTable.Insert()).BindStruct(animalMedia)
	if err := q.ExecRelease(); err != nil {
		panic(fmt.Errorf("error in exec query to insert animal media %w", err))
	}

	okResponse := ResponseMessage{Message: "media id " + animalMedia.MediaId + " for animal " + animalMedia.Id + " added"}
	json.NewEncoder(w).Encode(okResponse)

}

func getAnimalMedias(w http.ResponseWriter, r *http.Request, session gocqlx.Session) {
	slog.Info("Finding animal medias", "id", r.URL.Path)

	animalId := resolveAnimalId(r.URL.Path, ResourcePathMediaQueryRegex)

	if animalId == "" {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse := ResponseMessage{Message: "Please specify a valid animal id."}
		json.NewEncoder(w).Encode(errorResponse)
	} else {
		var animalMedias = findAnimalMediasById(animalId, session)
		if len(animalMedias) == 0 {
			w.WriteHeader(http.StatusNotFound)
		} else {
			json.NewEncoder(w).Encode(animalMedias)
		}
	}
}

func findAnimalMediasById(animalId string, session gocqlx.Session) []model.AnimalMedia {
	var medias []model.AnimalMedia

	q := qb.Select("uoc_animals.animal_media").Where(qb.EqLit("id", animalId)).Query(session)

	if err := q.SelectRelease(&medias); err != nil {
		panic(fmt.Errorf("error in exec query to select animal medias %w", err))
	}

	return medias
}

func deleteAnimalMedia(w http.ResponseWriter, r *http.Request, session gocqlx.Session) {
	slog.Info("Deleting animal", "id", r.URL.Path)
	animalId := resolveAnimalId(r.URL.Path, ResourcePathMediaRegex)
	mediaId := resolveMediaId(r.URL.Path, ResourcePathMediaRegex)

	if animalId == "" {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse := ResponseMessage{Message: "Bad url path for adding an animal"}
		json.NewEncoder(w).Encode(errorResponse)
	} else {
		q := session.Query(`DELETE FROM uoc_animals.animal_media WHERE id = ? and media_id = ?`,
			[]string{":id", ":media_id"}).
			BindMap(map[string]interface{}{
				":id":       animalId,
				":media_id": mediaId,
			})

		error := q.ExecRelease()
		if error != nil {
			w.WriteHeader(http.StatusBadRequest)
			errorResponse := ResponseMessage{Message: fmt.Errorf("error to exec delete query: %w", error).Error()}
			json.NewEncoder(w).Encode(errorResponse)
			return
		}
		slog.Info("Db deleted media ", "animalId", animalId, "media_id", mediaId)

		w.WriteHeader(http.StatusOK)
		errorResponse := ResponseMessage{Message: animalId + " deleted"}
		json.NewEncoder(w).Encode(errorResponse)
	}
}
