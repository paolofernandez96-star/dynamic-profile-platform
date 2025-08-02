// ----------------------------------------------------
// FILE: backend/store.go
// UPDATED with CreateUser and CreateProfile functions.
// ----------------------------------------------------
package main

import (
	"database/sql"
	"log"
	"os"
	_ "github.com/lib/pq"
)

type Store interface {
	Init() error
	GetProfileBySlug(slug string) (*Profile, error)
	CreateUser(email, password string) (*User, error)
	CreateProfile(profile *Profile) error
}

type DBStore struct {
	db *sql.DB
}

func NewDBStore(connStr string) (*DBStore, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	log.Println("Database connection successful")
	return &DBStore{db: db}, nil
}

func (s *DBStore) Init() error {
	schema, err := os.ReadFile("schema.sql")
	if err != nil {
		return err
	}
	_, err = s.db.Exec(string(schema))
	return err
}

func (s *DBStore) CreateUser(email, password string) (*User, error) {
	// In a real app, the password would be hashed.
	query := `INSERT INTO users (email, hashed_password) VALUES ($1, $2)
	          ON CONFLICT (email) DO NOTHING
			  RETURNING user_id, email, created_at`

	user := new(User)
	err := s.db.QueryRow(query, email, password).Scan(&user.ID, &user.Email, &user.CreatedAt)
	if err != nil {
		// If no rows are returned because of the conflict, fetch the existing user.
		if err == sql.ErrNoRows {
			log.Printf("User with email %s already exists, fetching...", email)
			return s.GetUserByEmail(email)
		}
		return nil, err
	}
	return user, nil
}

// Helper function to get a user if they already exist.
func (s *DBStore) GetUserByEmail(email string) (*User, error) {
	query := `SELECT user_id, email, created_at FROM users WHERE email = $1`
	user := new(User)
	err := s.db.QueryRow(query, email).Scan(&user.ID, &user.Email, &user.CreatedAt)
	return user, err
}

func (s *DBStore) CreateProfile(profile *Profile) error {
	query := `INSERT INTO profiles (user_id, profile_slug, name, title, bio)
	          VALUES ($1, $2, $3, $4, $5)
			  ON CONFLICT (profile_slug) DO NOTHING`

	_, err := s.db.Exec(
		query,
		profile.UserID,
		profile.Slug,
		profile.Name,
		profile.Title,
		profile.Bio,
	)
	return err
}

func (s *DBStore) GetProfileBySlug(slug string) (*Profile, error) {
	row := s.db.QueryRow("SELECT profile_id, profile_slug, name, title, bio, created_at FROM profiles WHERE profile_slug = $1", slug)

	p := new(Profile)
	err := row.Scan(
		&p.ID,
		&p.Slug,
		&p.Name,
		&p.Title,
		&p.Bio,
		&p.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("No profile found for slug: %s", slug)
		}
		return nil, err
	}
	return p, nil
}