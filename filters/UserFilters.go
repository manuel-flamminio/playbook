package filters

type UserQueryFilters struct {
	Page        int
	Username    string
	DisplayName string `json:"display_name"`
}
