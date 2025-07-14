package val

import (
	"fmt"
	"math/big"
	"regexp"
	"strings"

	"github.com/google/uuid"
)

var (
	isValidateUsername  = regexp.MustCompile(`^[a-z0-9_]+$`).MatchString
	isValidateCurrency  = regexp.MustCompile(`^(?i)(BNB|ETH|SOL)$`).MatchString
	isValidateNetwork   = regexp.MustCompile(`^(?i)(bsc|ethereum|solana)$`).MatchString
	isValidateEthandBsc = regexp.MustCompile(`^0x[a-fA-F0-9]{40}$`).MatchString
	isValidateSol       = regexp.MustCompile(`^[1-9A-HJ-NP-Za-km-z]{32,44}$`).MatchString
	isValidateTxHashEth = regexp.MustCompile(`^0x([A-Fa-f0-9]{64})$`).MatchString
	isValidateTxHashSol = regexp.MustCompile(`^[1-9A-HJ-NP-Za-km-z]{64,88}$`).MatchString
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

func ValidateCurrency(value string) error {

	if !isValidateCurrency(value) {
		return fmt.Errorf("must only contains one of BNB, ETH or SOL")
	}
	return nil
}

func ValidateNetwork(value string) error {

	if !isValidateNetwork(value) {
		return fmt.Errorf("must only contains one of bsc, ethereum or solana")
	}
	return nil
}

func ValidateWalletAddress(network, address string) error {
	address = strings.TrimSpace(address)

	switch network {
	case "ethereum", "bsc":
		if !isValidateEthandBsc(address) {
			return fmt.Errorf("invalid Ethereum/BSC wallet address format")
		}
	case "solana":
		if !isValidateSol(address) {
			return fmt.Errorf("invalid Solana wallet address format")
		}
	default:
		return fmt.Errorf("unsupported network")
	}

	return nil
}

func ValidateTxHash(network, address string) error {
	address = strings.TrimSpace(address)

	switch network {
	case "ethereum", "bsc":
		if !isValidateTxHashEth(address) {
			return fmt.Errorf("invalid Ethereum/BSC transaction hash format")
		}
	case "solana":
		if !isValidateTxHashSol(address) {
			return fmt.Errorf("invalid Solana transaction hash format")
		}
	default:
		return fmt.Errorf("unsupported network")
	}

	return nil
}
func ValidateAmount(amountStr string) error {
	amountStr = strings.TrimSpace(amountStr)
	if amountStr == "" {
		return fmt.Errorf("amount cannot be empty")
	}

	amount, ok := new(big.Float).SetString(amountStr)
	if !ok {
		return fmt.Errorf("invalid amount format")
	}

	zero := big.NewFloat(0)
	if amount.Cmp(zero) <= 0 {
		return fmt.Errorf("amount must be greater than zero")
	}

	return nil
}
