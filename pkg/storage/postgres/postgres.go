package postgres

import (
	"database/sql"

	"fmt"

	_ "github.com/lib/pq"
	"github.com/ysinjab/spotigo/pkg/albums"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "ysinjab"
	password = "123"
	dbname   = "spotigo"
)

type storage struct {
	db *sql.DB
}

func NewStorage() (*storage, error) {
	connectionString := fmt.Sprintf("user=%s password=%s dbname=%s port=%d sslmode=disable", user, password, dbname, port)
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	// Create table if not exists
	strQuery := "CREATE TABLE IF NOT EXISTS albums (id serial NOT NULL, name VARCHAR not NULL);"

	_, err = db.Exec(strQuery)
	if err != nil {
		return nil, err
	}
	return &storage{db: db}, nil
}

func (s *storage) GetAlbums() ([]albums.Album, error) {
	list := []albums.Album{}
	rows, err := s.db.Query("SELECT * FROM albums;")
	if err != nil {
		return nil, err
	}
	for rows.Next() {

		var id int32
		var name string
		err = rows.Scan(&id, &name)
		if err != nil {
			return nil, err
		}
		album := Album{id, name}
		fmt.Println("id | name ")
		fmt.Printf("%3v | %8v \n", album.id, album.name)
		list = append(list, albums.Album{Id: album.id, Name: album.name})
	}
	return list, nil
}
