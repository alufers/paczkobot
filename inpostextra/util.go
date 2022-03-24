package inpostextra

import "strings"

func NormalizePhoneNumber(phoneNumber string) string {
	phoneNumber = strings.Replace(phoneNumber, " ", "", -1)
	phoneNumber = strings.Replace(phoneNumber, "-", "", -1)

	// if len(phoneNumber) == 9 {
	// 	phoneNumber = phoneNumber
	// }
	return phoneNumber
}
