package benchmark

import (
	"banking/internal/http/json/respModels"
	"banking/pkg/models"
	"testing"
)

func BenchmarkViaStruct(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var card = models.Card{
			AccountID:  "123123123123",
			PAN:        "123123123123",
			CardSystem: "123123123123",
			ExpireDate: "123123123123",
			PIN:        "1234",
			CV2:        "123",
		}
		_ = respModels.NewCardMasked(card)
	}
}

func BenchmarkViaIter(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var card = models.Card{
			AccountID:  "123123123123",
			PAN:        "123123123123",
			CardSystem: "123123123123",
			ExpireDate: "123123123123",
			PIN:        "1234",
			CV2:        "123",
		}
		card.MaskSensitive()
	}
}
