package diskstore

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"time"

	"github.com/egonelbre/event"

	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type Store struct {
	db        *sql.DB
	publisher event.Publisher
}

func New(filename string, pub event.Publisher) (*Store, error) {
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		return nil, err
	}

	s := &Store{
		db:        db,
		publisher: pub,
	}

	return s, s.init()
}

func (store *Store) Close() {
	store.db.Close()
}

func (store *Store) init() error {
	_, err := store.db.Exec(`
		CREATE TABLE IF NOT EXISTS event (
			id         CHARACTER(16),
			timestamp  DATETIME,
			type       TEXT,
			version    INTEGER,
			data       BLOB
		)
	`)

	return err
}

type record struct {
	Id        event.AggregateId
	TimeStamp time.Time
	Type      string
	Version   int64
	Data      []byte
}

func (store *Store) Load() (err error) {
	rows, err := store.db.Query(`SELECT data FROM event`)
	if err != nil {
		return err
	}
	for rows.Next() {
		var data []byte
		err := rows.Scan(&data)
		if err != nil {
			return err
		}
		v, err := fromBytes(data)
		if err != nil {
			return err
		}
		store.publisher.Publish(v)
	}
	return nil
}

func (store *Store) Save(id event.AggregateId, expectedVersion int64, events ...event.Event) (err error) {
	var tx *sql.Tx
	tx, err = store.db.Begin()
	if err != nil {
		return err
	}

	var version int64

	row := tx.QueryRow("SELECT MAX(version) FROM event WHERE id = ?", id)
	var lastVersion sql.NullInt64
	if err := row.Scan(&lastVersion); err != nil {
		tx.Rollback()
		return fmt.Errorf("did not get last version: %v", err)
	}
	if lastVersion.Valid {
		version = lastVersion.Int64
	}

	stmt, err := tx.Prepare(`
		INSERT INTO 
			event  (id, timestamp, type, version, data)
			VALUES (?, ?, ?, ?, ?)
	`)
	if err != nil {
		tx.Rollback()
		return err
	}

	if version != expectedVersion && expectedVersion != -1 {
		tx.Rollback()
		return event.ConcurrencyError
	}

	for _, event := range events {
		version += 1
		name := fmt.Sprintf("%T", event)
		data, err := toBytes(event)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to convert %T: %v", event, err)
		}
		_, err = stmt.Exec(id, time.Now(), name, version, data)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to insert %T: %v", event, err)
		}
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit: %v", err)
	}

	for _, event := range events {
		store.publisher.Publish(event)
	}

	return nil
}

func (store *Store) SaveChanges(changes event.Changes) error {
	return store.Save(changes.GetId(), changes.GetVersion(), changes.GetChanges()...)
}

func (store *Store) List(id event.AggregateId) (events []event.Info, found bool) {
	rows, err := store.db.Query(`SELECT version, timestamp, data FROM event WHERE id = ?`, id)
	if err != nil {
		return nil, false
	}
	defer rows.Close()

	events = make([]event.Info, 0, 10)
	for rows.Next() {
		var version int64
		var timestamp time.Time
		var data []byte

		err := rows.Scan(&version, &timestamp, &data)
		if err != nil {
			return nil, false
		}

		v, err := fromBytes(data)
		if err != nil {
			return nil, false
		}
		events = append(events, event.Info{id, version, timestamp, v})
	}
	return events, true
}

func init() {
	gob.Register(event.Info{})
	gob.Register(wrapper{})
}

type wrapper struct {
	Value interface{}
}

func toBytes(e event.Event) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(wrapper{e})
	return buf.Bytes(), err
}

func fromBytes(data []byte) (event.Event, error) {
	buf := bytes.NewReader(data)
	dec := gob.NewDecoder(buf)
	var wrap wrapper
	err := dec.Decode(&wrap)
	return wrap.Value, err
}
