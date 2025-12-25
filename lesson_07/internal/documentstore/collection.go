package documentstore

type Collectable interface {
	Put(doc Document) error
	Get(key string) (*Document, bool)
	Delete(key string) bool
	List() []Document
}

type Collection struct {
	Config *CollectionConfig
	Items  map[string]*Document
}

type CollectionConfig struct {
	PrimaryKey string
}

func (s *Collection) Put(doc Document) error {
	if s.Config == nil {
		logger.Error("put document failed: collection config is nil")
		return ErrConfigNotFound
	}
	pk := s.Config.PrimaryKey
	docPrimaryKey, exist := doc.Fields[pk]
	if !exist || docPrimaryKey.Type != DocumentFieldTypeString {
		logger.Warn("put document failed: invalid primary key field", "primary_key", pk, "exist", exist, "type", func() any {
			if exist {
				return docPrimaryKey.Type
			}
			return nil
		}())
		return ErrUnsupportedDocumentField
	}

	key := docPrimaryKey.Value.(string)

	_, existed := s.Items[key]
	s.Items[key] = &doc

	if existed {
		logger.Info("document updated", "primary_key", pk, "key", key)
	} else {
		logger.Info("document created", "primary_key", pk, "key", key)
	}

	return nil
}

func (s *Collection) Get(key string) (*Document, bool) {
	item, exist := s.Items[key]
	if !exist {
		logger.Debug("document not found", "key", key)
		return nil, false
	}

	logger.Debug("document retrieved", "key", key)

	return item, true
}

func (s *Collection) Delete(key string) bool {
	if _, hasKey := s.Items[key]; !hasKey {
		logger.Warn("delete document failed: not found", "key", key)
		return false
	}

	delete(s.Items, key)
	logger.Info("document deleted", "key", key)

	return true
}

func (s *Collection) List() []Document {
	logger.Debug("list documents", "count", len(s.Items))
	docs := make([]Document, 0, len(s.Items))
	for _, d := range s.Items {
		docs = append(docs, *d)
	}
	return docs
}
