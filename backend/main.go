// ----------------------------------------------------
// FILE: backend/main.go
// UPDATED to seed the database on startup.
// ----------------------------------------------------
package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

func main() {
	// --- CONFIGURATION ---
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	instanceConnectionName := os.Getenv("INSTANCE_CONNECTION_NAME")
	port := os.Getenv("PORT")
	
	if port == "" {
		port = "8080"
	}
	
	var dsn string
	if instanceConnectionName != "" {
		dsn = fmt.Sprintf("host=/cloudsql/%s user=%s password=%s dbname=%s sslmode=disable",
			instanceConnectionName, dbUser, dbPass, dbName)
	} else {
		log.Println("INSTANCE_CONNECTION_NAME not set, using local proxy configuration.")
		dsn = "host=/tmp/cloudsql/dynamic-profile-platform:us-central1:profile-db-main user=postgres dbname=profiles_db password=Lapp179623!! sslmode=disable"
	}
	// --- END CONFIGURATION ---

	store, err := NewDBStore(dsn)
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}

	if err := store.Init(); err != nil {
		log.Fatalf("Could not initialize database schema: %v", err)
	}
	log.Println("Database schema initialized successfully")

	// This resilient seeder will now log warnings instead of crashing.
	seedDatabase(store)

	server := NewAPIServer(":"+port, store)
	if err := server.Start(); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}

// Resilient seeder function
func seedDatabase(store Store) {
	log.Println("Attempting to seed database...")

	// Attempt to create or fetch the user.
	user, err := store.CreateUser("jane.doe@example.com", "password123")
	if err != nil {
		// IMPORTANT: Log as a warning and continue, do not crash the app.
		log.Printf("Warning: could not create or fetch seed user: %v", err)
	}

	// Only try to create a profile if we successfully have a user.
	if user != nil {
		err = store.CreateProfile(&Profile{
			UserID: user.ID,
			Slug:   "jane-doe",
			Name:   "Jane Doe",
			Title:  sql.NullString{String: "Principal BI Architect", Valid: true},
			Bio:    sql.NullString{String: "A results-driven engineer specializing in data.", Valid: true},
		})
		if err != nil {
			// IMPORTANT: Log as a warning and continue.
			log.Printf("Warning: could not create seed profile: %v", err)
		}
	}
	log.Println("Seeding complete. Application will now start.")
}