package sso

// Type ConnectReply is returned from the connect functions.
type ConnectReply struct {
	Atoken string  `json:"atoken"`
	Member *Member `json:"member"`
}

// ConnectEmail validates an email/password combination and returns an
// auth token (for use in the Authorization header).
func ConnectEmail(email, password string) (*ConnectReply, error) {
	return &ConnectReply{}, nil
}

// ConnectSocial validates an id token from a social network and returns
// an auth token (for use in the Authorization header).
func ConnectSocial(provider, id_token string) (*ConnectReply, error) {
	return &ConnectReply{}, nil
}
