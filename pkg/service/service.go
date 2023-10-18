package service

import (
	"github.com/MarselBissengaliyev/advertising/pkg/model"
	"github.com/MarselBissengaliyev/advertising/pkg/repository"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type Advert interface {
	GetAll(limit int, page int, orderBy string) ([]model.GetAllAdvertsResponse, error)
	GetById(advertId int, fields []string) (model.Advert, error)
	Create(advert model.CreateAdvertsBody) (int, error)
}

type Service struct {
	Advert
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Advert: NewAdvertService(repos.Advert),
	}
}
