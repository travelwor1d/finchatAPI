package controller

import (
	"github.com/nyaruka/phonenumbers"
)

type Phone struct {
	CountryCode string `json:"countryCode" query:"countryCode" validate:"required"`
	Number      string `json:"phoneNumber" query:"phoneNumber" validate:"required"`
}

func (p Phone) Validate() bool {
	return validatePhonenumber(p.Number, p.CountryCode)
}

// formattedPhonenumber returns formatted phone number if .Number is valid phone number,
// otherwise returns empty string.
func (p Phone) formattedPhonenumber() string {
	num, err := phonenumbers.Parse(p.Number, p.CountryCode)
	if err != nil {
		return ""
	}
	return phonenumbers.Format(num, phonenumbers.E164)
}

func validatePhonenumber(val, countryCode string) bool {
	num, err := phonenumbers.Parse(val, countryCode)
	if err != nil {
		return false
	}
	return phonenumbers.IsValidNumber(num)
}
