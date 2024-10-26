package main

import (
	"context"
	"fmt"
	"jobstr-relay/policies"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/fiatjaf/eventstore/sqlite3"
	"github.com/fiatjaf/khatru"
	"github.com/nbd-wtf/go-nostr"
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
		policies.RestrictToSpecifiedKinds(),
		policies.PreventLargeTags(),
		policies.PreventTooManyIndexableTags(),
		func(ctx context.Context, event *nostr.Event) (reject bool, msg string) {
			lines := strings.Split(strings.TrimSpace(event.Content), "\n")
			if len(lines) != 2 ||
				!strings.HasPrefix(lines[0], "Job Title: ") ||
				!strings.HasPrefix(lines[1], "description: ") {
				return true, "Event rejected: invalid format. Use 'Job Title:' on the first line and 'description:' on the second line."
			}
			return false, ""
		})
	relay.StoreEvent = append(relay.StoreEvent, db.SaveEvent)
	relay.QueryEvents = append(relay.QueryEvents, db.QueryEvents)
	relay.CountEvents = append(relay.CountEvents, db.CountEvents)
	relay.DeleteEvent = append(relay.DeleteEvent, db.DeleteEvent)

	fmt.Println("running on :3334")
	http.ListenAndServe(":3334", relay)
}
