package services

import (
	"luggage-sys2/internal/database"
	"luggage-sys2/internal/models"
)

type LogService struct{}

func NewLogService() *LogService {
	return &LogService{}
}

func (s *LogService) GetStoredLogs(hotelID uint) ([]models.StoredLog, error) {
	var logs []models.StoredLog
	if err := database.DB.Where("hotel_id = ?", hotelID).Order("stored_at DESC").Find(&logs).Error; err != nil {
		return nil, err
	}
	return logs, nil
}

func (s *LogService) GetUpdatedLogs(hotelID uint) ([]models.UpdatedLog, error) {
	var logs []models.UpdatedLog
	if err := database.DB.Where("hotel_id = ?", hotelID).Order("updated_at DESC").Find(&logs).Error; err != nil {
		return nil, err
	}
	return logs, nil
}

func (s *LogService) GetRetrievedLogs(hotelID uint) ([]models.RetrievedLog, error) {
	var logs []models.RetrievedLog
	if err := database.DB.Where("hotel_id = ?", hotelID).Order("retrieved_at DESC").Find(&logs).Error; err != nil {
		return nil, err
	}
	return logs, nil
}
