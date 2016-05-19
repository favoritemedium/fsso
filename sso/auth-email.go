package sso

import (
	"database/sql"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

// AutheEmail verifies an email/password combination and returns the id
// of the associated member record.
func AuthEmail(email, pw string) (int64, error) {

	var mid int64
	var pwhash []byte

	if err := db.QueryRow(
		"SELECT member_id, pwhash FROM "+emailAuthTable+" WHERE email=?",
		email).Scan(&mid, &pwhash); err != nil {
		if err == sql.ErrNoRows {
			return 0, ErrAuthenticationFailure
		}
		return 0, err
	}

	err := bcrypt.CompareHashAndPassword(pwhash, []byte(pw))
	if err != nil {
		return 0, ErrAuthenticationFailure
	}

	return mid, nil
}

// GenerateVcode creates a unique token that maps back to the specified email
// address.  Send this token as part of a link in a confirmation email, and
// use GetVerifiedEmail to change it back into a (now verified) email address.
func GenerateVcode(email string) (string, error) {

	// Check only the very basic email format. The real validation happens when
	// we actually send to the email address. This will allow  unconventional
	// email addresses to still be used.
	if strings.Count(email, "@") != 1 || email[0] == '@' || email[len(email)-1] == '@' {
		return "", ErrInvalidEmail
	}

  // Have the verify code be valid for 24 hours.
	expiry := timestamp() + 86400

	for {
		vcode := RandomToken(32)
		if _, err := db.Exec(
			"INSERT INTO "+emailVerifyTable+" (vtoken, email, expires_at) VALUES (?,?,?)",
			vcode, email, expiry); err != nil {
			if isDuplicate(err) {
				continue
			}
			return "", err
		}
		return vcode, nil
	}
}

// GetVerifiedEmail gets a verified email address form the VerifyEmailTable
// using a verify code.
//
// GetVerifiedEmail has the side effect of delaying the expiration of of the
// verify code if it's about to expire.  After calling GetVerifiedEmail, there
// is always at least 10 minutes of validity left for the code, as the code
// needs to be used again after the user fills in and submits their personal
// information.
func GetVerifiedEmail(vcode string) (string, error) {

	var email string
	var expiry int64
	now := timestamp()

	if err := db.QueryRow(
		"SELECT email, expires_at FROM "+emailVerifyTable+" WHERE vtoken=?",
		vcode, now).Scan(&email, &expiry); err != nil {
		if err == sql.ErrNoRows {
			return "", ErrInvalidVerifyCode
		}
		return "", err
	}
	if expiry < now {
		return "", ErrInvalidVerifyCode
	}

	// Ensure we have at least 10 minutes of validity on this verify code.
	if expiry-now < 600 {
		db.Exec(
			"UPDATE "+emailVerifyTable+" SET expires_at=? WHERE vtoken=?",
			now+600, vcode)
	}

	return email, nil
}

// addEmailAuth creates a new email auth record and points it to an existing
// member record.
func addEmailAuth(email, pw string, mid int64, isPrimary bool) (err error) {

	t, err := db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err == nil {
			t.Rollback()
		} else {
			t.Commit()
		}
	}()

	pwhash, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	if _, err := t.Exec(
		"INSERT INTO "+emailAuthTable+" (member_id, email, pwhash, pwchanged_at, is_primary) VALUES (?, ?, ?, ?)",
		mid, email, pwhash, timestamp(), isPrimary); err != nil {
		if isDuplicate(err) {
			return ErrDuplicateEmail
		}
		return err
	}

	if isPrimary && email != "" {
		if _, err := t.Exec(
			"UPDATE "+memberTable+" SET email=? WHERE id=?",
			email, mid); err != nil {
			return err
		}
	}

	return nil
}
