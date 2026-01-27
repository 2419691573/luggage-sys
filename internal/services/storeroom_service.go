package services

import (
	"errors"

	"luggage-sys2/internal/database"
	"luggage-sys2/internal/models"
)

type StoreroomService struct{}

func NewStoreroomService() *StoreroomService {
	return &StoreroomService{}
}

type CreateStoreroomRequest struct {
	Name     string `json:"name" binding:"required"`
	Location string `json:"location"`
	Capacity int    `json:"capacity" binding:"required"`
	IsActive bool   `json:"is_active"`
}

type UpdateStoreroomRequest struct {
	IsActive bool `json:"is_active"`
}

func (s *StoreroomService) ListStorerooms(hotelID uint) ([]models.Storeroom, error) {
	var storerooms []models.Storeroom
	if err := database.DB.Where("hotel_id = ?", hotelID).Find(&storerooms).Error; err != nil {
		return nil, err
	}

	// 计算每个寄存室的已存数量和剩余容量
	for i := range storerooms {
		var count int64
		database.DB.Model(&models.Luggage{}).
			Where("storeroom_id = ? AND status = ?", storerooms[i].ID, "stored").
			Count(&count)
		storerooms[i].StoredCount = int(count)
		storerooms[i].RemainingCapacity = storerooms[i].Capacity - int(count)
	}

	return storerooms, nil
}

func (s *StoreroomService) CreateStoreroom(req CreateStoreroomRequest, hotelID uint) (*models.Storeroom, error) {
	storeroom := models.Storeroom{
		HotelID:  hotelID,
		Name:     req.Name,
		Location: req.Location,
		Capacity: req.Capacity,
		IsActive: req.IsActive,
	}

	if err := database.DB.Create(&storeroom).Error; err != nil {
		return nil, errors.New("create storeroom failed")
	}

	return &storeroom, nil
}

func (s *StoreroomService) UpdateStoreroom(id uint, req UpdateStoreroomRequest, hotelID uint) error {
	var storeroom models.Storeroom
	if err := database.DB.Where("id = ? AND hotel_id = ?", id, hotelID).First(&storeroom).Error; err != nil {
		return errors.New("invalid storeroom id")
	}

	storeroom.IsActive = req.IsActive
	if err := database.DB.Save(&storeroom).Error; err != nil {
		return errors.New("update storeroom status failed")
	}

	return nil
}

func (s *StoreroomService) GetStoreroomOrders(id uint, hotelID uint, status string) ([]models.Luggage, error) {
	// 验证寄存室是否属于当前酒店
	var storeroom models.Storeroom
	if err := database.DB.Where("id = ? AND hotel_id = ?", id, hotelID).First(&storeroom).Error; err != nil {
		return nil, errors.New("invalid storeroom id")
	}

	var luggages []models.Luggage
	query := database.DB.Where("storeroom_id = ?", id)
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if err := query.Find(&luggages).Error; err != nil {
		return nil, err
	}

	return luggages, nil
}
