package function

import (
	"fmt"
	"time"
	"net/http"
	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v3"
)

func (s Song) String() string {
	return fmt.Sprintf("Id: %s\nTitle: %s\nArtist: %s\nAlbum: %s\nCreated At: %s\n", s.Id, s.Title, s.Artist, s.Album, s.Created_at)
}

type Song struct {
	Id string
	Title string
	Artist string
	Album string
	Created_at time.Time
}

// Handle an HTTP Request.
func Handle(w http.ResponseWriter, r *http.Request) {
	/*
	 * YOUR CODE HERE
	 *
	 * Try running `go test`.  Add more test as you code in `handle_test.go`.
	 */
	 cluster := gocql.NewCluster("192.168.2.126")

	 cluster.Authenticator = gocql.PasswordAuthenticator{Username: "cassandra", Password: "cassandra"}
 
	session, err := gocqlx.WrapSession(cluster.CreateSession())

	if err != nil {
		panic("Connection fail")
	}

	var songs []Song

	q := session.Query("SELECT * FROM media_player.playlist", nil)

	if err := q.Select(&songs); err != nil {
		panic(fmt.Errorf("error in exec query to list playlists: %w", err))
	}

	fmt.Print(songs)
	
	// song := Song{}
	
	// q := session.Query(
	// 	`INSERT INTO media_player.playlist (id,title,artist,album,created_at) VALUES (now(),?,?,?,?)`,
	// 	[]string{":title", ":artist", ":album", ":created_at"}).
	// 	BindMap(map[string]interface{} {
	// 		":title":      song.Title,
	// 		":artist":     song.Artist,
	// 		":album":      song.Album,
	// 		":created_at": time.Now(),
	// 	})
	
	// err2 := q.Exec();
	// if err2 != nil {
	// 	panic(fmt.Errorf("error in exec query to insert a song in playlist %w", err2))
	// }

	fmt.Println("HOLA - Received request", session)
	fmt.Fprintf(w, "%q", "Boooo")
}


