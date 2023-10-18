package repository

import (
	"github.com/MarselBissengaliyev/advertising/pkg/model"
	"github.com/jmoiron/sqlx"
)

type Advert interface {
	GetAll(limit int, page int, orderBy string) ([]model.GetAllAdvertsResponse, error)
	GetById(advertId int, fields []string) (model.Advert, error)
	Create(advert model.CreateAdvertsBody) (int, error)
}

type Repository struct {
	Advert
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Advert: NewAdvertPostgres(db),
	}
}
