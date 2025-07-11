package model

import "time"

type BillingItem struct {
    ID                              uint      `gorm:"primaryKey"`
    InvoiceNumber                   string
    PartnerID                       uint
    CustomerID                      uint
    SubscriptionID                  uint
    ProductID                       uint
    MeterID                         uint
    ChargeStartDate                 time.Time
    ChargeEndDate                   time.Time
    UsageDate                       string
    ResourceLocation                string
    ResourceGroup                   string
    ResourceURI                     string
    ConsumedService                 string
    ChargeType                      string
    Unit                            string
    UnitPrice                       float64
    Quantity                        float64
    UnitType                        string
    BillingPreTaxTotal              float64
    BillingCurrency                 string
    PricingPreTaxTotal              float64
    PricingCurrency                 string
    EffectiveUnitPrice              float64
    PCToBCExchangeRate              float64
    PCToBCExchangeRateDate          time.Time
    EntitlementID                   string
    EntitlementDescription          string
    PartnerEarnedCreditPercentage   int
    CreditPercentage                int
    CreditType                      string
    Tags                            string
    AdditionalInfo                  string
    ServiceInfo1                    string
    ServiceInfo2                    string
}