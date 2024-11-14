package email

import "net/mail"

func IsEmailFormatValid(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func FilterOutDistinctValidAndInvalidEmails(emails []string) (validEmails []string, invalidEmails []string) {
	validEmailSet := make(map[string]bool)
	invalidEmailSet := make(map[string]bool)
	for _, e := range emails {
		if IsEmailFormatValid(e) {
			validEmailSet[e] = true
		} else {
			invalidEmailSet[e] = true
		}
	}
	for e := range validEmailSet {
		validEmails = append(validEmails, e)
	}
	for e := range invalidEmailSet {
		invalidEmails = append(invalidEmails, e)
	}
	return
}
