package repository

import (
	"encoding/json"
	"fmt"

	"github.com/MarselBissengaliyev/advertising/pkg/model"
	"github.com/jmoiron/sqlx"
)

type AdvertisingPostgres struct {
	db *sqlx.DB
}

func NewAdvertPostgres(db *sqlx.DB) *AdvertisingPostgres {
	return &AdvertisingPostgres{db: db}
}

func (r *AdvertisingPostgres) GetAll(limit int, page int, orderBy string) ([]model.GetAllAdvertsResponse, error) {
	var adverts []model.GetAllAdvertsResponse
	var query string
	offset := limit * (page - 1)

	switch orderBy {
	case "price_desc":
		query = fmt.Sprintf(
			"SELECT id, title, price, photos->>0 AS main_photo FROM %s ORDER BY price DESC LIMIT $1 OFFSET $2",
			advertsTable,
		)
	case "price_asc":
		query = fmt.Sprintf(
			"SELECT id, title, price, photos->>0 AS main_photo FROM %s ORDER BY price ASC LIMIT $1 OFFSET $2",
			advertsTable,
		)
	case "created_at_desc":
		query = fmt.Sprintf(
			"SELECT id, title, price, photos->>0 AS main_photo FROM %s ORDER BY created_at DESC LIMIT $1 OFFSET $2",
			advertsTable,
		)
	case "created_at_asc":
		query = fmt.Sprintf(
			"SELECT id, title, price, photos->>0 AS main_photo FROM %s ORDER BY created_at ASC LIMIT $1 OFFSET $2",
			advertsTable,
		)
	default:
		query = fmt.Sprintf(
			"SELECT id, title, price, photos->>0 AS main_photo FROM %s ORDER BY price ASC LIMIT $1 OFFSET $2",
			advertsTable,
		)
	}

	err := r.db.Select(&adverts, query, limit, offset)

	return adverts, err
}

func (r *AdvertisingPostgres) GetById(advertId int, fields []string) (model.Advert, error) {
	var advert model.Advert
	var query string

	selectedFields := "id, title, price, photos->>0 AS main_photo"

	for _, field := range fields {
		switch field {
		case "description":
			selectedFields += ", description"
		case "photos":
			selectedFields += ", photos"
		}
	}

	query = fmt.Sprintf(`SELECT %s FROM %s WHERE id = $1`, selectedFields, advertsTable)
	err := r.db.Get(&advert, query, advertId)

	return advert, err
}

func (r *AdvertisingPostgres) Create(advert model.CreateAdvertsBody) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var id int

	photos, err := json.Marshal(advert.Photos)
	if err != nil {
		return 0, err
	}

	createAdvertQuery := fmt.Sprintf(
		"INSERT INTO %s (title, description, photos, price) VALUES ($1, $2, $3, $4) RETURNING id",
		advertsTable,
	)
	row := tx.QueryRow(createAdvertQuery, advert.Title, advert.Description, photos, advert.Price)
	if err := row.Scan(&id); err != nil {
		err = tx.Rollback()
		
		if err != nil {
			return 0, err
		}

		return 0, err
	}

	return id, tx.Commit()
}
