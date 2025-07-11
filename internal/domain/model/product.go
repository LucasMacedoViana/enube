package model

type Product struct {
    ID              uint   `gorm:"primaryKey"`
    ProductID       string
    SkuID           string
    AvailabilityID  string
    SkuName         string
    ProductName     string
    PublisherName   string
    PublisherID     string
}