package val

import (
	"fmt"
	"regexp"
)

var (
	isValidateName         = regexp.MustCompile(`^[a-zA-Z0-9\s\-]+$`).MatchString
	isValidateBrand        = regexp.MustCompile(`^[a-zA-Z0-9\s\-]+$`).MatchString
	isValidateModel        = regexp.MustCompile(`^[a-zA-Z0-9\s\-]+$`).MatchString
	isValidateYear         = regexp.MustCompile(`^\d{4}$`).MatchString
	isValidateMileage      = regexp.MustCompile(`^(0|[1-9]\d{0,6})$`).MatchString
	isValidateTransmission = regexp.MustCompile(`^(?i)(manual|automatic)$`).MatchString
	isValidateFuelType     = regexp.MustCompile(`^(?i)(petrol|diesel|electric|hybrid)$`).MatchString
	isValidateCurrency     = regexp.MustCompile(`^(?i)(BNB|ETH|SOL)$`).MatchString
	isValidateLocation     = regexp.MustCompile(`^[a-zA-Z0-9\s,\.-]+$`).MatchString
	isValidateDescription  = regexp.MustCompile(`^.{10,}$`).MatchString
)

func ValidateString(value string, minLength int, maxlength int) error {

	n := len(value)

	if n < minLength || n > maxlength {
		return fmt.Errorf("must content from %d-%d characters", minLength, maxlength)
	}

	return nil
}

func ValidateName(value string) error {
	if err := ValidateString(value, 3, 100); err != nil {
		return err
	}

	if !isValidateName(value) {
		return fmt.Errorf("must contains only letters or spaces")
	}
	return nil
}

func ValidateBrand(value string) error {
	if err := ValidateString(value, 3, 100); err != nil {
		return err
	}

	if !isValidateBrand(value) {
		return fmt.Errorf("must contains letters or spaces")
	}
	return nil
}

func ValidateModel(value string) error {
	if err := ValidateString(value, 3, 100); err != nil {
		return err
	}

	if !isValidateModel(value) {
		return fmt.Errorf("must contains only letters or spaces")
	}
	return nil
}

func ValidateYear(value string) error {

	if !isValidateYear(value) {
		return fmt.Errorf("must contains only Number")
	}
	return nil
}

func ValidateMileage(value string) error {

	if !isValidateMileage(value) {
		return fmt.Errorf("must contains only number")
	}
	return nil
}

func ValidateTransmission(value string) error {
	if err := ValidateString(value, 3, 100); err != nil {
		return err
	}

	if !isValidateTransmission(value) {
		return fmt.Errorf("must contains only manual or automatic")
	}
	return nil
}

func ValidateFullType(value string) error {
	if err := ValidateString(value, 3, 100); err != nil {
		return err
	}

	if !isValidateFuelType(value) {
		return fmt.Errorf("must only contains one of petrol, diesel, electric, or hybrid")
	}
	return nil
}

func ValidateLocation(value string) error {
	if err := ValidateString(value, 3, 100); err != nil {
		return err
	}

	if !isValidateLocation(value) {
		return fmt.Errorf("must contains only letters or spaces")
	}
	return nil
}

func ValidateDescription(value string) error {
	if err := ValidateString(value, 3, 100); err != nil {
		return err
	}

	if !isValidateDescription(value) {
		return fmt.Errorf("must contains only letters or spaces")
	}
	return nil
}

func ValidateCurrency(value string) error {

	if !isValidateCurrency(value) {
		return fmt.Errorf("must only contains one of BNB, ETH or SOL")
	}
	return nil
}
