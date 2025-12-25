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
		return ErrConfigNotFound
	}
	docPrimaryKey, exist := doc.Fields[s.Config.PrimaryKey]
	if !exist || docPrimaryKey.Type != DocumentFieldTypeString {
		return ErrUnsupportedDocumentField
	}

	s.Items[docPrimaryKey.Value.(string)] = &doc
	return nil
}

func (s *Collection) Get(key string) (*Document, bool) {
	item, exist := s.Items[key]
	return item, exist
}

func (s *Collection) Delete(key string) bool {
	_, hasKey := s.Items[key]
	delete(s.Items, key)
	return hasKey // True if the item successfully removed, False if it's not exist
}

func (s *Collection) List() []Document {
	docs := make([]Document, 0, len(s.Items))
	for _, d := range s.Items {
		docs = append(docs, *d)
	}
	return docs
}
