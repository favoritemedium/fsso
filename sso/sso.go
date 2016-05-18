// Package sso contains all the backend code to manage user authentication
// via email/password or social network.
package sso

import (
)

// Type Failure is the custom error type used for all API calls
type ErrorResponse struct {
	Code    string  `json:"code"`
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
	Id        int64  `json:"-"`
	Email     string `json:"email"`
	ShortName string `json:"shortname"`
	FullName  string `json:"fullname"`
}
