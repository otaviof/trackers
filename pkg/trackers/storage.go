package trackers

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	// sqlite databse drivers
	_ "github.com/mattn/go-sqlite3"
)

type storageInterface interface {
	Read() ([]*Tracker, error)
	Write(trackers []*Tracker) error
	Update(tracker *Tracker) error
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
    FROM trackers
ORDER BY hostname;
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
		tracker.Addresses = strings.Split(addresses, ",")
		tracker.Status = status

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
INSERT INTO trackers(announce, addresses, hostname, status)
     VALUES (?, ?, ?, ?);
    `); err != nil {
		return err
	}
	defer stmt.Close()

	for _, tracker = range trackers {
		log.Printf("Saving '%s'", tracker.Announce)
		if _, err = stmt.Exec(
			tracker.Announce,
			strings.Join(tracker.Addresses, ","),
			tracker.Hostname,
			tracker.Status,
		); err != nil {
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

// Update entries in table that matches hostname, setting up addresses and status.
func (s *Storage) Update(tracker *Tracker) error {
	var tx *sql.Tx
	var stmt *sql.Stmt
	var err error

	log.Printf("[Storage] Updating tracker '%s'", tracker.Announce)
	if tx, err = s.db.Begin(); err != nil {
		return err
	}

	if stmt, err = tx.Prepare(
		fmt.Sprintf(`
UPDATE trackers
   SET addresses = '%s',
       status = %d
 WHERE announce = ?
		`, strings.Join(tracker.Addresses, ","), tracker.Status),
	); err != nil {
		return err
	}
	defer stmt.Close()

	if _, err = stmt.Exec(tracker.Announce); err != nil {
		return err
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
    hostname    TEXT NOT NULL,
    status      INTEGER NOT NULL
);
    `); err != nil {
		return err
	}

	return nil
}

// NewStorage instantiate the storage backend.
func NewStorage(dbPath string) (*Storage, error) {
	var storage = &Storage{}
	var err error

	if storage.db, err = sql.Open("sqlite3", dbPath); err != nil {
		return nil, err
	}

	if err = storage.InitDB(); err != nil {
		return nil, err
	}

	return storage, nil
}
