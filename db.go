package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Commit struct {
	ID        int64
	Sha       string
	Additions int
	Deletions int
}

func CreateCommitsTable(db *sql.DB) {
	statement := `
	CREATE TABLE IF NOT EXISTS commits(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		sha TEXT NOT NULL UNIQUE,
		additions INTEGER,
		deletions INTEGER
	);
	`

	_, err := db.Exec(statement)
	CheckError(err)
}

func InsertCommit(db *sql.DB, sha string, additions, deletions int) int64 {
	statement := `
	INSERT INTO commits(sha, additions, deletions) VALUES (?, ?, ?)
	`

	result, err := db.Exec(statement, sha, additions, deletions)
	CheckError(err)

	id, err := result.LastInsertId()
	CheckError(err)

	return id
}

func ReadCommits(db *sql.DB) []Commit {
	rows, err := db.Query("SELECT * FROM commits")
	CheckError(err)

	defer rows.Close()

	var all []Commit

	for rows.Next() {
		var commit Commit
		if err := rows.Scan(&commit.ID, &commit.Sha, &commit.Additions, &commit.Deletions); err != nil {
			log.Fatal(err)
		}

		all = append(all, commit)
	}

	return all
}

func GetCommitById(db *sql.DB, sha string) Commit {
	row := db.QueryRow("SELECT * FROM commits WHERE sha = ?", sha)

	var commit Commit
	if err := row.Scan(&commit.ID, &commit.Sha, &commit.Additions, &commit.Deletions); err != nil {
		log.Fatal(err)
	}

	return commit
}

func isStored(db *sql.DB, sha string) bool {
	commit := GetCommitById(db, sha)

	return commit.ID != 0
}
