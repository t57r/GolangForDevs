package documentstore

import (
	"fmt"

	"github.com/google/btree"
)

type Collectable interface {
	Put(doc Document) error
	Get(key string) (*Document, bool)
	Delete(key string) bool
	List() []Document
}

type Collection struct {
	Config  CollectionConfig
	Items   map[string]*Document
	indexes map[string]*btree.BTree
}

type CollectionConfig struct {
	PrimaryKey string
}

type QueryParams struct {
	Desc     bool
	MinValue *string // inclusive
	MaxValue *string // inclusive
}

// indexItem sorts by value first, then docKey
type indexItem struct {
	value  string
	docKey string
}

func (a indexItem) Less(b btree.Item) bool {
	o := b.(indexItem)
	if a.value == o.value {
		return a.docKey < o.docKey
	}
	return a.value < o.value
}

func (s *Collection) Put(doc Document) error {
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

	old, existed := s.Items[key]

	if existed {
		logger.Info("document updated", "primary_key", pk, "key", key)
		// if updating existing, remove old index entries
		if old != nil {
			if err := s.removeFromIndexes(key, old); err != nil {
				return err
			}
		}
	} else {
		logger.Info("document created", "primary_key", pk, "key", key)
	}

	// store document
	s.Items[key] = &doc

	// add new index entries
	if err := s.addToIndexes(key, &doc); err != nil {
		return err
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
	old, hasKey := s.Items[key]
	if !hasKey || old == nil {
		logger.Warn("delete document failed: not found", "key", key)
		return false
	}

	_ = s.removeFromIndexes(key, old)
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

// --- Index related methods ---
func getIndexableStringField(doc *Document, fieldName string) (value string, ok bool, err error) {
	f, exists := doc.Fields[fieldName]
	if !exists {
		return "", false, nil
	}
	if f.Type != DocumentFieldTypeString {
		return "", false, ErrIndexFieldNotIndexable
	}
	s, ok := f.Value.(string)
	if !ok {
		return "", false, ErrIndexFieldNotIndexable
	}
	return s, true, nil
}

func (s *Collection) removeFromIndexes(docKey string, doc *Document) error {
	if s.indexes == nil {
		return nil
	}
	for fieldName, tree := range s.indexes {
		v, ok, err := getIndexableStringField(doc, fieldName)
		if err != nil {
			return err
		}
		if !ok {
			continue
		}
		tree.Delete(indexItem{value: v, docKey: docKey})
	}
	return nil
}

func (s *Collection) addToIndexes(docKey string, doc *Document) error {
	if s.indexes == nil {
		return nil
	}
	for fieldName, tree := range s.indexes {
		v, ok, err := getIndexableStringField(doc, fieldName)
		if err != nil {
			return err
		}
		if !ok {
			continue
		}
		tree.ReplaceOrInsert(indexItem{value: v, docKey: docKey})
	}
	return nil
}

func (s *Collection) CreateIndex(fieldName string) error {
	if s.indexes == nil {
		s.indexes = make(map[string]*btree.BTree)
	}

	if _, exists := s.indexes[fieldName]; exists {
		return ErrIndexAlreadyExist
	}

	tree := btree.New(32)

	for docKey, doc := range s.Items {
		v, ok, err := getIndexableStringField(doc, fieldName)
		if err != nil {
			return fmt.Errorf("CreateIndex(%s): %w", fieldName, err)
		}
		if !ok {
			continue // field missing -> just skip
		}
		tree.ReplaceOrInsert(indexItem{value: v, docKey: docKey})
	}

	s.indexes[fieldName] = tree
	return nil
}

func (s *Collection) DeleteIndex(fieldName string) error {
	if s.indexes == nil {
		return ErrIndexNotFound
	}
	if _, exists := s.indexes[fieldName]; !exists {
		return ErrIndexNotFound
	}
	delete(s.indexes, fieldName)
	return nil
}

func (s *Collection) Query(fieldName string, params QueryParams) ([]Document, error) {
	if s.indexes == nil {
		return nil, ErrIndexNotFound
	}
	tree, exists := s.indexes[fieldName]
	if !exists {
		return nil, ErrIndexNotFound
	}

	out := make([]Document, 0)

	withinMax := func(v string) bool {
		return params.MaxValue == nil || v <= *params.MaxValue
	}
	withinMin := func(v string) bool {
		return params.MinValue == nil || v >= *params.MinValue
	}

	emit := func(it indexItem) bool {
		if !withinMin(it.value) {
			if params.Desc {
				return false
			}
			return true
		}
		if !withinMax(it.value) {
			if !params.Desc {
				return false
			}
			return true
		}

		docPtr := s.Items[it.docKey]
		if docPtr == nil {
			return true
		}
		out = append(out, *docPtr)
		return true
	}

	if !params.Desc {
		// Ascending
		if params.MinValue != nil {
			start := indexItem{value: *params.MinValue, docKey: ""} // "" gives earliest docKey for that value
			tree.AscendGreaterOrEqual(start, func(i btree.Item) bool {
				return emit(i.(indexItem))
			})
		} else {
			tree.Ascend(func(i btree.Item) bool {
				return emit(i.(indexItem))
			})
		}
		return out, nil
	}

	// Descending
	if params.MaxValue != nil {
		start := indexItem{value: *params.MaxValue, docKey: "\U0010FFFF"} // high docKey to include equals
		tree.DescendLessOrEqual(start, func(i btree.Item) bool {
			return emit(i.(indexItem))
		})
	} else {
		tree.Descend(func(i btree.Item) bool {
			return emit(i.(indexItem))
		})
	}

	return out, nil
}
