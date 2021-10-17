package validators

/*
AppCreateData -> app data for creation
*/
type AppCreateData struct {
	Name         string   `json:"name" binding:"required"`
	IconUrl      string   `json:"icon_url" binding:"omitempty,url"`
	RedirectUrls []string `json:"redirect_urls" binding:"omitempty"`
}

/*
AppUpdateData -> app data for update
*/
type AppUpdateData struct {
	Name         string   `json:"name" binding:"omitempty"`
	IconUrl      string   `json:"icon_url" binding:"omitempty,url"`
	RedirectUrls []string `json:"redirect_urls" binding:"omitempty"`
}
