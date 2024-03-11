package model

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Name          string  `gorm:"type:varchar(160);not null;uniqueIndex"`                             // "Shoe #1", "Shoe #2", ...
	Description   string  `gorm:"type:text;not null"`                                                 // "A very nice shoe", ...
	Files         []File  `gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"` // 1 to many; downloadables, ...
	Images        []File  `gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"` // 1 to many; images for the product, ...
	PriceAmount   float64 `gorm:"type:decimal(10,2);not null"`                                        // 100.00
	PriceCurrency string  `gorm:"type:varchar(3);not null"`                                           // usd, eur, ...
}

// make sure to put this behind role permissions so only admins can create plans
func (model *Product) Create(database *gorm.DB) error {
	return database.Create(model).Error
}

func (model *Product) ReadAll(database *gorm.DB) ([]Product, error) {
	var plans []Product
	err := database.Find(&plans).Error
	return plans, err
}

func (model *Product) Read(database *gorm.DB) error {
	return database.First(model, model.ID).Error
}

// make sure to put this behind role permissions so only admins can update plans
func (model *Product) Update(database *gorm.DB) error {
	return database.Save(model).Error
}

// make sure to put this behind role permissions so only admins can delete plans
func (model *Product) Delete(database *gorm.DB) error {
	return database.Delete(model).Error
}

type File struct {
	gorm.Model
	Name      string `gorm:"type:varchar(160);not null;uniqueIndex"` // "shoe1.jpg", "shoe2.jpg", ...
	FileBlob  []byte `gorm:"type:blob;not null"`
	Type      string `gorm:"type:varchar(100);not null"` // "image/jpeg", "application/pdf", ...
	Extension string `gorm:"type:varchar(10);not null"`  // "jpg", "png", "pdf", ...
}

// make sure to put this behind role permissions so only admins can create files
func (model *File) Create(database *gorm.DB) error {
	return database.Create(model).Error
}

func (model *File) ReadAll(database *gorm.DB) ([]File, error) {
	var files []File
	err := database.Find(&files).Error
	return files, err
}

func (model *File) Read(database *gorm.DB) error {
	return database.First(model, model.ID).Error
}

// make sure to put this behind role permissions so only admins can update files
func (model *File) Update(database *gorm.DB) error {
	return database.Save(model).Error
}

// make sure to put this behind role permissions so only admins can delete files
func (model *File) Delete(database *gorm.DB) error {
	return database.Delete(model).Error
}

type OrderStatus string

const (
	Waiting       OrderStatus = "waiting"
	Confirming    OrderStatus = "confirming"
	Confirmed     OrderStatus = "confirmed"
	Sending       OrderStatus = "sending"
	PartiallyPaid OrderStatus = "partially_paid"
	Finished      OrderStatus = "finished"
	Failed        OrderStatus = "failed"
	Expired       OrderStatus = "expired"
)

type Order struct {
	gorm.Model
	ExternalId    string      `json:"external_id" gorm:"unique"` // from NOWPayments
	UserId        uint        `json:"user_id"`                   // the user who made the payment
	Status        OrderStatus `json:"status"`                    // waiting, confirming, confirmed, sending, partially_paid, finished, failed, expired
	ProductId     uint        `json:"product_id"`                // the product the user bought
	PriceAmount   float64     `json:"price_amount"`              // the amount the user paid (it's important to store because prices can change)
	PriceCurrency string      `json:"price_currency"`            // the currency in which we denominate the amount (e.g. usd)
	InvoiceUrl    string      `json:"invoice_url" gorm:"unique"` // from NOWPayments
}

func (model *Order) Create(database *gorm.DB) error {
	return database.Create(model).Error
}

func (model *Order) ReadAll(database *gorm.DB) ([]Order, error) {
	var orders []Order
	err := database.Find(&orders).Error
	return orders, err
}

func (model *Order) Read(database *gorm.DB) error {
	return database.First(model, model.ID).Error
}

func (model *Order) Update(database *gorm.DB) error {
	return database.Save(model).Error
}

func (model *Order) Delete(database *gorm.DB) error {
	return database.Delete(model).Error
}
