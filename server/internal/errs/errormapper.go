package errs

import (
	"errors" // Import fmt for creating new errors
	"net/http"

	"github.com/gin-gonic/gin"
)

type APIError struct {
	Message string `json:"error"`
	Details string `json:"details,omitempty"`
}

type HTTPError interface {
	Error() string
	StatusCode() int
	PublicMessage() string
}

// baseHTTPError is a concrete type that implements the HTTPError interface.
type baseHTTPError struct {
	msg        string
	statusCode int
}

func (e *baseHTTPError) Error() string {
	return e.msg
}

func (e *baseHTTPError) StatusCode() int {
	return e.statusCode
}

func (e *baseHTTPError) PublicMessage() string {
	return e.msg
}

func NewHTTPError(msg string, statusCode int) HTTPError {
	return &baseHTTPError{
		msg:        msg,
		statusCode: statusCode,
	}
}

var (
	ErrUserNotFound           = NewHTTPError("user not found", http.StatusNotFound)
	ErrInvalidUserDetails     = NewHTTPError("invalid user details", http.StatusBadRequest)
	ErrInvalidCredentials     = NewHTTPError("invalid credentials", http.StatusBadRequest)
	ErrEmailPasswordRequired  = NewHTTPError("email and password are required", http.StatusBadRequest)
	ErrFailedToCreateUser     = NewHTTPError("failed to create user", http.StatusBadRequest)
	ErrFailedToUpdateUser     = NewHTTPError("failed to update user", http.StatusBadRequest)
	ErrPasswordsDoNotMatch    = NewHTTPError("passwords do not match", http.StatusBadRequest)
	ErrPasswordCannotBeSame   = NewHTTPError("new password cannot be the same as old password", http.StatusBadRequest)
	ErrInvalidPassword        = NewHTTPError("invalid password", http.StatusBadRequest)
	ErrFailedToChangePassword = NewHTTPError("failed to change password", http.StatusBadRequest)

	ErrTokenNotFound           = NewHTTPError("token not found", http.StatusNotFound)
	ErrInvalidRefreshToken     = NewHTTPError("invalid refresh token", http.StatusBadRequest)
	ErrRefreshTokenNotFound    = NewHTTPError("refresh token not found", http.StatusNotFound)
	ErrFailedToGenToken        = NewHTTPError("failed to generate token", http.StatusInternalServerError)
	ErrFailedToGenRefreshToken = NewHTTPError("failed to generate refresh token", http.StatusInternalServerError)
	ErrFailedToStoreToken      = NewHTTPError("failed to store token", http.StatusInternalServerError)
	ErrFailedToHashPassword    = NewHTTPError("failed to hash password", http.StatusInternalServerError)

	ErrAuctionNotFound             = NewHTTPError("auction not found", http.StatusNotFound)
	ErrInvalidAuctionDetails       = NewHTTPError("invalid auction details", http.StatusBadRequest)
	ErrFailedToCreateAuction       = NewHTTPError("failed to create auction", http.StatusBadRequest)
	ErrFailedToUpdateAuction       = NewHTTPError("failed to update auction", http.StatusBadRequest)
	ErrFailedToDeleteAuction       = NewHTTPError("failed to delete auction", http.StatusBadRequest)
	ErrAuctionNotOpenForBids       = NewHTTPError("auction not open for bids", http.StatusBadRequest)
	ErrBidTooLow                   = NewHTTPError("bid too low", http.StatusBadRequest)
	ErrBidBySeller                 = NewHTTPError("seller cannot bid on their own auction", http.StatusBadRequest)
	ErrPermissionDenied            = NewHTTPError("permission denied", http.StatusUnauthorized)
	ErrAuctionAlreadyClosed        = NewHTTPError("auction already closed", http.StatusNotFound)
	ErrBidNotFound                 = NewHTTPError("bid not found", http.StatusNotFound)
	ErrFailedToSaveBid             = NewHTTPError("failed to save bid", http.StatusBadRequest)
	ErrFailedToGetHighestBid       = NewHTTPError("failed to get highest bid", http.StatusNotFound)
	ErrNotificationNotFound        = NewHTTPError("notification not found", http.StatusNotFound)
	ErrDutchBidMustMatchCurrent    = NewHTTPError("dutch bid must match current auction price", http.StatusBadRequest)
	ErrDutchAuctionAlreadyWon      = NewHTTPError("dutch auction already won", http.StatusBadRequest)
	ErrDuplicateSealedBid          = NewHTTPError("duplicate sealed bid", http.StatusBadRequest)
	ErrFailedToDeleteBids          = NewHTTPError("failed to delete bids", http.StatusBadRequest)
	ErrFailedToDeleteNotifications = NewHTTPError("failed to delete notifications", http.StatusBadRequest)

	// Payment related errors
	ErrFailedToCreateStripeCheckout   = NewHTTPError("failed to create Stripe checkout session", http.StatusInternalServerError)
	ErrAmountCannotBeNegative         = NewHTTPError("amount cannot be negative", http.StatusBadRequest)
	ErrFailedToGetPaymentSession      = NewHTTPError("failed to get payment session", http.StatusNotFound)
	ErrMissingRequiredMetadata        = NewHTTPError("missing required metadata", http.StatusBadRequest)
	ErrMissingRequiredSessionMetadata = NewHTTPError("missing required session metadata", http.StatusBadRequest)
	ErrFailedToUnmarshalEvent         = NewHTTPError("failed to unmarshal Stripe event", http.StatusBadRequest)
	ErrFailedToGetPayment             = NewHTTPError("failed to get payment record", http.StatusNotFound)
	ErrFailedToUpdatePayment          = NewHTTPError("failed to update payment record", http.StatusInternalServerError)
	ErrFailedToCreatePayment          = NewHTTPError("failed to create payment record", http.StatusInternalServerError)
)

// MapServiceErrors maps service-level errors to appropriate HTTP responses.
func MapServiceErrors(c *gin.Context, err error) {
	var apiError APIError
	var statusCode int

	var httpErr HTTPError
	if errors.As(err, &httpErr) {
		apiError = APIError{
			Message: httpErr.PublicMessage(),
		}
		statusCode = httpErr.StatusCode()
	} else {
		apiError = APIError{
			Message: "an unexpected error occurred",
		}
		statusCode = http.StatusInternalServerError
	}

	c.JSON(statusCode, apiError)
}
