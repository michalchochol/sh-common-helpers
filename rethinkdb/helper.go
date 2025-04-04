package rethinkdb

import (
	"fmt"
	"log"
	// "time"

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
	fmt.Printf("Stored object: %v\n", object)
}

func (rh *RethinkDBHelper) GetObject(table string, objectId string) interface{} {
	object, err := r.Table(table).Get(objectId).Run(rh.session)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Read object: %v\n", object)
	return object
}

func (rh *RethinkDBHelper) StoreObjectIfNotExists(table string, object interface{}) {
	_, err := r.Table(table).Insert(object, r.InsertOpts{Conflict: "error"}).Run(rh.session)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("Stored object: %v\n", object)
	}
}

// Funkcja zapisująca stan do bazy danych RethinkDB
// func StoreState(session *r.Session, state s.State, stateName string) {

// 	// jsonData, err := json.Marshal(state)
// 	// if err != nil {
// 	// 	fmt.Println("Error:", err)
// 	// 	return
// 	// }

// 	timestamp := time.Now()
// 	// _, err := r.Table("states").Insert(map[string]interface{}{
// 	_, err := r.Table("states").Get(stateName).Update(map[string]interface{}{
// 		"id": stateName,
// 		// "name":  stateName,
// 		"state": state,
// 		// "state":     string(jsonData),
// 		"timestamp": timestamp,
// 	}).Run(session)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Printf("Stored state: %v\n", stateName)
// }

// Funkcja subskrybująca zmiany w tabeli RethinkDB
// func SubscribeToChanges(session *r.Session, handleStateChange func(s.State)) {
// Zaczynamy nasłuchiwać zmian w tabeli 'states'
// cursor, err := r.Table("states").Changes().Run(session)
// if err != nil {
// 	log.Fatal(err)
// }
// defer cursor.Close()

// for cursor.Next() {
// 	var change r.ChangeResponse
// 	if err := cursor.Scan(&change); err != nil {
// 		log.Fatal(err)
// 	}

// 	if change.NewValue != nil && change.OldValue != nil {
// 		var newState State
// 		err := json.Unmarshal([]byte(change.NewValue.(string)), &newState)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		// Obsługuje logikę na podstawie zmienionego stanu
// 		handleStateChange(newState)
// 	}
// }
// if cursor.Err() != nil {
// 	log.Fatal(cursor.Err())
// }
// }
