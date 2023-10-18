package repository

import (
	"encoding/json"
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/MarselBissengaliyev/advertising/pkg/model"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAdvertPostgres_GetAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	advertisingRepo := NewAdvertPostgres(sqlx.NewDb(db, "sqlmock"))

	testCases := []struct {
		name           string
		limit, page    int
		orderBy        string
		mockExpect     func()
		expectedResult []model.GetAllAdvertsResponse
		expectedError  error
	}{
		{
			name:    "Success",
			limit:   10,
			page:    1,
			orderBy: "price_desc",
			mockExpect: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "price", "main_photo"}).
					AddRow(1, "Advert1", 100.0, "photo1").
					AddRow(2, "Advert2", 150.0, "photo2")

				query := fmt.Sprintf(
					"SELECT id, title, price, photos->>0 AS main_photo FROM %s ORDER BY price DESC LIMIT $1 OFFSET $2",
					advertsTable,
				)

				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(10, 0).
					WillReturnRows(rows)
			},
			expectedResult: []model.GetAllAdvertsResponse{
				{Id: 1, Title: "Advert1", Price: 100.0, MainPhoto: "photo1"},
				{Id: 2, Title: "Advert2", Price: 150.0, MainPhoto: "photo2"},
			},
			expectedError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockExpect()

			result, err := advertisingRepo.GetAll(tc.limit, tc.page, tc.orderBy)

			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedResult, result)
		})
	}
}

func TestAdvertPostgres_GetById(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	advertisingRepo := NewAdvertPostgres(sqlx.NewDb(db, "sqlmock"))
	testCases := []struct {
		name           string
		fields         []string
		advertId       int
		mockExpect     func()
		expectedResult model.Advert
		expectedError  error
	}{
		{
			name:     "Success",
			fields:   []string{"description", "photos"},
			advertId: 1,
			mockExpect: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "price", "main_photo", "description", "photos"}).
					AddRow(1, "Advert1", 100.0, "photo1", "description1", `["photo1.jpg", "photo2.jpg"]`).
					AddRow(2, "Advert2", 150.0, "photo2", "description2", `["photo1.jpg", "photo2.jpg"]`)

				selectedFields := "id, title, price, photos->>0 AS main_photo"

				for _, field := range []string{"description", "photos"} {
					switch field {
					case "description":
						selectedFields += ", description"
					case "photos":
						selectedFields += ", photos"
					}
				}

				query := fmt.Sprintf("SELECT %s FROM %s WHERE id = $1", selectedFields, advertsTable)

				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(1).
					WillReturnRows(rows)
			},
			expectedResult: model.Advert{
				Id: 1, Title: "Advert1",
				Price: 100.0, MainPhoto: "photo1",
				Description: "description1", Photos: `["photo1.jpg", "photo2.jpg"]`,
			},
			expectedError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockExpect()

			result, err := advertisingRepo.GetById(tc.advertId, tc.fields)

			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedResult, result)
		})
	}
}

func TestAdvertPostgres_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer db.Close()

	advertisingRepo := NewAdvertPostgres(sqlx.NewDb(db, "sqlmock"))
	testCases := []struct {
		name           string
		advert         model.CreateAdvertsBody
		mockExpect     func()
		expectedResult int
		expectedError  error
	}{
		{
			name: "Success",
			advert: model.CreateAdvertsBody{
				Title:       "Test",
				Description: "Test",
				Photos:      []string{"photo1.jpg", "photo2.jpg"},
				Price:       200.00,
			},
			mockExpect: func() {
				mock.ExpectBegin()
				createAdvertQuery := fmt.Sprintf(
					"INSERT INTO %s (title, description, photos, price) VALUES ($1, $2, $3, $4) RETURNING id",
					advertsTable,
				)

				photos, err := json.Marshal([]string{"photo1.jpg", "photo2.jpg"})
				if err != nil {
					require.NoError(t, err)
				}

				mock.ExpectQuery(regexp.QuoteMeta(createAdvertQuery)).WithArgs("Test", "Test", photos, 200.00).WillReturnRows(
					sqlmock.NewRows([]string{"id"}).AddRow(1))
				mock.ExpectCommit()
			},
			expectedResult: 1,
			expectedError:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockExpect()

			result, err := advertisingRepo.Create(tc.advert)

			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedResult, result)
		})
	}
}
