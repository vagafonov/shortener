package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/vagafonov/shortener/internal/contract"
	"github.com/vagafonov/shortener/pkg/entity"
)

var ErrURLNotAdded = errors.New("url not added")

type dbStorage struct {
	connection *sql.DB
}

func NewDBStorage(db *sql.DB) contract.Storage {
	return &dbStorage{
		connection: db,
	}
}

func (s *dbStorage) GetByHash(key string) (*entity.URL, error) {
	q := `SELECT id, short, original FROM url WHERE short = $1`
	row := s.connection.QueryRowContext(context.TODO(), q, key)
	var url entity.URL
	err := row.Scan(&url.UUID, &url.Short, &url.Original)
	if err != nil {
		return nil, fmt.Errorf("cannot get url by hash: %w", err)
	}

	return &url, nil
}

func (s *dbStorage) GetByURL(val string) (*entity.URL, error) {
	q := `SELECT id, short, original FROM url WHERE original = $1`
	row := s.connection.QueryRowContext(context.TODO(), q, val)
	if row.Err() != nil {
		return nil, fmt.Errorf("sql query error: %w", row.Err())
	}
	var url entity.URL
	err := row.Scan(&url.UUID, &url.Short, &url.Original)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil //nolint:nilnil
		}

		return nil, fmt.Errorf("cannot get url by original url: %w", err)
	}

	return &url, nil
}

func (s *dbStorage) Add(hash string, url string) (*entity.URL, error) {
	q := `INSERT INTO url (id, short, original) VALUES ($1, $2, $3)`
	id := uuid.New()
	res, err := s.connection.ExecContext(context.TODO(), q, id, hash, url)
	if err != nil {
		return nil, fmt.Errorf("cannot add url: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("cannot get affected rows whean add url: %w", err)
	}

	if rows == 0 {
		return nil, ErrURLNotAdded
	}

	return &entity.URL{
		UUID:     id,
		Short:    hash,
		Original: url,
	}, nil
}

func (s *dbStorage) GetAll() ([]entity.URL, error) {
	q := `SELECT id, short, original FROM url LIMIT 1000`
	rows, err := s.connection.QueryContext(context.TODO(), q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	urls := make([]entity.URL, 0)
	for rows.Next() {
		var u entity.URL
		err = rows.Scan(&u.UUID, &u.Short, &u.Original)
		if err != nil {
			return nil, fmt.Errorf("cannot get all urls: %w", err)
		}
		urls = append(urls, u)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("cannot scan all urls: %w", err)
	}

	return urls, nil
}

func (s *dbStorage) Truncate() {
}

func (s *dbStorage) Close() error {
	return nil
}
