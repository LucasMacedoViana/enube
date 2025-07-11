package model

type Partner struct {
    ID        uint   `gorm:"primaryKey"`
    PartnerID string `gorm:"uniqueIndex"`
    Name      string
}