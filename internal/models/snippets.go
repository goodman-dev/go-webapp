package models

import (
	"database/sql"
	"errors"
	"time"
)

// Holds the data for an individual snippet
type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

// Wraps the sql.DB connection pool
type SnippetModel struct {
	DB *sql.DB
}

// Insert a new snippet into the db
func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {

	statement := `INSERT INTO snippets (title, content, created, expires) 
	VALUES (?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	result, err := m.DB.Exec(statement, title, content, expires)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// Get a snippet by ID
func (m *SnippetModel) Get(id int) (Snippet, error) {

	statement := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > UTC_TIMESTAMP() AND id = ?`

	// QueryRow() returns a pointer to a sql.Row object
	// row.Scan copies the values of the row into the variables
	// Note that errors from QueryRow are defered until Scan is called
	var s Snippet
	err := m.DB.QueryRow(statement, id).Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		// Returns sql.ErrNoRows error if none exist
		if errors.Is(err, sql.ErrNoRows) {
			return Snippet{}, ErrNoRecord
		} else {
			return Snippet{}, err
		}
	}

	return s, nil
}

// Get the latest 10 snippets
func (m *SnippetModel) Latest() ([]Snippet, error) {

	statement := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > UTC_TIMESTAMP() ORDER BY id DESC LIMIT 10`

	// Query method returns sql.Rows resultset
	rows, err := m.DB.Query(statement)
	if err != nil {
		return nil, err
	}

	// ¡NOTE! sql.Rows needs closing safely, since it holds an
	// open database connection from our pool
	defer rows.Close()

	var snippets []Snippet

	for rows.Next() {
		var s Snippet

		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}

		snippets = append(snippets, s)
	}

	// When we've finished iterating with rows.Next(), use rows.Err() to
	// retrieve any errors encountered during iteration
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
