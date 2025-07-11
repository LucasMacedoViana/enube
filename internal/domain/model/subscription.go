package model

type Subscription struct {
    ID              uint   `gorm:"primaryKey"`
    SubscriptionID  string `gorm:"uniqueIndex"`
    Description     string
}