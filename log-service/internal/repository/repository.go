package repository

type Log interface {
}

type Repository struct {
	Log
}

func NewRepository() *Repository {
	return &Repository{}
}
