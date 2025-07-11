package model

type Meter struct {
    ID          uint   `gorm:"primaryKey"`
    MeterID     string
    Type        string
    Category    string
    SubCategory string
    Name        string
    Region      string
}