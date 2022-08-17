package models

import (
	"banking/pkg/types"
	"gorm.io/gorm"
	"log"
	"math/rand"
	"time"
)

type Card struct {
	BaseModel
	AccountID  string           `gorm:"notNull"`
	PAN        types.PAN        `gorm:"primaryKey;size:16;<-:create;notNull"`
	CardSystem types.CardSystem `gorm:"notNull"`
	expireDate string           `gorm:"size:5;notNull"`
}

func generatePAN() types.PAN {
	rand.Seed(time.Now().UnixNano())
	s := "44448888"
	b := make([]byte, 8)
	for i := range b {
		b[i] = types.Digits[rand.Intn(len(types.Digits))]
	}
	return types.PAN(s + string(b))
}

func generateExpireDate(years int) string {
	n := time.Now().AddDate(years, 0, 0)
	dateAsString := n.Format("01/06")
	return dateAsString
}

func (c *Card) BeforeCreate(tx *gorm.DB) (err error) {
	c.PAN = generatePAN()
	c.CardSystem = types.MasterCard
	c.expireDate = generateExpireDate(3)
	return
}

func (c *Card) ExpireDate() time.Time {
	t, err := time.Parse("01/06", c.expireDate)
	if err != nil {
		log.Fatal(err)
	}
	return t
}

func (c *Card) ExpireDateString() string {
	return c.expireDate
}
