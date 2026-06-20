package dto

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r *RegisterRequest) Validate() error {
	if r.Username == "" {
		return errField("username is required")
	}
	if len(r.Username) < 3 {
		return errField("username must be at least 3 characters")
	}
	if r.Email == "" {
		return errField("email is required")
	}
	if r.Password == "" {
		return errField("password is required")
	}
	if len(r.Password) < 6 {
		return errField("password must be at least 6 characters")
	}
	return nil
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r *LoginRequest) Validate() error {
	if r.Email == "" {
		return errField("email is required")
	}
	if r.Password == "" {
		return errField("password is required")
	}
	return nil
}

type AuthResponse struct {
	Token string      `json:"token"`
	User  interface{} `json:"user"`
}
