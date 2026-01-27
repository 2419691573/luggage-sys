package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"luggage-sys2/internal/database"
	"luggage-sys2/internal/models"
	"luggage-sys2/internal/utils"

	"gorm.io/gorm"
)

type LuggageService struct{}

func NewLuggageService() *LuggageService {
	return &LuggageService{}
}

type LuggageItem struct {
	StoreroomID  uint     `json:"storeroom_id" binding:"required"`
	Description  string   `json:"description"`
	Quantity     int      `json:"quantity"`
	SpecialNotes string   `json:"special_notes"`
	PhotoURLs    []string `json:"photo_urls"`
	PhotoURL     string   `json:"photo_url"`
}

type CreateLuggageRequest struct {
	GuestName    string        `json:"guest_name" binding:"required"`
	StaffName    string        `json:"staff_name"`
	ContactPhone string        `json:"contact_phone"`
	ContactEmail string        `json:"contact_email"`
	Description  string        `json:"description"`   // 单件模式
	Quantity     int           `json:"quantity"`      // 单件模式
	SpecialNotes string        `json:"special_notes"` // 单件模式
	PhotoURLs    []string      `json:"photo_urls"`    // 单件模式
	PhotoURL     string        `json:"photo_url"`     // 单件模式
	StoreroomID  uint          `json:"storeroom_id"`  // 单件模式（多件模式时不需要）
	Items        []LuggageItem `json:"items"`         // 多件模式
}

type UpdateLuggageRequest struct {
	GuestName    string   `json:"guest_name"`
	ContactPhone string   `json:"contact_phone"`
	Description  string   `json:"description"`
	SpecialNotes string   `json:"special_notes"`
	PhotoURLs    []string `json:"photo_urls"`
	PhotoURL     string   `json:"photo_url"`
}

// generateUniqueRetrievalCode 生成唯一的取件码（检查数据库中是否已存在）
// 如果提供了事务，则在事务内检查；否则使用普通连接检查
func (s *LuggageService) generateUniqueRetrievalCode(tx *gorm.DB) string {
	if tx == nil {
		tx = database.DB
	}
	maxAttempts := 200 // 增加尝试次数
	for i := 0; i < maxAttempts; i++ {
		code := utils.GenerateRetrievalCode()
		// 检查数据库中是否已存在该取件码
		// 使用普通查询（不使用事务隔离），确保能看到已提交的数据
		var count int64
		if tx == database.DB {
			// 如果是普通连接，直接查询
			database.DB.Model(&models.Luggage{}).Where("retrieval_code = ?", code).Count(&count)
		} else {
			// 如果是事务，也查询主连接以确保看到已提交的数据
			database.DB.Model(&models.Luggage{}).Where("retrieval_code = ?", code).Count(&count)
		}
		if count == 0 {
			return code
		}
	}
	// 如果200次都失败，生成一个更长的数字码（12位）
	code := utils.GenerateRetrievalCode() + utils.GenerateRetrievalCode()
	return code
}

// isDuplicateKeyError 检查是否是重复键错误
func (s *LuggageService) isDuplicateKeyError(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	// MySQL错误：Error 1062 (23000): Duplicate entry
	return strings.Contains(errStr, "Duplicate entry") || strings.Contains(errStr, "1062")
}

