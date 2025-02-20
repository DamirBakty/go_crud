package services

import (
	"crud/market-service/models"
	"crud/market-service/repos"
	"errors"
)

type ItemService struct {
	repo *repos.ItemRepo
}

func NewItemService(repo *repos.ItemRepo) *ItemService {
	return &ItemService{repo: repo}
}

func (s *ItemService) CreateItem(item models.Item) (int, error) {
	// Validate input
	if item.Name == "" {
		return 0, errors.New("item name cannot be empty")
	}
	if item.Count < 0 {
		return 0, errors.New("item count cannot be negative")
	}
	if item.Price < 0 {
		return 0, errors.New("item price cannot be negative")
	}

	return s.repo.Create(item)
}

func (s *ItemService) GetItem(id int) (*models.Item, error) {
	if id <= 0 {
		return nil, errors.New("invalid item id")
	}

	item, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, errors.New("item not found")
	}

	return item, nil
}

func (s *ItemService) GetAllItems() ([]models.Item, error) {
	return s.repo.GetAll()
}

func (s *ItemService) UpdateItem(id int, item models.ItemEdit) error {
	// Validate input
	if id <= 0 {
		return errors.New("invalid item id")
	}
	if item.Name == "" {
		return errors.New("item name cannot be empty")
	}
	if item.Count < 0 {
		return errors.New("item count cannot be negative")
	}
	if item.Price < 0 {
		return errors.New("item price cannot be negative")
	}

	// Check if item exists
	existing, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("item not found")
	}

	return s.repo.Update(id, item)
}

func (s *ItemService) DeleteItem(id int) error {
	if id <= 0 {
		return errors.New("invalid item id")
	}

	err := s.repo.Delete(id)
	if err == nil {
		return nil
	}

	return err
}

func (s *ItemService) GetItemsByMarketID(marketID int) ([]models.Item, error) {
	if marketID <= 0 {
		return nil, errors.New("invalid market id")
	}

	items, err := s.repo.GetByMarketID(marketID)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (s *ItemService) UpdateItemCount(id int, newCount int) error {
	if id <= 0 {
		return errors.New("invalid item id")
	}
	if newCount < 0 {
		return errors.New("item count cannot be negative")
	}

	// Check if item exists
	existing, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("item not found")
	}

	itemEdit := models.ItemEdit{
		Name:  existing.Name,
		Count: newCount,
		Price: existing.Price,
	}

	return s.repo.Update(id, itemEdit)
}

func (s *ItemService) ValidateItemsExist(itemIDs []int) error {
	for _, id := range itemIDs {
		item, err := s.repo.GetByID(id)
		if err != nil {
			return err
		}
		if item == nil {
			return errors.New("one or more items not found")
		}
	}
	return nil
}
