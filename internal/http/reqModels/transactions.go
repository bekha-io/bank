package reqModels

import "banking/pkg/types"

type TransferMoneyReq struct {
	OrigAccountId   string             `json:"orig_account_id" binding:"required"`
	DestAccountId   *string            `json:"dest_account_id,omitempty"`
	DestPhoneNumber *types.PhoneNumber `json:"dest_phone_number,omitempty"`
	Amount          types.Money        `json:"amount" binding:"required"`
	Comment         string             `json:"comment,omitempty"`
}
