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
	"time"

	"github.com/scylladb/gocqlx/v3/qb"

	"github.com/scylladb/gocqlx/v3/table"

	"github.com/gocql/gocql"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/scylladb/gocqlx/v3"
)

const ResourcePathRegex = "^(/activities)?/([a-zA-Z0-9-]{36})$"
const ResourcePathUpdateRegex = "^(/activities)?/([a-zA-Z0-9-]{36})/status$"

type ResponseMessage struct {
	Message string `json:"message"`
}

// metadata specifies table name and columns it must be in sync with schema.
var activityMetadata = table.Metadata{
	Name:    "uoc_animals.activity",
	Columns: []string{"id", "shortcode", "description", "type", "status", "created_at", "updated_at"},
	PartKey: []string{"id"},
	SortKey: []string{"id"},
}

var activityLogMetadata = table.Metadata{
	Name:    "uoc_animals.activity_log",
	Columns: []string{"id", "status", "update_time", "updated_by", "description"},
	PartKey: []string{"id"},
	SortKey: []string{"id"},
}

// Handle an HTTP Request.
func Handle(w http.ResponseWriter, r *http.Request) {

	godotenv.Load()
	dbIp := os.Getenv("SCYLLADB_IP")
	dbUser := os.Getenv("SCYLLA_APPUSER")
	dbPwd := os.Getenv("SCYLLA_APPPWD")
	appVersion := os.Getenv("ACTIVITIES_VERSION")

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

	switch r.Method {
	case "GET":
		getActivity(w, r, session)
	case "POST":
		addActivity(w, r, session)
	case "PUT":
		modifyActivityStatus(w, r, session)
	case "DELETE":
		deleteActivity(w, r, session)
	default:
		w.WriteHeader(http.StatusBadRequest)
	}

}

func addActivity(w http.ResponseWriter, r *http.Request, session gocqlx.Session) {
	if (r.URL.Path != "/") && (r.URL.Path != "/activities") {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse := ResponseMessage{Message: "Bad url path for adding an activity"}
		json.NewEncoder(w).Encode(errorResponse)
		return
	}
	var ag model.Activity

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

	slog.Info("Activity received", "data", ag)

	w.WriteHeader(http.StatusOK)
	activityId := uuid.New()
	ag.Id = activityId.String()
	ag.CreatedAt = time.Now()

	var activityTable = table.New(activityMetadata)
	q := session.Query(activityTable.Insert()).BindStruct(ag)
	if err := q.ExecRelease(); err != nil {
		panic(fmt.Errorf("error in exec query to insert activity %w", err))
	}

	okResponse := ResponseMessage{Message: activityId.String() + " added"}
	json.NewEncoder(w).Encode(okResponse)

}

func modifyActivityStatus(w http.ResponseWriter, r *http.Request, session gocqlx.Session) {

	var actStatusUpdate model.ActivityStatusUpdate

	activityId := resolveActivityId(r.URL.Path, ResourcePathUpdateRegex)

	if activityId == "" {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse := ResponseMessage{Message: "Invalid service call."}
		json.NewEncoder(w).Encode(errorResponse)
	} else {
		var activities = findActivityById(activityId, session)
		if len(activities) == 0 {
			w.WriteHeader(http.StatusNotFound)
			return
		}
	}

	err := decodeJSONBody(w, r, &actStatusUpdate)
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

	slog.Info("Activity status update data received", "data", actStatusUpdate)

	actStatusUpdate.Id = activityId
	actStatusUpdate.UpdatedTime = time.Now()
	q := qb.Update("uoc_animals.activity").
		Set("status").
		Set("updated_at").
		Where(qb.Eq("id")).
		Query(session).
		BindStruct(actStatusUpdate)
	if err := q.ExecRelease(); err != nil {
		panic(fmt.Errorf("error in exec query to update actStatusUpdate %w", err))
	}

	var activityLogMetadata = table.New(activityLogMetadata)
	var actLog = model.ActivityLog{Id: activityId, Status: actStatusUpdate.Status, UpdatedAt: time.Now(),
		UpdateBy: actStatusUpdate.Issuer, Description: actStatusUpdate.Description}
	qLog := session.Query(activityLogMetadata.Insert()).BindStruct(actLog)
	if err := qLog.ExecRelease(); err != nil {
		panic(fmt.Errorf("error in exec query to insert activity log %w", err))
	}

	w.WriteHeader(http.StatusOK)
	okResponse := ResponseMessage{Message: activityId + " updated"}
	json.NewEncoder(w).Encode(okResponse)

}

