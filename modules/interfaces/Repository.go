package interfaces

type Repository interface {
	Save(document Document) (Document, error)
	GetMany() ([]Document, error)
}
