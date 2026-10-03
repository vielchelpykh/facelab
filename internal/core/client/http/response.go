package core_http_client

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var responseValidator = validator.New()

func DecodeAndValidateResponse(r *http.Response, dest any) error {
	if err := json.NewDecoder(r.Body).Decode(dest); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}

	if err := responseValidator.Struct(dest); err != nil {
		return fmt.Errorf("validate response: %w", err)
	}

	return nil
}
