package errs

import "errors"

var (
	ErrUserNotFound            = errors.New("user not found")
	ErrTokenNotFound           = errors.New("token not found")
	ErrRefreshTokenNotFound    = errors.New("refresh token not found")
	ErrInvalidRefreshToken     = errors.New("invalid refresh token")
	ErrFailedToCreateUser      = errors.New("failed to create user")
	ErrInvalidUserDetails      = errors.New("invalid user details")
	ErrInvalidCredentials      = errors.New("invalid credentials")
	ErrEmailPasswordRequired   = errors.New("email or password required")
	ErrFailedToChangePassword  = errors.New("failed to change password")
	ErrFailedToGenToken        = errors.New("failed to generate token")
	ErrFailedToGenRefreshToken = errors.New("failed to generate refresh token")
	ErrFailedToHashPassword    = errors.New("failed to hash password")
	ErrFailedToStoreToken      = errors.New("failed to store token in database")
	ErrFailedToUpdateUser      = errors.New("failed to update user profile")
	ErrPasswordsDoNotMatch     = errors.New("passwords do not match")
	ErrInvalidPassword         = errors.New("invalid password")
	ErrPasswordCannotBeSame    = errors.New("new password cannot be the same as the old password")
)
