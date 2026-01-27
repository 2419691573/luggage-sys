package services

import (
	"luggage-sys2/internal/database"
	"luggage-sys2/internal/models"
	"luggage-sys2/internal/utils"
)

type AuthService struct{}

func NewAuthService() *AuthService {
	return &AuthService{}
}

func (s *AuthService) Login(username, password string) (*models.User, string, error) {
	var user models.User
	if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, "", err
	}

	if !utils.CheckPassword(password, user.Password) {
		return nil, "", nil
	}

	token, err := utils.GenerateToken(user.ID, user.Username, user.HotelID)
	if err != nil {
		return nil, "", err
	}

	return &user, token, nil
}
