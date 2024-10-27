package main

import (
	"fmt"
	"jobstr-relay/policies"
	"net/http"
	"os"
	"path/filepath"

	"github.com/fiatjaf/eventstore/sqlite3"
	"github.com/fiatjaf/khatru"
)

func main() {
	relay := khatru.NewRelay()

	relay.Info.Name = "Jobstr Relay"
	relay.Info.PubKey = ""
	relay.Info.Description = "this is my custom relay"
	relay.Info.Icon = "https://i.ytimg.com/vi/XeTfcLTKvdA/hqdefault.jpg"

	dbPath := "./db/jobstr.db"
	dbDir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		panic(fmt.Sprintf("Failed to create database directory: %v", err))
	}
	if err := os.Chmod(dbDir, 0755); err != nil {
		panic(fmt.Sprintf("Failed to set database directory permissions: %v", err))
	}

	db := sqlite3.SQLite3Backend{DatabaseURL: dbPath}
	if err := db.Init(); err != nil {
		panic(fmt.Sprintf("Failed to initialize database: %v", err))
	}

	relay.RejectEvent = append(relay.RejectEvent,
		policies.RestrictToSpecifiedKinds(false, 1),
		policies.VerifyMessagePattern())
	relay.StoreEvent = append(relay.StoreEvent, db.SaveEvent)
	relay.QueryEvents = append(relay.QueryEvents, db.QueryEvents)
	relay.CountEvents = append(relay.CountEvents, db.CountEvents)
	relay.DeleteEvent = append(relay.DeleteEvent, db.DeleteEvent)

	fmt.Println("running on :3334")
	http.ListenAndServe(":3334", relay)
}
