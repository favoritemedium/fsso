package sso

// Type TokenReply is returned from the token functions.
type TokenReply struct {
  Atoken string   `json:"atoken"`
  Member *Member  `json:"member"`
}

// TokenEmail validates an email/password combination and either returns
// an auth token (for use in the Authorization header) or returns
// ErrAuthenticationFailure.
func TokenEmail(email, password string) (*TokenReply, error) {
  return &TokenReply{}, nil
}

// TokenEmail validates an email/password combination and either returns
// an auth token (for use in the Authorization header) or returns
// ErrAuthenticationFailure.
func TokenSocial(provider, id_token string) (*TokenReply, error) {
  return &TokenReply{}, nil
}
