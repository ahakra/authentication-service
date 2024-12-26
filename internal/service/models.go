package service

// ------------------------------

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginInResponse struct {
	AuthorizationToken string `json:"authorization_token"`
}

// -------------------------------

type AddPermissionInput struct {
	Permission string `json:"permission"`
}

type DeletePermissionInput struct {
	Permission string `json:"permission"`
}

//------------------------------

type RegenerateEmailTokenInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type ReGenerateEmailTokenResponse struct {
	ID      int64  `json:"-"`
	Email   string `json:"-"`
	IsMatch bool   `json:"-"`
	Token   string `json:"token"`
}
type ValidateTokenInput struct {
	Token string `json:"token"`
}
type ValidateTokenResponse struct {
	Token   string `json:"token"`
	IsValid bool   `json:"is_valid"`
}

//---------------------------------

type UserRegisterInput struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Activated bool   `json:"activated"`
}

type UserResponse struct {
	ID                int64  `json:"id"`
	Name              string `json:"name"`
	Email             string `json:"email"`
	VerificationToken string `json:"verification_token"`
	Password          []byte `json:"-"`
	Activated         bool   `json:"activated"`
}
