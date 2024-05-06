package storage

import (
	"bufio"
	"context"
	"encoding/json"
	"os"

	"github.com/google/uuid"
	"github.com/vagafonov/shortener/internal/contract"
	"github.com/vagafonov/shortener/pkg/entity"
)

type fileSystemStorage struct {
	file    *os.File
	encoder *json.Encoder
	scanner *bufio.Scanner
}

func NewFileSystemStorage(fileName string) (contract.Storage, error) {
	fss := fileSystemStorage{}
	var err error
	fss.file, err = os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666) //nolint:gofumpt, gomnd
	if err != nil {
		return nil, err
	}
	fss.encoder = json.NewEncoder(fss.file)
	fss.scanner = bufio.NewScanner(fss.file)

	return &fss, nil
}

func (fss *fileSystemStorage) GetByHash(ctx context.Context, hash string) (*entity.URL, error) {
	var e *entity.URL
	for fss.scanner.Scan() {
		err := json.Unmarshal(fss.scanner.Bytes(), &e)
		if err != nil {
			return nil, err
		}
		if e.Short == hash {
			return e, nil
		}
	}

	return nil, nil //nolint:nilnil
}

func (fss *fileSystemStorage) GetByURL(ctx context.Context, url string) (*entity.URL, error) {
	var e *entity.URL
	for fss.scanner.Scan() {
		err := json.Unmarshal(fss.scanner.Bytes(), &e)
		if err != nil {
			return nil, err
		}
		if e.Original == url {
			return e, nil
		}
	}

	return nil, nil //nolint:nilnil
}

func (fss *fileSystemStorage) Add(
	ctx context.Context,
	key string,
	value string,
	userID uuid.UUID,
) (*entity.URL, error) {
	url := &entity.URL{
		UUID:     uuid.New(),
		Short:    key,
		Original: value,
		UserID:   userID,
	}

	return url, fss.encoder.Encode(url)
}

func (fss *fileSystemStorage) GetAll(ctx context.Context) ([]*entity.URL, error) {
	res := make([]*entity.URL, 0)
	var e entity.URL
	for fss.scanner.Scan() {
		err := json.Unmarshal(fss.scanner.Bytes(), &e)
		if err != nil {
			return nil, err
		}
		res = append(res, &e)
	}

	return res, nil
}

func (fss *fileSystemStorage) AddBatch(ctx context.Context, b []*entity.URL) (int, error) {
	encoder := json.NewEncoder(fss.file)
	for _, v := range b {
		err := encoder.Encode(v)
		if err != nil {
			return 0, err
		}
	}

	return len(b), nil
}

func (fss *fileSystemStorage) GetAllURLsByUser(
	ctx context.Context,
	userID uuid.UUID,
	baseURL string,
) ([]*entity.URL, error) {
	res := make([]*entity.URL, 0)
	var e entity.URL
	for fss.scanner.Scan() {
		err := json.Unmarshal(fss.scanner.Bytes(), &e)
		if err != nil {
			return nil, err
		}

		if e.UserID == userID {
			res = append(res, &e)
		}
	}

	return res, nil
}

// DeleteURLsByUser TODO need implement.
func (fss *fileSystemStorage) DeleteURLsByUser(ctx context.Context, userID uuid.UUID, batch []string) error {
	return nil
}

func (fss *fileSystemStorage) Ping(ctx context.Context) error {
	return nil
}

func (fss *fileSystemStorage) Truncate() {
}

func (fss *fileSystemStorage) Close() error {
	return fss.file.Close()
}
