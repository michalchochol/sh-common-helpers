package rethinkdb

import (
	"encoding/json"
	"log"

	// e "github.com/michalchochol/sh-common-helpers/error"

	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

type RethinkDBHelper struct {
	session *r.Session
}

// Inicjalizacja połączenia z bazą danych RethinkDB
func (rh *RethinkDBHelper) Init() *r.Session {
	var err error
	rh.session, err = r.Connect(r.ConnectOpts{
		Address:  "localhost:30815", // Zmień na odpowiedni adres RethinkDB
		Database: "sh_state",
	})
	if err != nil {
		log.Fatal(err)
	}
	return rh.session
}

func (rh *RethinkDBHelper) StoreObject(table string, object interface{}) {
	_, err := r.Table(table).Insert(object, r.InsertOpts{Conflict: "update"}).Run(rh.session)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Stored object: %v\n", object)
}

func (rh *RethinkDBHelper) GetObject(table string, objectId string) interface{} {
	object, err := r.Table(table).Get(objectId).Run(rh.session)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Read object: %v\n", object)
	return object
}

func (rh *RethinkDBHelper) StoreObjectIfNotExists(table string, object interface{}) {
	_, err := r.Table(table).Insert(object, r.InsertOpts{Conflict: "error"}).Run(rh.session)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Stored object: %v\n", object)
	}
}

func (rh *RethinkDBHelper) SubscribeToChanges(table string, handler func([]byte, []byte)) {

	cursor, err := r.Table(table).Changes().Run(rh.session)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close()

	var change r.ChangeResponse

	for cursor.Next(&change) {
		if change.NewValue != nil && change.OldValue != nil {
			log.Printf("Change detected: %+v\n", change)

			oldJsonBytes, err := json.Marshal(change.OldValue)
			if err != nil {
				log.Println("Error marshalling new_val:", err)
				return
			}
			newJsonBytes, err := json.Marshal(change.NewValue)
			if err != nil {
				log.Println("Error marshalling new_val:", err)
				return
			}

			handler(oldJsonBytes, newJsonBytes)
		}
	}
	if cursor.Err() != nil {
		log.Fatal(cursor.Err())
	}
}