func getActivity(w http.ResponseWriter, r *http.Request, session gocqlx.Session) {
	slog.Info("Finding activity metadata", "url", r.URL.Path)

	shortCode := r.URL.Query().Get("shortcode")
	if shortCode != "" {
		slog.Info("Short code", "shortCode", shortCode)
		getActivityByCode(shortCode, w, r, session)
	} else {
		getActivityById(w, r, session)
	}

}

func getActivityById(w http.ResponseWriter, r *http.Request, session gocqlx.Session) {
	slog.Info("Finding activity metadata", "url", r.URL.Path)

	activityId := resolveActivityId(r.URL.Path, ResourcePathRegex)

	if activityId == "" {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse := ResponseMessage{Message: "Please specify a valid activity id."}
		json.NewEncoder(w).Encode(errorResponse)
	} else {
		var activities = findActivityById(activityId, session)
		if len(activities) == 0 {
			w.WriteHeader(http.StatusNotFound)
		} else {
			json.NewEncoder(w).Encode(activities[0])
		}
	}

}

func getActivityByCode(activitycode string, w http.ResponseWriter, r *http.Request, session gocqlx.Session) {
	slog.Info("Finding activity by short code", "shortcode", activitycode)

	var activities = findActivityByCode(activitycode, session)
	if len(activities) == 0 {
		w.WriteHeader(http.StatusNotFound)
	} else {
		json.NewEncoder(w).Encode(activities[0])
	}

}

func findActivityById(activityId string, session gocqlx.Session) []model.Activity {
	var activities []model.Activity
	var activityTable = table.New(activityMetadata)
	q := session.Query(activityTable.Select()).BindMap(qb.M{"id": activityId})
	if err := q.SelectRelease(&activities); err != nil {
		panic(fmt.Errorf("error in exec query to get activity by id: %w", err))
	}
	return activities
}

func findActivityByCode(activityCode string, session gocqlx.Session) []model.Activity {
	var activities []model.Activity

	row := &model.Activity{}
	q := qb.Select("uoc_animals.activity").Where(qb.EqLit("shortcode", "a34fsdr")).Query(session)

	if err := q.Iter().Get(row); err != nil {
		panic(fmt.Errorf("error in exec query to get activity by code: %w", err))
	} else {
		activities = append(activities, *row)
	}

	return activities
}

func deleteActivity(w http.ResponseWriter, r *http.Request, session gocqlx.Session) {
	slog.Info("Deleting activity", "id", r.URL.Path)
	activityId := resolveActivityId(r.URL.Path, ResourcePathRegex)

	if activityId == "" {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse := ResponseMessage{Message: "Bad url path for adding an activity"}
		json.NewEncoder(w).Encode(errorResponse)
	} else {
		var activities = findActivityById(activityId, session)
		if len(activities) == 0 {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		q := session.Query(`DELETE FROM uoc_animals.activity WHERE id = ?`,
			[]string{":id"}).
			BindMap(map[string]interface{}{
				":id": activityId,
			})

		error := q.ExecRelease()
		if error != nil {
			w.WriteHeader(http.StatusBadRequest)
			errorResponse := ResponseMessage{Message: fmt.Errorf("error to exec delete query: %w", error).Error()}
			json.NewEncoder(w).Encode(errorResponse)
			return
		}
		slog.Info("Db deleted activity", "id", activityId)

		w.WriteHeader(http.StatusOK)
		errorResponse := ResponseMessage{Message: activityId + " deleted"}
		json.NewEncoder(w).Encode(errorResponse)
	}
}

func resolveActivityId(urlPath string, pattern string) string {
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(urlPath)
	if len(matches) <= 0 {
		return ""
	} else {
		return matches[2]
	}
}
