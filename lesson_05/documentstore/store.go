package documentstore

type Store struct {
	Collections map[string]Collectable
}

func NewStore() *Store {
	return &Store{
		Collections: make(map[string]Collectable),
	}
}

func (s *Store) CreateCollection(name string, cfg *CollectionConfig) (Collectable, error) {
	if cfg == nil {
		return nil, ErrConfigNotFound
	}
	_, alreadyExist := s.Collections[name]
	if alreadyExist {
		return nil, ErrCollectionAlreadyExist
	}

	collection := &Collection{
		Config: cfg,
		Items:  make(map[string]*Document),
	}
	s.Collections[name] = collection
	return collection, nil
}

func (s *Store) GetCollection(name string) (Collectable, error) {
	collection, exist := s.Collections[name]
	if !exist {
		return nil, ErrCollectionNotFound
	}
	return collection, nil
}

func (s *Store) DeleteCollection(name string) error {
	_, hasKey := s.Collections[name]
	delete(s.Collections, name)
	if !hasKey {
		return ErrCollectionNotFound
	}
	return nil
}

