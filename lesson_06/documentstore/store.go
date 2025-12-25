package documentstore

import (
	"bytes"
	"encoding/gob"
	"errors"
	"os"
	"path/filepath"
)

type Store struct {
	Collections map[string]*Collection
}

func NewStore() *Store {
	return &Store{
		Collections: make(map[string]*Collection),
	}
}

func (s *Store) CreateCollection(name string, cfg *CollectionConfig) (*Collection, error) {
	logger.Info("CreateCollection requested", "name", name)

	if cfg == nil {
		logger.Warn("CreateCollection failed: config not found", "name", name)
		return nil, ErrConfigNotFound
	}

	_, alreadyExist := s.Collections[name]
	if alreadyExist {
		logger.Warn("CreateCollection failed: already exist", "name", name)
		return nil, ErrCollectionAlreadyExist
	}

	collection := &Collection{
		Config: cfg,
		Items:  make(map[string]*Document),
	}
	s.Collections[name] = collection

	logger.Info("CreateCollection success", "name", name, "primary_key", cfg.PrimaryKey)

	return collection, nil
}

func (s *Store) GetCollection(name string) (*Collection, error) {
	collection, exist := s.Collections[name]
	if !exist {
		logger.Warn("GetCollection failed: not found", "name", name)
		return nil, ErrCollectionNotFound
	}

	logger.Debug("GetCollection retrieved", "name", name)

	return collection, nil
}

func (s *Store) DeleteCollection(name string) error {
	logger.Info("DeleteCollection requested", "name", name)

	if _, exists := s.Collections[name]; !exists {
		logger.Warn("DeleteCollection failed: not found", "name", name)
		return ErrCollectionNotFound
	}

	delete(s.Collections, name)
	logger.Info("DeleteCollection success", "name", name)

	return nil
}

func NewStoreFromDump(dump []byte) (*Store, error) {
	if len(dump) == 0 {
		return nil, errors.New("NewStoreFromDump: empty data")
	}

	decoder := gob.NewDecoder(bytes.NewReader(dump))

	// Decode bytes into Store
	var store *Store
	if err := decoder.Decode(&store); err != nil {
		return nil, err
	}

	return store, nil
}

func (s *Store) Dump() ([]byte, error) {
	if s == nil {
		return nil, errors.New("Dump: store is nil")
	}

	// Encode Store into bytes
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	if err := encoder.Encode(s); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func NewStoreFromFile(filename string) (*Store, error) {
	if filename == "" {
		return nil, errors.New("NewStoreFromFile: filename is empty")
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return NewStoreFromDump(data)
}

func (s *Store) DumpToFile(filename string) error {
	if filename == "" {
		return errors.New("DumpToFile: filename is empty")
	}

	data, err := s.Dump()
	if err != nil {
		return err
	}

	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	// we will use atomic write: write to temp file then rename.
	tmp, err := os.CreateTemp(dir, ".store-*.tmp")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()

	defer func() {
		_ = tmp.Close()
		_ = os.Remove(tmpName)
	}()

	if _, err := tmp.Write(data); err != nil {
		return err
	}
	if err := tmp.Sync(); err != nil {
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}

	// On Windows rename fails if target exist, so remove it first and then retry
	if err := os.Rename(tmpName, filename); err != nil {
		err = os.Remove(filename)
		return os.Rename(tmpName, filename)
	}

	return nil
}
