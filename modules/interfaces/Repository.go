package interfaces

type Repository interface {
	Save(document Document) (Document, error)
	FindById(id string) (Document, error)
	FindMany(ids []string) ([]Document, error)
}
