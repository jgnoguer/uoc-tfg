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

	"github.com/gocql/gocql"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/scylladb/gocqlx/v3"
	"github.com/scylladb/gocqlx/v3/qb"
	"github.com/scylladb/gocqlx/v3/table"
)

const ResourcePathRegex = "^(/groups)?/([a-zA-Z0-9-]{36})$"
const ResourcePathMembersRegex = "^(/groups)?/([a-zA-Z0-9-]{36})/members/([a-zA-Z0-9-]{36})$"

type ResponseMessage struct {
	Message string `json:"message"`
}

// metadata specifies table name and columns it must be in sync with schema.
var groupMetadata = table.Metadata{
	Name:    "uoc_animals.group",
	Columns: []string{"id", "name", "description", "members", "created_at", "updated_at"},
	PartKey: []string{"id"},
	SortKey: []string{"id"},
}

// Handle an HTTP Request.
func Handle(w http.ResponseWriter, r *http.Request) {

	godotenv.Load()
	dbIp := os.Getenv("SCYLLADB_IP")
	dbUser := os.Getenv("SCYLLA_APPUSER")
	dbPwd := os.Getenv("SCYLLA_APPPWD")
	appVersion := os.Getenv("GROUPS_VERSION")

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

	if strings.Contains(r.URL.Path, "/members") {
		switch r.Method {
		case "PUT":
			addGroupMember(w, r, session)
		case "DELETE":
			removeGroupMember(w, r, session)
		default:
			w.WriteHeader(http.StatusBadRequest)
		}
	} else {
		switch r.Method {
		case "GET":
			getGroup(w, r, session)
		case "POST":
			addGroup(w, r, session)
		case "PUT":
			modifyGroup(w, r, session)
		case "DELETE":
			deleteGroup(w, r, session)
		default:
			w.WriteHeader(http.StatusBadRequest)
		}
	}

}

func resolveGroupId(urlPath string, pattern string) string {
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(urlPath)
	if len(matches) <= 0 {
		return ""
	} else {
		return matches[2]
	}
}

func resolveMemberId(urlPath string, pattern string) string {
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(urlPath)
	if len(matches) <= 0 {
		return ""
	} else {
		return matches[3]
	}
}

func addGroup(w http.ResponseWriter, r *http.Request, session gocqlx.Session) {
	if (r.URL.Path != "/") && (r.URL.Path != "/groups") {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse := ResponseMessage{Message: "Bad url path for adding an group"}
		json.NewEncoder(w).Encode(errorResponse)
		return
	}
	var ag model.Group

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

	slog.Info("Group received", "data", ag)

	w.WriteHeader(http.StatusOK)
	groupId := uuid.New()
	ag.Id = groupId.String()
	ag.CreatedAt = time.Now()

	var groupTable = table.New(groupMetadata)
	q := session.Query(groupTable.Insert()).BindStruct(ag)
	if err := q.ExecRelease(); err != nil {
		panic(fmt.Errorf("error in exec query to insert group %w", err))
	}

	okResponse := ResponseMessage{Message: groupId.String() + " added"}
	json.NewEncoder(w).Encode(okResponse)

}

func modifyGroup(w http.ResponseWriter, r *http.Request, session gocqlx.Session) {

	var grp model.Group

	groupId := resolveGroupId(r.URL.Path, ResourcePathRegex)

	if groupId == "" {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse := ResponseMessage{Message: "Please specify a valid group id."}
		json.NewEncoder(w).Encode(errorResponse)
	} else {
		var groups = findGroupById(groupId, session)
		if len(groups) == 0 {
			w.WriteHeader(http.StatusNotFound)
			return
		}
	}

	err := decodeJSONBody(w, r, &grp)
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

	slog.Info("Group data received", "data", grp)

	grp.Id = groupId
	q := qb.Update("uoc_animals.group").
		Set("firstname", "lastname", "type", "status").
		Where(qb.Eq("id")).
		Query(session).
		BindStruct(grp)
	if err := q.ExecRelease(); err != nil {
		panic(fmt.Errorf("error in exec query to update group %w", err))
	}

	w.WriteHeader(http.StatusOK)
	okResponse := ResponseMessage{Message: groupId + " updated"}
	json.NewEncoder(w).Encode(okResponse)

}