func (s *LuggageService) CreateLuggage(req CreateLuggageRequest, hotelID uint) (*models.Luggage, string, error) {
	// 判断是单件模式还是多件模式
	if len(req.Items) > 0 {
		// 多件模式：创建多个行李记录，共用同一个取件码
		// 使用事务确保原子性：要么全部成功，要么全部回滚
		tx := database.DB.Begin()
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
			}
		}()

		// 在事务内生成唯一的取件码（多件模式时，所有行李共用同一个取件码）
		retrievalCode := s.generateUniqueRetrievalCode(nil)

		var firstLuggage *models.Luggage
		for _, item := range req.Items {
			// 检查寄存室是否存在
			var storeroom models.Storeroom
			if err := tx.Where("id = ? AND hotel_id = ? AND is_active = ?", item.StoreroomID, hotelID, true).First(&storeroom).Error; err != nil {
				// 提供更详细的错误信息，帮助定位问题
				// 先检查寄存室是否存在（不检查 hotel_id 和 is_active）
				var checkStoreroom models.Storeroom
				if err2 := tx.Where("id = ?", item.StoreroomID).First(&checkStoreroom).Error; err2 != nil {
					tx.Rollback()
					return nil, "", fmt.Errorf("storeroom not found: storeroom_id %d does not exist", item.StoreroomID)
				}
				// 检查是否属于当前酒店
				if checkStoreroom.HotelID != hotelID {
					tx.Rollback()
					return nil, "", fmt.Errorf("storeroom not found: storeroom_id %d does not belong to your hotel (hotel_id: %d)", item.StoreroomID, hotelID)
				}
				// 检查是否启用
				if !checkStoreroom.IsActive {
					tx.Rollback()
					return nil, "", fmt.Errorf("storeroom not found: storeroom_id %d is not active", item.StoreroomID)
				}
				tx.Rollback()
				return nil, "", fmt.Errorf("storeroom not found: storeroom_id %d", item.StoreroomID)
			}

			// 检查容量
			var count int64
			tx.Model(&models.Luggage{}).Where("storeroom_id = ? AND status = ?", item.StoreroomID, "stored").Count(&count)
			if int(count) >= storeroom.Capacity {
				tx.Rollback()
				return nil, "", errors.New("storeroom is full")
			}

			// 处理图片
			photoURLs := item.PhotoURLs
			if len(photoURLs) == 0 && item.PhotoURL != "" {
				photoURLs = []string{item.PhotoURL}
			}
			photoURL := item.PhotoURL
			if photoURL == "" && len(photoURLs) > 0 {
				photoURL = photoURLs[0]
			}

			// 设置数量，如果未设置则默认为1
			quantity := item.Quantity
			if quantity <= 0 {
				quantity = 1
			}

			// 创建行李记录
			luggage := models.Luggage{
				GuestName:     req.GuestName,
				StaffName:     req.StaffName,
				ContactPhone:  req.ContactPhone,
				ContactEmail:  req.ContactEmail,
				Description:   item.Description,
				Quantity:      quantity,
				SpecialNotes:  item.SpecialNotes,
				PhotoURLs:     models.StringSlice(photoURLs),
				PhotoURL:      photoURL,
				StoreroomID:   item.StoreroomID,
				RetrievalCode: retrievalCode, // 共用同一个取件码
				Status:        "stored",
			}

			// 直接插入，多个行李可以共用同一个取件码
			if err := tx.Create(&luggage).Error; err != nil {
				tx.Rollback()
				return nil, "", err
			}

			// 保存第一个创建的行李记录用于返回
			if firstLuggage == nil {
				firstLuggage = &luggage
			}

			// 创建寄存记录
			storedLog := models.StoredLog{
				HotelID:   hotelID,
				LuggageID: luggage.ID,
				GuestName: luggage.GuestName,
				Status:    "stored",
			}
			if err := tx.Create(&storedLog).Error; err != nil {
				tx.Rollback()
				return nil, "", err
			}
		}

		// 提交事务
		if err := tx.Commit().Error; err != nil {
			return nil, "", err
		}

		return firstLuggage, retrievalCode, nil
	} else {
		// 单件模式：原有逻辑
		// 验证单件模式必填字段
		if req.StoreroomID == 0 {
			return nil, "", errors.New("storeroom_id is required in single-item mode")
		}

		// 检查寄存室是否存在
		var storeroom models.Storeroom
		if err := database.DB.Where("id = ? AND hotel_id = ? AND is_active = ?", req.StoreroomID, hotelID, true).First(&storeroom).Error; err != nil {
			return nil, "", errors.New("storeroom not found")
		}

		// 检查容量
		var count int64
		database.DB.Model(&models.Luggage{}).Where("storeroom_id = ? AND status = ?", req.StoreroomID, "stored").Count(&count)
		if int(count) >= storeroom.Capacity {
			return nil, "", errors.New("storeroom is full")
		}

		// 生成唯一的取件码（单件模式）
		retrievalCode := s.generateUniqueRetrievalCode(nil)

		// multi-photo compatibility:
		// - prefer photo_urls when provided
		// - fallback to photo_url (single) -> photo_urls=[photo_url]
		photoURLs := req.PhotoURLs
		if len(photoURLs) == 0 && req.PhotoURL != "" {
			photoURLs = []string{req.PhotoURL}
		}
		photoURL := req.PhotoURL
		if photoURL == "" && len(photoURLs) > 0 {
			photoURL = photoURLs[0]
		}

		// 设置数量，如果未设置则默认为1
		quantity := req.Quantity
		if quantity <= 0 {
			quantity = 1
		}

		luggage := models.Luggage{
			GuestName:     req.GuestName,
			StaffName:     req.StaffName,
			ContactPhone:  req.ContactPhone,
			ContactEmail:  req.ContactEmail,
			Description:   req.Description,
			Quantity:      quantity,
			SpecialNotes:  req.SpecialNotes,
			PhotoURLs:     models.StringSlice(photoURLs),
			PhotoURL:      photoURL,
			StoreroomID:   req.StoreroomID,
			RetrievalCode: retrievalCode,
			Status:        "stored",
		}

		if err := database.DB.Create(&luggage).Error; err != nil {
			return nil, "", err
		}

		// 创建寄存记录
		storedLog := models.StoredLog{
			HotelID:   hotelID,
			LuggageID: luggage.ID,
			GuestName: luggage.GuestName,
			Status:    "stored",
		}
		database.DB.Create(&storedLog)

		return &luggage, retrievalCode, nil
	}
}

