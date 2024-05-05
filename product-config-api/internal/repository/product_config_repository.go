package repository

import (
	"errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"product-config-api/internal/config"
	"product-config-api/internal/model"
	"product-config-api/internal/model/request"
)

type ProductConfigRepository interface {
	CreateProductConfig(requestModel request.ProductConfigRequestModel) error
	GetActiveProductConfig() (model.ProductConfigEntity, error)
}

type ProductConfigRepositoryImpl struct {
	applicationConfig config.ApplicationConfigManager
	db                *gorm.DB
}

func (p *ProductConfigRepositoryImpl) GetActiveProductConfig() (model.ProductConfigEntity, error) {
	var productConfigEntity model.ProductConfigEntity
	if err := p.db.Model(&model.ProductConfigEntity{}).Where("active = ?", true).First(&productConfigEntity).Error; err != nil {
		return model.ProductConfigEntity{}, err
	}
	return productConfigEntity, nil
}

func (p *ProductConfigRepositoryImpl) CreateProductConfig(requestModel request.ProductConfigRequestModel) error {
	productConfigEntity := &model.ProductConfigEntity{
		FieldName: requestModel.FieldName,
		OrderType: requestModel.OrderType,
		Active:    true,
	}

	transaction := p.db.Begin()

	var existingProductConfigEntity model.ProductConfigEntity

	if err := transaction.Model(&model.ProductConfigEntity{}).Where("active = ?", true).First(&existingProductConfigEntity).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			transaction.Rollback()
			return err
		}
	}

	if existingProductConfigEntity.ID != 0 {
		existingProductConfigEntity.Active = false
		if err := transaction.Save(&existingProductConfigEntity).Error; err != nil {
			transaction.Rollback()
			return err
		}
	}

	if err := transaction.Create(&productConfigEntity).Error; err != nil {
		transaction.Rollback()
		return err
	}

	if err := transaction.Commit().Error; err != nil {
		return err
	}
	return nil
}

func NewProductConfigRepository(applicationConfig config.ApplicationConfigManager) ProductConfigRepository {
	db, err := gorm.Open(postgres.Open(applicationConfig.PostgreSql.Url), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&model.ProductConfigEntity{})
	if err != nil {
		panic(err)
	}

	return &ProductConfigRepositoryImpl{
		applicationConfig: applicationConfig,
		db:                db,
	}
}
