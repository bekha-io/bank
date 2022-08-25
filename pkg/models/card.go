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
	AccountID  string           `gorm:"notNull" json:"accountId"`
	PAN        types.PAN        `gorm:"primaryKey;size:16;<-:create;notNull" json:"pan"`
	CardSystem types.CardSystem `gorm:"notNull" json:"cardSystem"`
	ExpireDate string           `gorm:"size:5;notNull" json:"expireDate"`
	PIN        types.PIN        `gorm:"size:4;notNull" json:"pin,omitempty"`
	CV2        types.CV2        `gorm:"size:3;notNull;default:123" json:"cv2,omitempty"`
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

func generatePIN() types.PIN {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, 4)
	for i := range b {
		b[i] = types.Digits[rand.Intn(len(types.Digits))]
	}
	return types.PIN(b)
}

func generateCV2() types.CV2 {
	rand.Seed(time.Now().UnixMilli())
	b := make([]byte, 3)
	for i := range b {
		b[i] = types.Digits[rand.Intn(len(types.Digits))]
	}
	return types.CV2(b)
}

func generateExpireDate(years int) string {
	n := time.Now().AddDate(years, 0, 0)
	dateAsString := n.Format("01/06")
	return dateAsString
}

func (c *Card) BeforeCreate(tx *gorm.DB) (err error) {
	c.PAN = generatePAN()
	c.PIN = generatePIN()
	c.CV2 = generateCV2()
	c.ExpireDate = generateExpireDate(3)

	if c.CardSystem == "" {
		c.CardSystem = types.KortiMilli
	}

	return
}

func (c *Card) ExpireAsDate() time.Time {
	t, err := time.Parse("01/06", c.ExpireDate)
	if err != nil {
		log.Fatal(err)
	}
	return t
}

func (c *Card) MaskSensitive() {
	c.PIN = ""
	c.CV2 = ""
}