func (s *LuggageService) GetLuggageByCode(code string, hotelID uint) ([]models.Luggage, error) {
	var luggages []models.Luggage
	if err := database.DB.Where("retrieval_code = ?", code).Find(&luggages).Error; err != nil {
		return nil, err
	}

	if len(luggages) == 0 {
		return nil, errors.New("luggage not found")
	}

	// 验证是否属于当前酒店（检查所有行李的寄存室）
	validLuggages := make([]models.Luggage, 0)
	for _, luggage := range luggages {
		var storeroom models.Storeroom
		if err := database.DB.Where("id = ? AND hotel_id = ?", luggage.StoreroomID, hotelID).First(&storeroom).Error; err == nil {
			validLuggages = append(validLuggages, luggage)
		}
	}

	if len(validLuggages) == 0 {
		return nil, errors.New("luggage not found in this hotel")
	}

	return validLuggages, nil
}

func (s *LuggageService) CheckoutLuggage(code string, username string, hotelID uint) ([]uint, error) {
	// 获取所有同取件码的行李记录
	var luggages []models.Luggage
	if err := database.DB.Where("retrieval_code = ?", code).Find(&luggages).Error; err != nil {
		return nil, errors.New("luggage not found")
	}

	if len(luggages) == 0 {
		return nil, errors.New("luggage not found")
	}

	// 验证是否属于当前酒店并取走所有在存状态的行李
	var retrievedIDs []uint
	now := time.Now()

	for i := range luggages {
		luggage := &luggages[i]

		// 验证是否属于当前酒店
		var storeroom models.Storeroom
		if err := database.DB.Where("id = ? AND hotel_id = ?", luggage.StoreroomID, hotelID).First(&storeroom).Error; err != nil {
			continue // 跳过不属于当前酒店的行李
		}

		// 只取走在存状态的行李
		if luggage.Status != "stored" {
			continue
		}

		// 更新状态
		luggage.Status = "retrieved"
		luggage.RetrievedAt = &now
		luggage.RetrievedBy = username

		if err := database.DB.Save(luggage).Error; err != nil {
			return nil, err
		}

		retrievedIDs = append(retrievedIDs, luggage.ID)

		// 创建取出记录
		retrievedLog := models.RetrievedLog{
			HotelID:     hotelID,
			LuggageID:   luggage.ID,
			GuestName:   luggage.GuestName,
			RetrievedBy: username,
		}
		database.DB.Create(&retrievedLog)
	}

	if len(retrievedIDs) == 0 {
		return nil, errors.New("no stored luggage found for this code in this hotel")
	}

	return retrievedIDs, nil
}

