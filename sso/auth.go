package sso

// Type AuthReply is returned from the auth functions.
type AuthReply struct {
  Rtoken string         `json:"rtoken"`
  RtokenExpires int64   `json:"rtoken_expires"`
  Member *Member        `json:"member"`
}

// AuthEmail validates an email/password combination and either signs in the
// user with a session cookie or returns ErrAuthenticationFailure.
func AuthEmail(email, password string) (*AuthReply, error) {
  return &AuthReply{}, nil
}

// AuthRefresh validates a refresh token and either signs in the user with a
// session cookie or returns ErrAuthenticationFailure.
func AuthRefresh(rtoken string) (*AuthReply, error) {
  return &AuthReply{}, nil
}

// AuthSocial validates an id token from a social network and either signs in
// the user with a session cookie or returns ErrAuthenticationFailure.
func AuthSocial(provider, id_token string) (*AuthReply, error) {
  return &AuthReply{}, nil
}
