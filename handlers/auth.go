package handlers

import "errors"

var (
	// ErrNoAuthParams means that authentication parameters were not specified.
	// ! This should be depricated. We souldn't check if params were specifed. An error will be
	// ! returned when user was not found or password was incorrect which is actually what will
	// ! be used to show the user.
	ErrNoAuthParams = errors.New("auth: no auth params were specified")
	// ErrUserExists means that specified user already exists.
	ErrUserExists = errors.New("auth: user already exists")

	// ! We should move to JWT instead of using cookies and server side storage to store session
	// ! info.

	// ErrSessionExpired means that specified session is already expired.
	ErrSessionExpired = errors.New("auth: session has expired")
	// ErrSessionInvalid means that specified session is invalid.
	ErrSessionInvalid = errors.New("auth: session is invalid")
)

type authCreds struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