func (s *LuggageService) GetGuestList(hotelID uint) ([]string, error) {
	var guestNames []string
	if err := database.DB.Model(&models.Luggage{}).
		Distinct("guest_name").
		Joins("JOIN storerooms ON luggages.storeroom_id = storerooms.id").
		Where("storerooms.hotel_id = ? AND luggages.status = ?", hotelID, "stored").
		Pluck("guest_name", &guestNames).Error; err != nil {
		return nil, err
	}
	return guestNames, nil
}

func (s *LuggageService) GetLuggageByGuestName(guestName string, hotelID uint) ([]models.Luggage, error) {
	var luggages []models.Luggage
	if err := database.DB.
		Joins("JOIN storerooms ON luggages.storeroom_id = storerooms.id").
		Where("luggages.guest_name = ? AND storerooms.hotel_id = ? AND luggages.status = ?", guestName, hotelID, "stored").
		Find(&luggages).Error; err != nil {
		return nil, err
	}
	return luggages, nil
}

func (s *LuggageService) UpdateLuggage(id uint, req UpdateLuggageRequest, hotelID uint, username string) error {
	var luggage models.Luggage
	if err := database.DB.Where("id = ?", id).First(&luggage).Error; err != nil {
		return errors.New("invalid luggage id")
	}

	// 验证是否属于当前酒店
	var storeroom models.Storeroom
	if err := database.DB.Where("id = ? AND hotel_id = ?", luggage.StoreroomID, hotelID).First(&storeroom).Error; err != nil {
		return errors.New("luggage not found in this hotel")
	}

	// 保存旧数据用于日志
	oldData, _ := json.Marshal(luggage)

	// 更新字段
	if req.GuestName != "" {
		luggage.GuestName = req.GuestName
	}
	if req.ContactPhone != "" {
		luggage.ContactPhone = req.ContactPhone
	}
	if req.Description != "" {
		luggage.Description = req.Description
	}
	if req.SpecialNotes != "" {
		luggage.SpecialNotes = req.SpecialNotes
	}
	// photo_urls has higher priority (replace all)
	if len(req.PhotoURLs) > 0 {
		luggage.PhotoURLs = models.StringSlice(req.PhotoURLs)
		luggage.PhotoURL = req.PhotoURLs[0]
	} else if req.PhotoURL != "" {
		// backward compatibility: update single photo
		luggage.PhotoURL = req.PhotoURL
		if len(luggage.PhotoURLs) > 0 {
			luggage.PhotoURLs[0] = req.PhotoURL
		} else {
			luggage.PhotoURLs = models.StringSlice{req.PhotoURL}
		}
	}

	if err := database.DB.Save(&luggage).Error; err != nil {
		return err
	}

	// 创建修改记录
	newData, _ := json.Marshal(luggage)
	updatedLog := models.UpdatedLog{
		HotelID:   hotelID,
		LuggageID: luggage.ID,
		UpdatedBy: username,
		OldData:   string(oldData),
		NewData:   string(newData),
	}
	database.DB.Create(&updatedLog)

	return nil
}
