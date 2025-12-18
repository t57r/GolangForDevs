package documentstore

type DocumentFieldType string

const (
	DocumentFieldTypeString DocumentFieldType = "string"
	DocumentFieldTypeNumber DocumentFieldType = "number"
	DocumentFieldTypeBool   DocumentFieldType = "bool"
	DocumentFieldTypeArray  DocumentFieldType = "array"
	DocumentFieldTypeObject DocumentFieldType = "object"
)

type DocumentField struct {
	Type  DocumentFieldType
	Value interface{}
}

type Document struct {
	Fields map[string]DocumentField
}

var documents = map[string]*Document{}

func Put(doc *Document) {
	field, exist := doc.Fields["key"]
	if !exist {
		return // document must contain field "key"
	}

	if field.Type != DocumentFieldTypeString {
		return // field "key" must be of type string
	}

	key, ok := field.Value.(string)
	if !ok {
		return // field "key" must be of type string
	}

	documents[key] = doc
}

func Get(key string) (*Document, bool) {
	document, exist := documents[key]
	if !exist {
		return nil, false
	}

	return document, true
}

func Delete(key string) bool {
	_, exist := documents[key]
	delete(documents, key) // no-op if key doesn't exist
	return exist
}

func List() []*Document {
	docs := make([]*Document, 0, len(documents))
	for _, d := range documents {
		docs = append(docs, d)
	}
	return docs
}
