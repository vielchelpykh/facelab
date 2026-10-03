package core_http_client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	domain "github.com/vielchelpykh/facelab/internal/core/domains"
)

type ClientDTO struct {
	FileName string `json:"fileName"`
	FilePath string `json:"filePath"`
}

func (c *Client) BlurVideo(ctx context.Context, fileName string, filePath string) (domain.ClientDomain, error) {
	clientDTO := ClientDTO{
		FileName: fileName,
		FilePath: filePath,
	}

	jsonBody, err := json.Marshal(clientDTO)
	if err != nil {
		return domain.ClientDomain{}, fmt.Errorf("create json body for request: %w", err)
	}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.baseURL+"/videos",
		bytes.NewReader(jsonBody),
	)
	if err != nil {
		return domain.ClientDomain{}, fmt.Errorf("create http request: %w", err)
	}

	response, err := c.http.Do(request)
	if err != nil {
		return domain.ClientDomain{}, fmt.Errorf("do http request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return domain.ClientDomain{}, fmt.Errorf("blur service return status code: %d", response.StatusCode)
	}

	var clientDomain domain.ClientDomain
	if err := DecodeAndValidateResponse(response, clientDomain); err != nil {
		return domain.ClientDomain{}, fmt.Errorf("decode and validate response: %w", err)
	}

	return clientDomain, nil
}
