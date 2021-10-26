package validators

import "gandalf/bindings"

// Validator struct for oauth authorize app request
type OauthAuthorizeData struct {
	ClientID    string           `json:"client_id" binding:"required,uuid4" example:"4722679b-5a48-4e85-9084-605e8df610f4"`
	RedirectURI string           `json:"redirect_uri" binding:"required,url" example:"http://yourredirecturl.dev"`
	Scopes      []bindings.Scope `json:"scopes" binding:"required" example:"user:read"`
	State       string           `json:"state" binding:"omitempty" example:"iuywerghiuhg3487"`
}

// Validator struct for oauth token exchange
type OauthExchangeToken struct {
	GrantType         string `json:"grant_type" form:"grant_type" binding:"required,oneof='authorization_code'" example:"authorization_code"`
	ClientID          string `json:"client_id" form:"client_id" binding:"required" example:"4722679b-5a48-4e85-9084-605e8df610f4"`
	ClientSecret      string `json:"client_secret" form:"client_secret" binding:"required" example:"3i4u5h234ui5234bniuoo4i55543oi5jhio"`
	AuthorizationCode string `json:"code" form:"code" binding:"required" example:"iwuqebgrfweiur4"`
	RedirectUrl       string `json:"redirect_uri" form:"redirect_uri" binding:"required,url" example:"http://callback"`
}
