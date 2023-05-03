package interfaces

type Indexer interface {
	IndexFromDocument(document Document) (*Document, error)
}
