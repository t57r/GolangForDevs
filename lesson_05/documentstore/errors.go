package documentstore

import "errors"

var ErrDocumentNotFound = errors.New("document not found")
var ErrCollectionAlreadyExist = errors.New("collection already exists")
var ErrCollectionNotFound = errors.New("collection not found")
var ErrUnsupportedDocumentField = errors.New("unsupported document field")
