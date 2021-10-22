package validators

import "gandalf/bindings"

// Validator struct for oauth authorize app request
type OauthAuthorizeData struct {
	ClientID    string           `json:"client_id" binding:"required,uuid4" example:"4722679b-5a48-4e85-9084-605e8df610f4"`
	RedirectURI string           `json:"redirect_uri" binding:"required,url" example:"http://yourredirecturl.dev"`
	Scopes      []bindings.Scope `json:"scopes" binding:"required" example:"user:read"`
	State       string           `json:"state" binding:"omitempty" example:"iuywerghiuhg3487"`
}
