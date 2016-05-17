// Package sso contains all the backend code to manage user authentication
// via email/password or social network.
package sso

import (
)

// Type Member contains basic member information.

type Member struct {
	Id        int64  `json:"-"`
	Email     string `json:"email"`
	ShortName string `json:"shortname"`
	FullName  string `json:"fullname"`
}
