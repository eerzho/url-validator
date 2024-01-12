package url

import "url-validator/internal/lib/api/response"

type ValidateRequest struct {
	Domain string   `json:"domain" validate:"required"`
	Urls   []string `json:"urls" validate:"required,min=1,max=300"`
}

type ValidateResponse struct {
	response.Response
	Validated map[string]int32 `json:"validated"`
}
