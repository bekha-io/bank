package models

//
//type Merchant struct {
//	BaseModel
//	ID               uint    `gorm:"primaryKey"`
//	AccountID        string  `gorm:"notNull"`
//	Name             string  `gorm:"notNull"`
//	IsVerified       bool    `gorm:"default:false"` // Статус верификации мерчанта
//	BankInterest     float64 `gorm:"default:1"`     // Комиссия банка за каждый обработанный платеж мерчанта (в %)
//	CashbackInterest float64 `gorm:"default:2"`     // Процентная ставка кэшбека, выплачиваемая клиенту на счет в виде бонусов
//}
//
//func (m Merchant) GetAccounts() *[]Account {
//	var a []Account
//	Db.Find(&a, " = ?", m.ID)
//	return &a
//}
