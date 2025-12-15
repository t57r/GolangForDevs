package documentstore

type Store struct {
	Collections map[string]*Collection
}

func NewStore() *Store {
	return &Store{
		Collections: make(map[string]*Collection),
	}
}

func (s *Store) CreateCollection(name string, cfg *CollectionConfig) (bool, *Collection) {
	_, alreadyExist := s.Collections[name]
	if alreadyExist {
		return false, nil
	}

	collection := &Collection{
		Config: cfg,
		Items:  make(map[string]*Document),
	}
	s.Collections[name] = collection
	return true, collection
}

func (s *Store) GetCollection(name string) (*Collection, bool) {
	collection, exist := s.Collections[name]
	return collection, exist
}

func (s *Store) DeleteCollection(name string) bool {
	_, hasKey := s.Collections[name]
	delete(s.Collections, name)
	return hasKey
}
