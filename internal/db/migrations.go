package db

import (
	"database/sql"
	"fmt"
	"time"
)

type migration struct {
	version int
	apply   func(tx *sql.Tx) error
}

var migrations = []migration{
	{
		version: 1,
		apply: func(tx *sql.Tx) error {
			if _, err := tx.Exec(`
				CREATE TABLE IF NOT EXISTS users (
					id INTEGER PRIMARY KEY AUTOINCREMENT,
					username TEXT UNIQUE NOT NULL,
					passhash TEXT NOT NULL,
					xp INTEGER DEFAULT 0,
					coins INTEGER DEFAULT 0,
					social TEXT DEFAULT '',
					upload_token TEXT DEFAULT ''
				);
			`); err != nil {
				return fmt.Errorf("create users table: %w", err)
			}

			if _, err := tx.Exec(`
				CREATE TABLE IF NOT EXISTS posts (
					id INTEGER PRIMARY KEY AUTOINCREMENT,
					user_id INTEGER NOT NULL,
					created_at INTEGER NOT NULL,
					updated_at INTEGER DEFAULT 0,
					title TEXT DEFAULT '',
					has_image INTEGER DEFAULT 0,
					path TEXT DEFAULT '',
					content TEXT DEFAULT '',
					thumbnail_path TEXT DEFAULT '',
					color TEXT DEFAULT ''
				);
			`); err != nil {
				return fmt.Errorf("create posts table: %w", err)
			}
			return nil
		},
	},
	{
		version: 2,
		apply: func(tx *sql.Tx) error {
			// Ensure legacy abcms.sqlite databases from earlier Lapis migrations
			// have all columns expected by the Go queries.
			userCols := []struct {
				name string
				ddl  string
			}{
				{"xp", "INTEGER DEFAULT 0"},
				{"coins", "INTEGER DEFAULT 0"},
				{"social", "TEXT DEFAULT ''"},
				{"upload_token", "TEXT DEFAULT ''"},
			}
			for _, col := range userCols {
				if err := ensureColumn(tx, "users", col.name, col.ddl); err != nil {
					return err
				}
			}

			postCols := []struct {
				name string
				ddl  string
			}{
				{"title", "TEXT DEFAULT ''"},
				{"has_image", "INTEGER DEFAULT 0"},
				{"path", "TEXT DEFAULT ''"},
				{"content", "TEXT DEFAULT ''"},
				{"thumbnail_path", "TEXT DEFAULT ''"},
				{"color", "TEXT DEFAULT ''"},
				{"updated_at", "INTEGER DEFAULT 0"},
			}
			for _, col := range postCols {
				if err := ensureColumn(tx, "posts", col.name, col.ddl); err != nil {
					return err
				}
			}
			return nil
		},
	},
	{
		version: 3,
		apply: func(tx *sql.Tx) error {
			if _, err := tx.Exec(`CREATE INDEX IF NOT EXISTS idx_posts_created_at ON posts(created_at DESC, id DESC)`); err != nil {
				return fmt.Errorf("create idx_posts_created_at: %w", err)
			}
			if _, err := tx.Exec(`CREATE INDEX IF NOT EXISTS idx_posts_user_id ON posts(user_id)`); err != nil {
				return fmt.Errorf("create idx_posts_user_id: %w", err)
			}
			return nil
		},
	},
}

// Migrate brings the SQLite schema up to date while preserving existing Lapis data.
func Migrate(conn *sql.DB) error {
	if _, err := conn.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version INTEGER PRIMARY KEY,
			applied_at INTEGER NOT NULL
		);
	`); err != nil {
		return fmt.Errorf("create schema_migrations: %w", err)
	}

	for _, m := range migrations {
		var exists int
		err := conn.QueryRow(`SELECT COUNT(1) FROM schema_migrations WHERE version = ?`, m.version).Scan(&exists)
		if err != nil {
			return fmt.Errorf("check migration %d: %w", m.version, err)
		}
		if exists > 0 {
			continue
		}

		tx, err := conn.Begin()
		if err != nil {
			return fmt.Errorf("begin migration %d: %w", m.version, err)
		}
		if err := m.apply(tx); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("apply migration %d: %w", m.version, err)
		}
		if _, err := tx.Exec(`INSERT INTO schema_migrations (version, applied_at) VALUES (?, ?)`, m.version, time.Now().Unix()); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("record migration %d: %w", m.version, err)
		}
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("commit migration %d: %w", m.version, err)
		}
	}

	return nil
}

func ensureColumn(tx *sql.Tx, table, column, ddl string) error {
	rows, err := tx.Query(fmt.Sprintf("PRAGMA table_info(%s)", table))
	if err != nil {
		return fmt.Errorf("inspect %s columns: %w", table, err)
	}
	defer rows.Close()

	for rows.Next() {
		var cid int
		var name, colType string
		var notNull int
		var dflt sql.NullString
		var pk int
		if err := rows.Scan(&cid, &name, &colType, &notNull, &dflt, &pk); err != nil {
			return fmt.Errorf("scan table_info(%s): %w", table, err)
		}
		if name == column {
			return nil
		}
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate table_info(%s): %w", table, err)
	}

	if _, err := tx.Exec(fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s", table, column, ddl)); err != nil {
		return fmt.Errorf("add column %s.%s: %w", table, column, err)
	}
	return nil
}
