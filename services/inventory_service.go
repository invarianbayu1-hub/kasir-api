package services

import (
	"kasir-api/models"
	"kasir-api/repositories"
)

type InventoryService struct {
	repo *repositories.InventoryRepository
}

func NewInventoryService(repo *repositories.InventoryRepository) *InventoryService {
	return &InventoryService{repo: repo}
}

func (s *InventoryService) GetAll() ([]models.Inventory, error) {
	return s.repo.GetAll()
}

func (s *InventoryService) Create(data *models.Inventory) error {
	return s.repo.Create(data)
}

func (s *InventoryService) GetByID(id int) (*models.Inventory, error) {
	return s.repo.GetByID(id)
}

func (s *InventoryService) Update(inventory *models.Inventory) error {
	return s.repo.Update(inventory)
}

func (s *InventoryService) Delete(id int) error {
	return s.repo.Delete(id)
}
