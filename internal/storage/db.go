package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/vagafonov/shortener/internal/contract"
	"github.com/vagafonov/shortener/internal/customerror"
	"github.com/vagafonov/shortener/pkg/entity"
)

const batchInsertSize = 100

type dbStorage struct {
	connection *sql.DB
}

// NewDBStorage Constructor for DBStorage.
func NewDBStorage(db *sql.DB) contract.Storage {
	return &dbStorage{
		connection: db,
	}
}

// GetByHash get short urls by hash from database.
func (s *dbStorage) GetByHash(ctx context.Context, key string) (*entity.URL, error) {
	q := `SELECT id, short, original, user_id, deleted_at FROM urls WHERE short = $1`
	row := s.connection.QueryRowContext(ctx, q, key)
	var url entity.URL
	err := row.Scan(&url.UUID, &url.Short, &url.Original, &url.UserID, &url.DeletedAt)
	if err != nil {
		return nil, fmt.Errorf("cannot get url by hash: %w", err)
	}

	if url.DeletedAt != nil {
		return nil, customerror.ErrURLDeleted
	}

	return &url, nil
}

// GetByURL get short urls by url from database.
func (s *dbStorage) GetByURL(ctx context.Context, val string) (*entity.URL, error) {
	q := `SELECT id, short, original FROM urls WHERE original = $1`
	row := s.connection.QueryRowContext(ctx, q, val)
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

// Add create new short url in database.
func (s *dbStorage) Add(ctx context.Context, hash string, url string, userID uuid.UUID) (*entity.URL, error) {
	q := `INSERT INTO urls (id, short, original, user_id) VALUES ($1, $2, $3, $4)`
	id := uuid.New()
	res, err := s.connection.ExecContext(ctx, q, id, hash, url, userID)
	if err != nil {
		return nil, fmt.Errorf("cannot add url: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("cannot get affected rows whean add url: %w", err)
	}

	if rows == 0 {
		return nil, customerror.ErrURLNotAdded
	}

	return &entity.URL{
		UUID:     id,
		Short:    hash,
		Original: url,
		UserID:   userID,
	}, nil
}

// GetAll get all short urls from database.
func (s *dbStorage) GetAll(ctx context.Context) ([]*entity.URL, error) {
	q := `SELECT id, short, original FROM urls LIMIT 1000`
	rows, err := s.connection.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	urls := make([]*entity.URL, 0)
	for rows.Next() {
		var u entity.URL
		err = rows.Scan(&u.UUID, &u.Short, &u.Original)
		if err != nil {
			return nil, fmt.Errorf("cannot get all urls: %w", err)
		}
		urls = append(urls, &u)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("cannot scan all urls: %w", err)
	}

	return urls, nil
}

// AddBatch create multiple short  urls in database.
func (s *dbStorage) AddBatch(ctx context.Context, b []*entity.URL) (int, error) {
	bufIns := make([]*entity.URL, 0)
	inserted := 0
	for _, v := range b {
		v.UUID = uuid.New()
		bufIns = append(bufIns, v)
		if len(bufIns) == batchInsertSize {
			if err := s.batchInsert(ctx, bufIns); err != nil {
				return inserted, err
			}
			inserted = +len(bufIns)
			bufIns = nil
		}
	}

	if err := s.batchInsert(ctx, bufIns); err != nil {
		return 0, err
	}
	inserted = +len(bufIns)

	return inserted, nil
}

func (s *dbStorage) batchInsert(ctx context.Context, urls []*entity.URL) error {
	tx, err := s.connection.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback() //nolint:errcheck

	stmt, err := tx.PrepareContext(ctx, "INSERT INTO urls (id, short, original, user_id) VALUES($1, $2, $3, $4)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, u := range urls {
		_, err := stmt.ExecContext(ctx, u.UUID, u.Short, u.Original, u.UserID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// GetAllURLsByUser Get all urls by user from database.
func (s *dbStorage) GetAllURLsByUser(ctx context.Context, userID uuid.UUID, baseURL string) ([]*entity.URL, error) {
	q := `SELECT id, short, original FROM urls WHERE user_id = $1 LIMIT 1000`

	rows, err := s.connection.QueryContext(ctx, q, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	urls := make([]*entity.URL, 0)
	for rows.Next() {
		var u entity.URL
		err = rows.Scan(&u.UUID, &u.Short, &u.Original)
		if err != nil {
			return nil, fmt.Errorf("cannot get all urls by user: %w", err)
		}
		u.Short = fmt.Sprintf("%s/%s", baseURL, u.Short)
		urls = append(urls, &u)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("cannot scan all urls by user: %w", err)
	}

	return urls, nil
}

// DeleteURLsByUser delete URLS by user.
func (s *dbStorage) DeleteURLsByUser(ctx context.Context, userID uuid.UUID, batch []string) error {
	tx, err := s.connection.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction for delete urls: %w", err)
	}

	defer func() {
		err = tx.Rollback()
	}()
	if err != nil {
		return fmt.Errorf("failed to rollback transaction for delete urls: %w", err)
	}

	stmt, err := s.connection.PrepareContext(ctx, "UPDATE urls SET deleted_at = NOW() WHERE short = $1")
	if err != nil {
		return fmt.Errorf("failed to prepare context for delete urls: %w", err)
	}
	defer func() {
		err = stmt.Close()
	}()
	if err != nil {
		return fmt.Errorf("failed to close stmt for delete urls: %w", err)
	}

	for _, v := range batch {
		// _, err := stmt.ExecContext(ctx, userID, v)
		_, err := stmt.ExecContext(ctx, v)
		if err != nil {
			return fmt.Errorf("failed to exec context for delete urls: %w", err)
		}
	}

	return tx.Commit()
}

// Ping not implemented.
func (s *dbStorage) Ping(ctx context.Context) error {
	return s.connection.PingContext(ctx)
}

// Truncate not implemented.
func (s *dbStorage) Truncate() {
}

// Close not implemented.
func (s *dbStorage) Close() error {
	return nil
}
