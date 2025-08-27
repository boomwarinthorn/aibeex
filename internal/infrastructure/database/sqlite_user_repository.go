package database

import (
	"database/sql"
	"log"

	"fiber-hello-world/internal/domain/entity"

	_ "modernc.org/sqlite"
)

// SQLiteUserRepository implements UserRepository interface for SQLite
type SQLiteUserRepository struct {
	db *sql.DB
}

// NewSQLiteUserRepository creates a new SQLite user repository
func NewSQLiteUserRepository(db *sql.DB) *SQLiteUserRepository {
	return &SQLiteUserRepository{db: db}
}

// Create saves a new user and returns the created user with ID
func (r *SQLiteUserRepository) Create(user *entity.User) (*entity.User, error) {
	query := `
	INSERT INTO users (email, password, full_name, phone_number, birthday, created_at)
	VALUES (?, ?, ?, ?, ?, ?)
	RETURNING id`

	var id int
	err := r.db.QueryRow(query, user.Email, user.Password, user.FullName, user.PhoneNumber, user.Birthday, user.CreatedAt).Scan(&id)
	if err != nil {
		return nil, err
	}

	user.ID = id
	return user, nil
}

// GetByEmail retrieves a user by email
func (r *SQLiteUserRepository) GetByEmail(email string) (*entity.User, error) {
	query := `SELECT id, email, password, full_name, phone_number, birthday, created_at FROM users WHERE email = ?`

	var user entity.User
	err := r.db.QueryRow(query, email).Scan(&user.ID, &user.Email, &user.Password, &user.FullName, &user.PhoneNumber, &user.Birthday, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetByID retrieves a user by ID
func (r *SQLiteUserRepository) GetByID(id int) (*entity.User, error) {
	query := `SELECT id, email, password, full_name, phone_number, birthday, created_at FROM users WHERE id = ?`

	var user entity.User
	err := r.db.QueryRow(query, id).Scan(&user.ID, &user.Email, &user.Password, &user.FullName, &user.PhoneNumber, &user.Birthday, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Update updates user information
func (r *SQLiteUserRepository) Update(user *entity.User) error {
	query := `
	UPDATE users SET email = ?, full_name = ?, phone_number = ?, birthday = ?
	WHERE id = ?`

	_, err := r.db.Exec(query, user.Email, user.FullName, user.PhoneNumber, user.Birthday, user.ID)
	return err
}

// Delete removes a user by ID
func (r *SQLiteUserRepository) Delete(id int) error {
	query := `DELETE FROM users WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}

// InitDatabase initializes SQLite database and creates tables
func InitDatabase() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "users.db")
	if err != nil {
		return nil, err
	}

	// Create users table if it doesn't exist
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		email TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		full_name TEXT NOT NULL,
		phone_number TEXT NOT NULL,
		birthday TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	_, err = db.Exec(createTableQuery)
	if err != nil {
		return nil, err
	}

	log.Println("Database initialized successfully")
	return db, nil
}
