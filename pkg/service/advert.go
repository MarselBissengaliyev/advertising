package service

import (
	"github.com/MarselBissengaliyev/advertising/pkg/model"
	"github.com/MarselBissengaliyev/advertising/pkg/repository"
)

type AdvertService struct {
	repo repository.Advert
}

func NewAdvertService(repo repository.Advert) *AdvertService {
	return &AdvertService{repo: repo}
}

func (s *AdvertService) GetAll(limit int, page int, orderBy string) ([]model.GetAllAdvertsResponse, error) {
	return s.repo.GetAll(limit, page, orderBy)
}

func (s *AdvertService) GetById(advertId int, fields []string) (model.Advert, error) {
	return s.repo.GetById(advertId, fields)
}

func (s *AdvertService) Create(advert model.CreateAdvertsBody) (int, error) {
	return s.repo.Create(advert)
}
