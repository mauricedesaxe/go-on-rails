package model

import "gorm.io/gorm"

type ProductModel struct {
	gorm.Model
	Name          string      `gorm:"type:varchar(160);not null;uniqueIndex"`                             // "Shoe #1", "Shoe #2", ...
	Description   string      `gorm:"type:text;not null"`                                                 // "A very nice shoe", ...
	Files         []FileModel `gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"` // 1 to many; downloadables, ...
	Images        []FileModel `gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"` // 1 to many; images for the product, ...
	PriceAmount   float64     `gorm:"type:decimal(10,2);not null"`                                        // 100.00
	PriceCurrency string      `gorm:"type:varchar(3);not null"`                                           // usd, eur, ...
}

// make sure to put this behind role permissions so only admins can create plans
func (model *ProductModel) Create(database *gorm.DB) error {
	return database.Create(model).Error
}

func (model *ProductModel) ReadAll(database *gorm.DB) ([]ProductModel, error) {
	var plans []ProductModel
	err := database.Find(&plans).Error
	return plans, err
}

func (model *ProductModel) Read(database *gorm.DB) error {
	return database.First(model, model.ID).Error
}

// make sure to put this behind role permissions so only admins can update plans
func (model *ProductModel) Update(database *gorm.DB) error {
	return database.Save(model).Error
}

// make sure to put this behind role permissions so only admins can delete plans
func (model *ProductModel) Delete(database *gorm.DB) error {
	return database.Delete(model).Error
}

type FileModel struct {
	gorm.Model
	Name      string `gorm:"type:varchar(160);not null;uniqueIndex"` // "shoe1.jpg", "shoe2.jpg", ...
	FileBlob  []byte `gorm:"type:blob;not null"`
	Type      string `gorm:"type:varchar(100);not null"` // "image/jpeg", "application/pdf", ...
	Extension string `gorm:"type:varchar(10);not null"`  // "jpg", "png", "pdf", ...
}

// make sure to put this behind role permissions so only admins can create files
func (model *FileModel) Create(database *gorm.DB) error {
	return database.Create(model).Error
}

func (model *FileModel) ReadAll(database *gorm.DB) ([]FileModel, error) {
	var files []FileModel
	err := database.Find(&files).Error
	return files, err
}

func (model *FileModel) Read(database *gorm.DB) error {
	return database.First(model, model.ID).Error
}

// make sure to put this behind role permissions so only admins can update files
func (model *FileModel) Update(database *gorm.DB) error {
	return database.Save(model).Error
}

// make sure to put this behind role permissions so only admins can delete files
func (model *FileModel) Delete(database *gorm.DB) error {
	return database.Delete(model).Error
}
