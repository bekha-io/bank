package reqModels

import (
	"banking/pkg/models"
	"banking/pkg/types"
)

type CreateUser struct {
	Login           string                   `json:"login" binding:"required"`
	Password        string                   `json:"password" binding:"required"`
	PhoneNumber     types.PhoneNumber        `json:"phone_number" binding:"required"`
	IndividualInfo  *IndividualCustomerInfo  `json:"individual_info,omitempty" structs:"individual_info,omitempty"`
	LegalEntityInfo *LegalEntityCustomerInfo `json:"legal_entity_info,omitempty" structs:"legal_entity_info,omitempty"`
}

type LegalEntityCustomerInfo struct {
	models.LegalEntityCustomerRaw
}

type IndividualCustomerInfo struct {
	models.IndividualCustomerRaw
}

type RefreshToken struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}
