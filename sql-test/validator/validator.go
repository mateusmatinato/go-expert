package validator

import "fmt"

func ValidateDeleteParams(params []string) error {
	if len(params) != 1 {
		return fmt.Errorf("expected 1 param, got %d", len(params))
	}
	return nil
}

func ValidateUserParams(params []string) error {
	if len(params) != 2 {
		return fmt.Errorf("expected 2 params, got %d", len(params))
	}
	return nil
}
