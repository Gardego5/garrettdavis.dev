//go:generate msgp
package model

type GHAccessToken struct {
	AccessToken string `json:"access_token" validate:"required"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type" validate:"required"`
}
