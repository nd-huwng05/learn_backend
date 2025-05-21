package banana

import "errors"

var (
	UserError      = errors.New("User Exists")
	SignUpFailed   = errors.New("Sign Up Failed")
	UserNotFound   = errors.New("User Not Found")
	UserNotUpdated = errors.New("User Not Updated")
)
