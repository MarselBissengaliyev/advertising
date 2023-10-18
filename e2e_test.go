package e2e_test

import (
	"bytes"
	"net/http"
	"os"
	"testing"

	"github.com/MarselBissengaliyev/advertising/pkg/handler"
	"github.com/MarselBissengaliyev/advertising/pkg/model"
	"github.com/MarselBissengaliyev/advertising/pkg/repository"
	"github.com/MarselBissengaliyev/advertising/pkg/service"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/spf13/viper"
)

func TestHandler_E2E(t *testing.T) {
	Convey("Should init server and connect to test database", t, func() {
		logrus.SetFormatter(new(logrus.JSONFormatter))

		err := initConfig()
		Convey("Config should be initialized", func() {
			So(err, ShouldEqual, nil)
		})

		err = godotenv.Load("test/.env")
		Convey("Env should be loaded", func() {
			So(err, ShouldEqual, nil)
		})

		db, err := repository.NewPostgresDB(repository.Config{
			Host:     viper.GetString("db.host"),
			Port:     viper.GetString("db.port"),
			UserName: viper.GetString("db.username"),
			DBName:   viper.GetString("db.dbname"),
			SSLMode:  viper.GetString("db.sslmode"),
			Password: os.Getenv("DB_PASSWORD"),
		})
		Convey("New db should not be nil", func() {
			So(db, ShouldNotBeNil)
		})
		Convey("New db err should be nil", func() {
			So(err, ShouldEqual, nil)
		})

		repos := repository.NewRepository(db)
		services := service.NewService(repos)
		handlers := handler.NewHandler(services)
		srv := new(model.Server)
		go func() {
			port := viper.GetString("port")
			err = srv.Run(port, handlers.InitRoutes())
			Convey("Server should not be nil", t, func() {
				So(srv, ShouldNotBeNil)
			})
		}()

		getAllAdvertsResponse, err := http.Get("http://localhost:8080/api/adverts/")
		So(err, ShouldEqual, nil)
		So(getAllAdvertsResponse.StatusCode, ShouldEqual, http.StatusOK)

		json_data := bytes.NewBufferString(`{
			"title": "Test",
			"description": "test",
			"photos": ["https://sun9-34.userapi.com/impg/3WOMrjs0H5io1nFdLaHv_NiOi5lxrz9qkk7RXg/-IRxizgUjRQ.jpg?size=960x1280&quality=95&sign=5f33c827a0d6a1ce884d82fe1202541d&type=album","https://sun9-68.userapi.com/impg/18OF1APOug-EIq6K63oIjxqR2wYN43DifTF-zw/PLhXQGZHl4c.jpg?size=1200x1600&quality=95&sign=94a5e71255e270fe46ba5d8f68f02770&type=album"],
			"price": 100.00
		}`)

		createAdvertsResponse, err := http.Post("http://localhost:8080/api/adverts/", "application/json", json_data)
		So(err, ShouldEqual, nil)
		So(createAdvertsResponse.StatusCode, ShouldEqual, http.StatusCreated)

		getAdvertByIdResponse, err := http.Get("http://localhost:8080/api/adverts/1")
		So(err, ShouldEqual, nil)
		So(getAdvertByIdResponse.StatusCode, ShouldEqual, http.StatusOK)

		logrus.Print("Advertising App Started")
	})
}

func initConfig() error {
	viper.AddConfigPath("test")
	viper.SetConfigType("yml")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
