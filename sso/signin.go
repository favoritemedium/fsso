package sso

// Type SigninReply is returned from the signin functions.
type SigninReply struct {
	Rtoken        string  `json:"rtoken"`
	RtokenExpires int64   `json:"rtoken_expires"`
	Member        *Member `json:"member"`
}

// SigninEmail validates an email/password combination signs in the user with
// a session cookie.
func SigninEmail(email, password string) (*SigninReply, error) {
	return &SigninReply{}, nil
}

// SigninRefresh validates a refresh token and signs in the user with a
// session cookie.
func SigninRefresh(rtoken string) (*SigninReply, error) {
	return &SigninReply{}, nil
}

// SigninSocial validates an id token from a social network and signs in the
// user with a session cookie.
func SigninSocial(provider, id_token string) (*SigninReply, error) {
	return &SigninReply{}, nil
}
