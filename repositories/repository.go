package repositories

type User interface {
}

type Repository struct {
	User
}

func NewRepository() *Repository {
	return &Repository{}
}
