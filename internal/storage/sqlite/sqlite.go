package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/kavshevova/project_restapi/internal/storage"
	"github.com/mattn/go-sqlite3"
	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB //коннект до базы
}

func New(storagePath string) (*Storage, error) {
	//функция-конструктор, которая соберет и вернет объект стораджа
	const op = "storage.sqlite.New" //чтобы врапить ошибку и сразу видеть где именно ошибка

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	//что делает урлшортнер
	// есть длинная некрасивая ссылка url , есть маленькая красивая ссылка alias
	// по alias будет происходить редирект (перенаправление пользователя с одного url-адреса на другой) на url

	stmt, err := db.Prepare(`
CREATE TABLE IF NOT EXISTS urls (
    id INTEGER PRIMARY KEY,
    alias TEXT NOT NULL UNIQUE,
    url TEXT NOT NULL);
CREATE INDEX IF NOT EXISTS idx_alias ON urls(alias);
`)
	if err != nil {
		return nil, fmt.Errorf("#{op}: #{err}")
	}
	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &Storage{db: db}, nil
}

func (s *Storage) SaveURL(urlTOsave string, alias string) (int64, error) {
	const op = "storage.sqlite.SaveURL"

	stmt, err := s.db.Prepare("INSERT INTO urls (url, alias) VALUES (?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.Exec(urlTOsave, alias)
	if err != nil {
		//приводим ошибку к внутреннему типу склайт и смотрим
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrURLExists)
		}

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: failed to get last insert id: %w", op, err)
	}

	return id, nil
}

func (s *Storage) GetURL(alias string) (string, error) {
	const op = "storage.sqlite.GetURL"

	stmt, err := s.db.Prepare("SELECT url FROM urls WHERE alias = ?")
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	var resURL string
	err = stmt.QueryRow(alias).Scan(&resURL)
	if errors.Is(err, sql.ErrNoRows) { return "", storage.ErrURLNotFound
	}
	if err != nil { return "", fmt.Errorf("%s: execute statement: %w", op, err) }

	return resURL, nil
}

func (s *Storage) DeleteURL(alias  string) error {
	const op = "storage.sqlite.DeleteURL"
	stmt, err := s.db.Prepare("DELETE FROM urls WHERE alias = ?")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	_, err = stmt.Exec(alias)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
 return nil
}