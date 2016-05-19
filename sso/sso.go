// Package sso contains all the backend code to manage user authentication
// via email/password or social network.
package sso

import (
	"database/sql"
	"net/http"
	"time"
)

// Type Failure is the custom error type used for all API calls
type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e ErrorResponse) Error() string {
	return e.Message
}

var (
	ErrAuthenticationFailure   = ErrorResponse{"authfail", "Authentication failure."}
	ErrDisabledAccount         = ErrorResponse{"disabled", "That account is disabled."}
	ErrInvalidRtoken           = ErrorResponse{"rtoken", "Invalid or expired refresh token."}
	ErrUnknownProvider         = ErrorResponse{"provider", "Unknown social network provider."}
	ErrInvalidItoken           = ErrorResponse{"itoken", "Invalid id token."}
	ErrInvalidEmail            = ErrorResponse{"email", "Email address is not valid."}
	ErrDuplicateEmail          = ErrorResponse{"dupemail", "That email address is already registered."}
	ErrDuplicateAccount        = ErrorResponse{"dupaccount", "That account is already registered."}
	ErrInvalidVerifyCode       = ErrorResponse{"vcode", "Invalid or expired email verification code."}
	ErrBadPassword             = ErrorResponse{"password", "Password is not valid."}
	ErrMemberDetails           = ErrorResponse{"member", "Member details are incomplete."}
	ErrReauthenticationFailure = ErrorResponse{"reauthfail", "Reauthentication failure."}
	ErrNoEmail                 = ErrorResponse{"noemail", "This member doesn't have email/password signin."}
	ErrNoChange                = ErrorResponse{"nochange", "Old and new passwords are the same."}
	ErrYoureNotSure            = ErrorResponse{"notsure", "You're not sure."}
	ErrInvalidAccount          = ErrorResponse{"account", "Account is invalid."}
	ErrAlreadyPrimary          = ErrorResponse{"alreadyprimary", "That account is already primary."}
	ErrCantRemovePrimary       = ErrorResponse{"cantremoveprimary", "Primary account can't be removed."}
)

// Type Member contains basic member information.
type Member struct {
	id        int64
	Email     string `json:"email"`
	ShortName string `json:"shortname"`
	FullName  string `json:"fullname"`
	data      string
	roles     uint32
	aToken    string
}

// GetId returns the unique internal ID for the member.
func (m *Member) GetId() int64 {
	return m.id
}

// HasRole tells us if this member has any of the roles given, i.e.
//     m.HasRole(SuperRole | AdminRole)
//     m.HasRole(SuperRole) || m.HasRole(AdminRole)  // equivalent
func (m *Member) HasRole(roles uint32) bool {
	return m.roles&roles != 0
}

// HasRoles returns true if ths memberr has all of the roles given.
//     m.HasRole(SuperRole | AdminRole)
//     m.HasRole(SuperRole) && m.HasRole(AdminRole)  // equivalent
func (m *Member) HasRoles(roles uint32) bool {
	return m.roles&roles == roles
}

// SetSessionData writes to the active table an arbitrary string,
// which may be read back on a later request.
func (m *Member) SetSessionData(data string) error {
	if _, err := db.Exec(
		"UPDATE "+activeTable+" SET data=$2 WHERE atoken=$1",
		m.aToken, data); err != nil {
		return err
	}
	m.data = data
	return nil
}

// GetSessionData retrieves whatever was written previously using SetSessionData.
func (m *Member) GetSessionData() string {
	return m.data
}

// CurrentMember finds the currently connected member by checking the
// authorization header if present and otherwise the session cookie.
// Returns nil, nil if there is no current member.
func CurrentMember(r *http.Request) (*Member, error) {

	var (
		isCookie bool
		atoken   string
	)

	// These values are loaded from the active table
	var (
		memberId  int64
		useragent string
		isSession bool
		data      string
	)

	// Authorization header takes priority, so check it first
	atoken = r.Header.Get("Authorization")
	if atoken != "" {
		// We've found an Authorization header; check its validity
		if len(atoken) < 7 || atoken[0:6] != "token " {
			return nil, nil // Invalid header; ignore it and we're done.
		}
		atoken = atoken[6:len(atoken)]
	} else {
		// No Authorization header; look for a session cookie
		for _, cookie := range r.Cookies() {
			if cookie.Name == "sess" {
				atoken = cookie.Value
				break
			}
		}
		if atoken == "" {
			return nil, nil // No session cookie; we're done.
		}
		isCookie = true
	}

	// Attempt to update the last activity BEFORE we read it.  If we read it
	// first and the stale session cleaner kicks in before we've had a chance to
	// update the last active time, then our session will get dropped on the NEXT
	// call, which would be really awful nearly impossible-to-find sporadic  bug.
	// If the session has already been dropped, then this update has no effect.
	db.Exec("UPDATE "+activeTable+" SET active_at=?, ip=? WHERE atoken=?",
		timestamp(), r.RemoteAddr, atoken)

	if err := db.QueryRow(
		"SELECT member_id, useragent, is_session, data FROM "+activeTable+" WHERE atoken=?",
		atoken).Scan(&memberId, &useragent, &isSession, &data); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	// Check for fishiness.  If a session cookie is supplied as an Authorization
	// token or vice versa, that's fishy.  If the user agent has changed, that's
	// fishy.  If something is fishy, kill the session now.
	if isSession != isCookie { // TODO: user agent
		db.Exec("DELETE from "+activeTable+" WHERE atoken=?", atoken)
		return nil, nil
	}

	m := Member{}
	isActive := false
	if err := db.QueryRow(
		"SELECT email, fullname, shortname, is_active, roles FROM "+memberTable+" WHERE id=?",
		memberId).Scan(&m.Email, &m.FullName, &m.ShortName, &isActive, &m.roles); err != nil {
		return nil, err
	}
	if !isActive {
		// The account has been disabled since the last access, so cancel the session.
		db.Exec("DELETE from "+activeTable+" WHERE atoken=?", atoken)
		return nil, nil
	}

	m.id = memberId
	m.aToken = atoken
	m.data = data
	return &m, nil
}

// timestamp returns the rurrent unix time.  Tests could override this.
func timestamp() int64 {
  return time.Now().Unix()
}