func getGroup(w http.ResponseWriter, r *http.Request, session gocqlx.Session) {
	slog.Info("Finding group metadata", "id", r.URL.Path)

	groupId := resolveGroupId(r.URL.Path, ResourcePathRegex)

	if groupId == "" {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse := ResponseMessage{Message: "Please specify a valid group id."}
		json.NewEncoder(w).Encode(errorResponse)
	} else {
		var groups = findGroupById(groupId, session)
		if len(groups) == 0 {
			w.WriteHeader(http.StatusNotFound)
		} else {
			json.NewEncoder(w).Encode(groups[0])
		}
	}

}

func findGroupById(groupId string, session gocqlx.Session) []model.Group {
	var groups []model.Group
	var groupTable = table.New(groupMetadata)
	q := session.Query(groupTable.Select()).BindMap(qb.M{"id": groupId})
	if err := q.SelectRelease(&groups); err != nil {
		panic(fmt.Errorf("error in exec query to get group: %w", err))
	}
	return groups
}

func deleteGroup(w http.ResponseWriter, r *http.Request, session gocqlx.Session) {
	slog.Info("Deleting group", "id", r.URL.Path)
	groupId := resolveGroupId(r.URL.Path, ResourcePathRegex)

	if groupId == "" {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse := ResponseMessage{Message: "Bad url path for adding an group"}
		json.NewEncoder(w).Encode(errorResponse)
	} else {
		var groups = findGroupById(groupId, session)
		if len(groups) == 0 {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		q := session.Query(`DELETE FROM uoc_animals.group WHERE id = ?`,
			[]string{":id"}).
			BindMap(map[string]interface{}{
				":id": groupId,
			})

		error := q.ExecRelease()
		if error != nil {
			w.WriteHeader(http.StatusBadRequest)
			errorResponse := ResponseMessage{Message: fmt.Errorf("error to exec delete query: %w", error).Error()}
			json.NewEncoder(w).Encode(errorResponse)
			return
		}
		slog.Info("Db deleted group", "id", groupId)

		w.WriteHeader(http.StatusOK)
		errorResponse := ResponseMessage{Message: groupId + " deleted"}
		json.NewEncoder(w).Encode(errorResponse)
	}
}
func addGroupMember(w http.ResponseWriter, r *http.Request, session gocqlx.Session) {

	groupId := resolveGroupId(r.URL.Path, ResourcePathMembersRegex)
	memberId := resolveMemberId(r.URL.Path, ResourcePathMembersRegex)

	if (groupId == "") || (memberId == "") {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse := ResponseMessage{Message: "Bad url path for adding a member to a group"}
		json.NewEncoder(w).Encode(errorResponse)
		return
	}
	slog.Info("Adding member to group", "member", memberId, "group", groupId)

	w.WriteHeader(http.StatusOK)

	var grp model.Group
	grp.Id = groupId
	grp.UpdatedAt = time.Now()
	grp.Members = []string{memberId}
	q := qb.Update("uoc_animals.group").Add("members").Set("updated_at").Where(qb.Eq("id")).Query(session).BindStruct(grp)

	if err := q.ExecRelease(); err != nil {
		panic(fmt.Errorf("error in exec query to add a member to a group %w", err))
	}

	okResponse := ResponseMessage{Message: "member " + memberId + " added to group " + groupId}
	json.NewEncoder(w).Encode(okResponse)

}

func removeGroupMember(w http.ResponseWriter, r *http.Request, session gocqlx.Session) {
	groupId := resolveGroupId(r.URL.Path, ResourcePathMembersRegex)
	memberId := resolveMemberId(r.URL.Path, ResourcePathMembersRegex)

	if (groupId == "") || (memberId == "") {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse := ResponseMessage{Message: "Bad url path for adding a member to a group"}
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	slog.Info("Remove member from group", "member", memberId, "group", groupId)

	w.WriteHeader(http.StatusOK)

	var grp model.Group
	grp.Id = groupId
	grp.UpdatedAt = time.Now()
	grp.Members = []string{memberId}
	q := qb.Update("uoc_animals.group").Remove("members").Set("updated_at").Where(qb.Eq("id")).Query(session).BindStruct(grp)

	if err := q.ExecRelease(); err != nil {
		panic(fmt.Errorf("error in exec query to remove a member from a group %w", err))
	}

	okResponse := ResponseMessage{Message: "member " + memberId + " removed from group " + groupId}
	json.NewEncoder(w).Encode(okResponse)
}
