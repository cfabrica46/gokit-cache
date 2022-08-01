package entity

// IDUsernameEmailSecretRequest ...
type IDUsernameEmailSecretRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Secret   string `json:"secret"`
	ID       int    `json:"id"`
}

// TokenSecretRequest ...
type TokenSecretRequest struct {
	Token  string `json:"token"`
	Secret string `json:"secret"`
}

// Token ...
type Token struct {
	Token string `json:"token"`
}

// IDUsernameEmailErrResponse ...
type IDUsernameEmailErrResponse struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Err      string `json:"err,omitempty"`
	ID       int    `json:"id"`
}

// ErrorResponse ...
type ErrorResponse struct {
	Err string `json:"err,omitempty"`
}

// CheckErrResponse ...
type CheckErrResponse struct {
	Err   string `json:"err,omitempty"`
	Check bool   `json:"check"`
}
