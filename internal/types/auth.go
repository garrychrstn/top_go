package types

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ForgotPassRequest struct {
	Password    string `json:"password"`
	OldPassword string `json:"old_password"`
}
