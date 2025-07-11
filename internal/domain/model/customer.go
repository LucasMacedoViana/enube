package model

type Customer struct {
    ID           uint   `gorm:"primaryKey"`
    CustomerID   string `gorm:"uniqueIndex"`
    Name         string
    DomainName   string
    Country      string
}