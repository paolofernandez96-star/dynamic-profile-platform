// ----------------------------------------------------
// FILE: backend/types.go
// UPDATED to handle nullable text fields from the database.
// ----------------------------------------------------
package main

import (
	"database/sql" // Added for sql.NullString
	"time"
	"github.com/google/uuid"
)

type Profile struct {
	ID                uuid.UUID      `json:"id"`
	UserID            uuid.UUID      `json:"-"`
	Slug              string         `json:"slug"`
	Name              string         `json:"name"`
	Title             sql.NullString `json:"title"` // Handles nullable fields
	Bio               sql.NullString `json:"bio"`   // Handles nullable fields
	ProfilePictureURL string         `json:"profilePictureUrl"`
	CreatedAt         time.Time      `json:"createdAt"`
}

type User struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"createdAt"`
}