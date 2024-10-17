package db

import (
	"database/sql"
	"log"
	"os"
	"path"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

var connection *sql.DB = nil

func open() *sql.DB {
	dbDir, err := os.Executable()
	if err != nil {
		panic(err)
	}

	p := path.Join(filepath.Dir(dbDir), "gomato.db")
	log.Println(p)
	db, err := sql.Open("sqlite3", p)
	if err != nil {
		log.Fatal(err)
	}

	init_db(db)
	return db
}

func GetDb() *sql.DB {
	if connection == nil {
		connection = open()
	}
	return connection
}

func init_db(db *sql.DB) {
	q := `
CREATE TABLE IF NOT EXISTS configuration (
	id INTEGER PRIMARY KEY,
	[name] text UNIQUE NOT NULL,
	work_time INTEGER NOT NULL,
	pause_time INTEGER NOT NULL,
	-- 0: none, 1: audio, 2: audio and visual
	notification_level INTEGER NOT NULL DEFAULT (2)
);

CREATE TABLE IF NOT EXISTS defaults (
    id INTEGER PRIMARY KEY,
	default_config INTEGER REFERENCES configuration
);

INSERT INTO configuration
(
    id, [name], work_time, pause_time, notification_level
) VALUES (
  1, 'default', 1500000000000, 300000000000, 2
) ON CONFLICT(id) DO NOTHING;

INSERT INTO defaults 
(
	id, default_config
) VALUES (
	1, 1
) ON CONFLICT(id) DO NOTHING;`
	_, err := db.Exec(q)
	if err != nil {
		panic(err)
	}
}
