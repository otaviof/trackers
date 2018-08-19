package trackers

import (
	"database/sql"
	"log"
	"strings"

	// sqlite databse drivers
	_ "github.com/mattn/go-sqlite3"
)

type storageInterface interface {
	Read() ([]*Tracker, error)
	Write(trackers []*Tracker) error
}

// Storage persist trackers in sqlite
type Storage struct {
	db *sql.DB
}

// Read data from sqlite.
func (s *Storage) Read() ([]*Tracker, error) {
	var rows *sql.Rows
	var trackers []*Tracker
	var tracker *Tracker
	var err error

	if rows, err = s.db.Query(`
SELECT announce, addresses, status
  FROM trackers;
    `); err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var announce string
		var addresses string
		var status int

		if err = rows.Scan(&announce, &addresses, &status); err != nil {
			return nil, err
		}
		if tracker, err = NewTracker(announce); err != nil {
			return nil, err
		}
		if status != 99 {
			tracker.Addresses = strings.Split(addresses, ",")
		}

		trackers = append(trackers, tracker)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return trackers, nil
}

// Write write a list of trackers to the database.
func (s *Storage) Write(trackers []*Tracker) error {
	var tx *sql.Tx
	var tracker *Tracker
	var stmt *sql.Stmt
	var err error

	if tx, err = s.db.Begin(); err != nil {
		return err
	}

	if stmt, err = tx.Prepare(`
INSERT INTO trackers(announce, addresses, status)
     VALUES (?, ?, ?)
    `); err != nil {
		return err
	}
	defer stmt.Close()

	for _, tracker = range trackers {
		var addresses string
		var status int

		if len(tracker.Addresses) > 0 {
			addresses = strings.Join(tracker.Addresses, ",")
		} else {
			addresses = "0.0.0.0"
			status = 99
		}

		log.Printf("Saving '%s'", tracker.Announce)
		if _, err = stmt.Exec(tracker.Announce, addresses, status); err != nil {
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

// InitDB creates the initial structure to store tracker records.
func (s *Storage) InitDB() error {
	var err error

	if _, err = s.db.Exec(`
CREATE TABLE IF NOT EXISTS trackers (
    announce    TEXT NOT NULL PRIMARY KEY,
    addresses   TEXT NOT NULL,
    status      INTEGER NOT NULL
);
    `); err != nil {
		return err
	}

	return nil
}

// NewStorage instantiate the storage backend.
func NewStorage() (*Storage, error) {
	var storage = &Storage{}
	var err error

	if storage.db, err = sql.Open("sqlite3", "/var/tmp/test.sqlite"); err != nil {
		return nil, err
	}

	if err = storage.InitDB(); err != nil {
		return nil, err
	}

	return storage, nil
}
