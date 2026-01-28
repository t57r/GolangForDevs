package documentstore

import "errors"

var (
	ErrConfigNotFound           = errors.New("config must be initialized")
	ErrDocumentNotFound         = errors.New("document not found")
	ErrCollectionAlreadyExist   = errors.New("collection already exists")
	ErrCollectionNotFound       = errors.New("collection not found")
	ErrUnsupportedDocumentField = errors.New("unsupported document field")

	ErrIndexAlreadyExist      = errors.New("index already exists")
	ErrIndexNotFound          = errors.New("index not found")
	ErrIndexFieldNotIndexable = errors.New("only string field can be indexed")
)
