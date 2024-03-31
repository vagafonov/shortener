package storage

import (
	"bufio"
	"encoding/json"
	"os"

	"github.com/google/uuid"
	"github.com/vagafonov/shortener/pkg/entity"
)

type fileSystemStorage struct {
	file    *os.File
	encoder *json.Encoder
	scanner *bufio.Scanner
}

func NewFileSystemStorage(fileName string) (Storage, error) {
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

func (fss *fileSystemStorage) GetByHash(hash string) (*entity.URL, error) {
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

func (fss *fileSystemStorage) GetByURL(url string) (*entity.URL, error) {
	var e *entity.URL
	for fss.scanner.Scan() {
		err := json.Unmarshal(fss.scanner.Bytes(), &e)
		if err != nil {
			return nil, err
		}
		if e.Full == url {
			return e, nil
		}
	}

	return nil, nil //nolint:nilnil
}

func (fss *fileSystemStorage) Add(key string, value string) (*entity.URL, error) {
	url := &entity.URL{
		UUID:  uuid.New(),
		Short: key,
		Full:  value,
	}

	return url, fss.encoder.Encode(url)
}

func (fss *fileSystemStorage) GetAll() ([]entity.URL, error) {
	res := make([]entity.URL, 0)
	var e entity.URL
	for fss.scanner.Scan() {
		err := json.Unmarshal(fss.scanner.Bytes(), &e)
		if err != nil {
			return nil, err
		}
		res = append(res, e)
	}

	return res, nil
}

func (fss *fileSystemStorage) Truncate() {
}

func (fss *fileSystemStorage) Close() error {
	return fss.file.Close()
}
