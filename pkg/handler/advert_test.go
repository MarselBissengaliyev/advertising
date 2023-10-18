package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MarselBissengaliyev/advertising/pkg/model"
	"github.com/MarselBissengaliyev/advertising/pkg/service"
	mock_service "github.com/MarselBissengaliyev/advertising/pkg/service/mocks"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHandler_createAdvert(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAdvert, advert model.CreateAdvertsBody)

	testTable := []struct {
		name                string
		inputBody           string
		inputAdvert         model.CreateAdvertsBody
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name: "OK",
			inputBody: `{
				"title": "Test",
				"description": "test",
				"photos": ["https://sun9-34.userapi.com/impg/3WOMrjs0H5io1nFdLaHv_NiOi5lxrz9qkk7RXg/-IRxizgUjRQ.jpg?size=960x1280&quality=95&sign=5f33c827a0d6a1ce884d82fe1202541d&type=album","https://sun9-68.userapi.com/impg/18OF1APOug-EIq6K63oIjxqR2wYN43DifTF-zw/PLhXQGZHl4c.jpg?size=1200x1600&quality=95&sign=94a5e71255e270fe46ba5d8f68f02770&type=album"],
				"price": 100.00
		}`,
			inputAdvert: model.CreateAdvertsBody{
				Title:       "Test",
				Description: "test",
				Photos: []string{
					"https://sun9-34.userapi.com/impg/3WOMrjs0H5io1nFdLaHv_NiOi5lxrz9qkk7RXg/-IRxizgUjRQ.jpg?size=960x1280&quality=95&sign=5f33c827a0d6a1ce884d82fe1202541d&type=album",
					"https://sun9-68.userapi.com/impg/18OF1APOug-EIq6K63oIjxqR2wYN43DifTF-zw/PLhXQGZHl4c.jpg?size=1200x1600&quality=95&sign=94a5e71255e270fe46ba5d8f68f02770&type=album",
				},
				Price: 100.00,
			},
			mockBehavior: func(s *mock_service.MockAdvert, advert model.CreateAdvertsBody) {
				s.EXPECT().Create(advert).Return(28, nil)
			},
			expectedStatusCode:  http.StatusCreated,
			expectedRequestBody: `{"id":28}` + "\n",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			// Init deps
			c := gomock.NewController(t)
			defer c.Finish()

			advert := mock_service.NewMockAdvert(c)
			testCase.mockBehavior(advert, testCase.inputAdvert)

			services := &service.Service{Advert: advert}
			handler := NewHandler(services)

			// Test server
			r := echo.New()
			r.POST("/adverts", handler.createAdvert)

			// Test Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest(
				"POST", "/adverts",
				bytes.NewBufferString(testCase.inputBody),
			)
			req.Header.Set("Content-Type", "application/json")

			// Perform Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}

func TestHandler_getAdvertById(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAdvert, advertId int, fields []string)

	testTable := []struct {
		name                string
		paramID             string
		queryFields         string
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:        "OK",
			paramID:     "1",
			queryFields: "title,price,description,photos",
			mockBehavior: func(s *mock_service.MockAdvert, advertId int, fields []string) {
				expectedFields := []string{"title", "price", "description", "photos"}
				assert.ElementsMatch(t, expectedFields, fields)

				advert := model.Advert{
					Id:          1,
					Title:       "Test Advert",
					Price:       100.00,
					Description: "This is a test advert",
					MainPhoto:   "https://example.com/main-photo.jpg",
					Photos:      `[ "https://example.com/photo1.jpg", "https://example.com/photo2.jpg" ]`,
				}

				s.EXPECT().GetById(advertId, fields).Return(advert, nil)
			},
			expectedStatusCode:  http.StatusOK,
			expectedRequestBody: `{"description":"This is a test advert","id":1,"main_photo":"https://example.com/main-photo.jpg","photos":["https://example.com/photo1.jpg","https://example.com/photo2.jpg"],"price":100,"title":"Test Advert"}` + "\n",
		},
		// Add more test cases as needed
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			// Init deps
			c := gomock.NewController(t)
			defer c.Finish()

			advert := mock_service.NewMockAdvert(c)
			testCase.mockBehavior(advert, 1, []string{"title", "price", "description", "photos"})

			services := &service.Service{Advert: advert}
			handler := NewHandler(services)

			// Test server
			r := echo.New()
			r.GET("/adverts/:id", handler.getAdvertById)

			// Test Request
			reqURL := "/adverts/" + testCase.paramID
			if testCase.queryFields != "" {
				reqURL += "?fields=" + testCase.queryFields
			}

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", reqURL, nil)

			// Perform Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}

func TestHandler_getAllAdverts(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAdvert, limit int, page int, orderBy string)

	testTable := []struct {
		name                string
		queryParams         string
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:        "OK",
			queryParams: "sort=price_desc&page=1",
			mockBehavior: func(s *mock_service.MockAdvert, limit int, page int, orderBy string) {
				assert.Equal(t, 10, limit)
				assert.Equal(t, 1, page)
				assert.Equal(t, "price_desc", orderBy)

				adverts := []model.GetAllAdvertsResponse{
					{
						Id:        1,
						Title:     "Test Advert 1",
						Price:     100.00,
						MainPhoto: "https://example.com/main-photo-1.jpg",
					},
					{
						Id:        2,
						Title:     "Test Advert 2",
						Price:     150.00,
						MainPhoto: "https://example.com/main-photo-2.jpg",
					},
				}

				s.EXPECT().GetAll(limit, page, orderBy).Return(adverts, nil)
			},
			expectedStatusCode:  http.StatusOK,
			expectedRequestBody: `{"data":[{"id":1,"title":"Test Advert 1","main_photo":"https://example.com/main-photo-1.jpg","price":100},{"id":2,"title":"Test Advert 2","main_photo":"https://example.com/main-photo-2.jpg","price":150}]}` + "\n",
		},
		// Add more test cases as needed
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			// Init deps
			c := gomock.NewController(t)
			defer c.Finish()

			advert := mock_service.NewMockAdvert(c)
			testCase.mockBehavior(advert, 10, 1, "price_desc")

			services := &service.Service{Advert: advert}
			handler := NewHandler(services)

			// Test server
			r := echo.New()
			r.GET("/adverts", handler.getAllAdverts)

			// Test Request
			reqURL := "/adverts"
			if testCase.queryParams != "" {
				reqURL += "?" + testCase.queryParams
			}

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", reqURL, nil)

			// Perform Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}



