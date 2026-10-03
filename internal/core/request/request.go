package core_http_request

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var requestValidator = validator.New()

func DecodeAndValidateRequest(r *http.Request, videoDTO any) error {
	if err := json.NewDecoder(r.Body).Decode(videoDTO); err != nil {
		return fmt.Errorf("decode request: %w", err)
	}

	if err := requestValidator.Struct(videoDTO); err != nil {
		return fmt.Errorf("validate request: %w", err)
	}

	return nil
}
