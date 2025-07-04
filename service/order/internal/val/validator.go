package val

import (
	"fmt"
	"regexp"

	"github.com/google/uuid"
)

var (
	isValidateUsername = regexp.MustCompile(`^[a-z0-9_]+$`).MatchString
)

func ValidateString(value string, minLength int, maxlength int) error {

	n := len(value)

	if n < minLength || n > maxlength {
		return fmt.Errorf("must content from %d-%d characters", minLength, maxlength)
	}

	return nil
}

func ValidateUsername(value string) error {
	if err := ValidateString(value, 3, 100); err != nil {
		return err
	}

	if !isValidateUsername(value) {
		return fmt.Errorf("must contains only lowercase letters, digit or underscore")
	}
	return nil
}

func ValidateUUID(value string) error {
	_, err := uuid.Parse(value)

	if err != nil {
		return fmt.Errorf("invalid id")
	}
	return nil
}
