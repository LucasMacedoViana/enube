package model

type User struct {
	ID       uint   `gorm:"primaryKey" json:"-"`
	Name     string `gorm:"uniqueIndex;not null" json:"name"`
	Password string `gorm:"not null" json:"-"`
}
