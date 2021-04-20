package controller

import (
	"fmt"

	"github.com/nyaruka/phonenumbers"
)

type Phone struct {
	CountryCode string `json:"countryCode" query:"countryCode" validate:"required"`
	Number      string `json:"phonenumber" query:"phonenumber" validate:"validatePhonenumber"`
}

func (p Phone) ValidatePhonenumber(val string) bool {
	return validatePhonenumber(val, p.CountryCode)
}

// formattedPhonenumber returns formatted phonenumber if .Number is valid phone number,
// otherwise returns empty string.
func (p Phone) formattedPhonenumber() string {
	num, err := phonenumbers.Parse(p.Number, p.CountryCode)
	if err != nil {
		return ""
	}
	return phonenumbers.Format(num, phonenumbers.NATIONAL)
}

func validatePhonenumber(val, countryCode string) bool {
	num, err := phonenumbers.Parse(val, countryCode)
	if err != nil {
		return false
	}
	fmt.Println(num)
	return phonenumbers.IsValidNumber(num)
}
