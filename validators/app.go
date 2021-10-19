package validators

// Validator struct for app creation
type AppCreateData struct {
	Name         string   `json:"name" binding:"required" example:"MySuperApp"`
	IconUrl      string   `json:"icon_url" binding:"omitempty,url" example:"http://youriconurl.dev"`
	RedirectUrls []string `json:"redirect_urls" binding:"omitempty" example:"http://yourredirecturl.dev"`
}

// Validator struct for app update
type AppUpdateData struct {
	Name         string   `json:"name" binding:"omitempty" example:"MySuperApp"`
	IconUrl      string   `json:"icon_url" binding:"omitempty,url" example:"http://youriconurl.dev"`
	RedirectUrls []string `json:"redirect_urls" binding:"omitempty" example:"http://yourredirecturl.dev"`
}
