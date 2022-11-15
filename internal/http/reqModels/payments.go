package reqModels

import (
	"banking/pkg/types"
	"time"
)

type GeneratePaymentReq struct {
	MerchantAccountID string      `json:"merchant_account_id" binding:"required"`
	MerchantPaymentID string      `json:"merchant_payment_id,omitempty"`
	Amount            types.Money `json:"amount" binding:"required"`
	ExpiresAt         *time.Time  `json:"expires_at"`
	SuccessUrl        string      `json:"success_url"`
	FailUrl           string      `json:"fail_url"`
	IsMultiple        bool        `json:"is_multiple"`
	Purpose           string      `json:"purpose"`
	// Currency will be taken from the MerchantAccount
}

type PaymentForm struct {
	PAN            types.PAN        `form:"pan" binding:"required"`
	ExpireDate     types.ExpireDate `form:"expire_date" binding:"required"`
	CV2            types.CV2        `form:"cv2" binding:"required"`
	CardholderName string           `form:"cardholder_name" binding:"required"`
}
