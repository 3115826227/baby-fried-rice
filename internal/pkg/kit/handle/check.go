package handle

import "regexp"

var (
	phoneRegular = "^((13[0-9])|(14[5,7])|(15[0-3,5-9])|(17[0,3,5-8])|(18[0-9])|166|198|199|(147))\\d{8}$"
	phoneReg     = regexp.MustCompile(phoneRegular)

	emailRegular = `^[0-9a-z][_.0-9a-z-]{0,31}@([0-9a-z][0-9a-z-]{0,30}[0-9a-z]\.){1,4}[a-z]{2,4}$`
	emailReg     = regexp.MustCompile(emailRegular)
)

func PhoneInvalid(phone string) bool {
	return phoneReg.MatchString(phone)
}

func EmailInvalid(email string) bool {
	return emailReg.MatchString(email)
}
