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

const ResourcePathRegex = "^/([a-zA-Z0-9-]{36})$"

type ResponseMessage struct {
	Message string `json:"message"`
}

// metadata specifies table name and columns it must be in sync with schema.
var agentMetadata = table.Metadata{
	Name:    "uoc_animals.agent",
	Columns: []string{"id", "firstname", "lastname", "type", "status", "created_at"},
	PartKey: []string{"id"},
	SortKey: []string{"id"},
}

// Handle an HTTP Request.
func Handle(w http.ResponseWriter, r *http.Request) {

	godotenv.Load()
	dbIp := os.Getenv("SCYLLADB_IP")
	dbUser := os.Getenv("SCYLLA_APPUSER")
	dbPwd := os.Getenv("SCYLLA_APPPWD")
	appVersion := os.Getenv("AGENTS_VERSION")

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
		getAgent(w, r, session)
	case "POST":
		addAgent(w, r, session)
	case "PUT":
		modifyAgent(w, r, session)
	case "DELETE":
		deleteAgent(w, r, session)
	default:
		w.WriteHeader(http.StatusBadRequest)
	}

}

func addAgent(w http.ResponseWriter, r *http.Request, session gocqlx.Session) {
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse := ResponseMessage{Message: "Bad url path for adding an agent"}
		json.NewEncoder(w).Encode(errorResponse)
		return
	}
	var ag model.Agent

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

	slog.Info("Agent received", "data", ag)

	w.WriteHeader(http.StatusOK)
	agentId := uuid.New()
	ag.Id = agentId.String()
	ag.CreatedAt = time.Now()

	var agentTable = table.New(agentMetadata)
	q := session.Query(agentTable.Insert()).BindStruct(ag)
	if err := q.ExecRelease(); err != nil {
		panic(fmt.Errorf("error in exec query to insert agent %w", err))
	}

	okResponse := ResponseMessage{Message: agentId.String() + " added"}
	json.NewEncoder(w).Encode(okResponse)

}

func modifyAgent(w http.ResponseWriter, r *http.Request, session gocqlx.Session) {

	var ag model.Agent

	agentId := resolveAgentId(r.URL.Path, ResourcePathRegex)

	if agentId == "" {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse := ResponseMessage{Message: "Please specify a valid agent id."}
		json.NewEncoder(w).Encode(errorResponse)
	} else {
		var agents = findAgentById(agentId, session)
		if len(agents) == 0 {
			w.WriteHeader(http.StatusNotFound)
			return
		}
	}

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

	slog.Info("Agent data received", "data", ag)

	ag.Id = agentId
	q := qb.Update("uoc_animals.agent").
		Set("firstname", "lastname", "type", "status").
		Where(qb.Eq("id")).
		Query(session).
		BindStruct(ag)
	if err := q.ExecRelease(); err != nil {
		panic(fmt.Errorf("error in exec query to update agent %w", err))
	}

	w.WriteHeader(http.StatusOK)
	okResponse := ResponseMessage{Message: agentId + " updated"}
	json.NewEncoder(w).Encode(okResponse)

}

func getAgent(w http.ResponseWriter, r *http.Request, session gocqlx.Session) {
	slog.Info("Finding agent metadata", "id", r.URL.Path)

	agentId := resolveAgentId(r.URL.Path, ResourcePathRegex)

	if agentId == "" {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse := ResponseMessage{Message: "Please specify a valid agent id."}
		json.NewEncoder(w).Encode(errorResponse)
	} else {
		var agents = findAgentById(agentId, session)
		if len(agents) == 0 {
			w.WriteHeader(http.StatusNotFound)
		} else {
			json.NewEncoder(w).Encode(agents[0])
		}
	}

}

func findAgentById(agentId string, session gocqlx.Session) []model.Agent {
	var agents []model.Agent
	var agentTable = table.New(agentMetadata)
	q := session.Query(agentTable.Select()).BindMap(qb.M{"id": agentId})
	if err := q.SelectRelease(&agents); err != nil {
		panic(fmt.Errorf("error in exec query to get agent: %w", err))
	}
	return agents
}

func deleteAgent(w http.ResponseWriter, r *http.Request, session gocqlx.Session) {
	slog.Info("Deleting agent", "id", r.URL.Path)
	agentId := resolveAgentId(r.URL.Path, ResourcePathRegex)

	if agentId == "" {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse := ResponseMessage{Message: "Bad url path for adding an agent"}
		json.NewEncoder(w).Encode(errorResponse)
	} else {
		var agents = findAgentById(agentId, session)
		if len(agents) == 0 {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		q := session.Query(`DELETE FROM uoc_animals.agent WHERE id = ?`,
			[]string{":id"}).
			BindMap(map[string]interface{}{
				":id": agentId,
			})

		error := q.ExecRelease()
		if error != nil {
			w.WriteHeader(http.StatusBadRequest)
			errorResponse := ResponseMessage{Message: fmt.Errorf("error to exec delete query: %w", error).Error()}
			json.NewEncoder(w).Encode(errorResponse)
			return
		}
		slog.Info("Db deleted agent", "id", agentId)

		w.WriteHeader(http.StatusOK)
		errorResponse := ResponseMessage{Message: agentId + " deleted"}
		json.NewEncoder(w).Encode(errorResponse)
	}
}

func resolveAgentId(urlPath string, pattern string) string {
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(urlPath)
	if len(matches) <= 0 {
		return ""
	} else {
		return matches[1]
	}
}
