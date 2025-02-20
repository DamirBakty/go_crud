package services

import (
	"crud/market-service/models"
	"crud/market-service/repos"
)

type MarketService struct {
	repo *repos.MarketRepo
}

func NewMarketService(repo *repos.MarketRepo) *MarketService {
	return &MarketService{repo: repo}
}

func (s *MarketService) CreateMarket(market models.MarketEdit) (int, error) {
	return s.repo.Create(market)
}

func (s *MarketService) GetMarket(id int) (*models.MarketView, error) {
	return s.repo.Get(id)
}

func (s *MarketService) UpdateMarket(id int, market models.MarketEdit) error {
	return s.repo.Update(id, market)
}

func (s *MarketService) DeleteMarket(id int) error {
	return s.repo.Delete(id)
}
