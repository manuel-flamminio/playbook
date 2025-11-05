package requests

type CreateUserRequest struct {
	DisplayName      string `json:"display_name,omitempty" example:"The Lion"`
	Username         string `json:"username,omitempty" example:"the.lion@example.com"`
	EncodedUserImage string `json:"user_image,omitempty" example:"iVBORw0KGgoAAAANSUhEUgAAAAgAAAAIAQMAAAD+wSzIAAAABlBMVEX///+/v7+jQ3Y5AAAADklEQVQI12P4AIX8EAgALgAD/aNpbtEAAAAASUVORK5CYII"`
	Password         string `json:"password,omitempty" example:"SuperSecurePassword1234!"`
}

type UpdateUserDisplayNameRequest struct {
	DisplayName string `json:"display_name,omitempty" example:"The Lion"`
}

type UpdateUserPasswordRequest struct {
	Password string `json:"password,omitempty" example:"SuperSecurePassword1234!"`
}

type UserLoginRequest struct {
	Username string `json:"username,omitempty" example:"the.lion@example.com"`
	Password string `json:"password,omitempty" example:"SuperSecurePassword1234!"`
}

type UserFilters struct {
	Page        int    `form:"page"`
	Username    string `form:"username"`
	DisplayName string `form:"display_name"`
}
