package errs

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type APIError struct {
	Message string `json:"error"`
	Details string `json:"details,omitempty"`
}

func MapServiceErrors(c *gin.Context, err error) {
	var apiError APIError
	var statusCode int

	switch {
	case errors.Is(err, ErrUserNotFound):
		apiError = APIError{
			Message: err.Error(),
		}
		statusCode = http.StatusNotFound

	case errors.Is(err, ErrInvalidUserDetails):
		apiError = APIError{
			Message: err.Error(),
		}
		statusCode = http.StatusBadRequest

	case errors.Is(err, ErrTokenNotFound):
		apiError = APIError{
			Message: err.Error(),
		}
		statusCode = http.StatusNotFound

	case errors.Is(err, ErrInvalidCredentials):
		apiError = APIError{
			Message: err.Error(),
		}
		statusCode = http.StatusBadRequest

	case errors.Is(err, ErrRefreshTokenNotFound):
		apiError = APIError{
			Message: err.Error(),
		}
		statusCode = http.StatusNotFound

	case errors.Is(err, ErrInvalidRefreshToken):
		apiError = APIError{
			Message: err.Error(),
		}
		statusCode = http.StatusBadRequest

	case errors.Is(err, ErrEmailPasswordRequired):
		apiError = APIError{
			Message: err.Error(),
		}
		statusCode = http.StatusBadRequest

	case errors.Is(err, ErrFailedToCreateUser):
		apiError = APIError{
			Message: err.Error(),
		}
		statusCode = http.StatusBadRequest

	case errors.Is(err, ErrFailedToGenToken):
		apiError = APIError{
			Message: err.Error(),
		}
		statusCode = http.StatusBadRequest

	case errors.Is(err, ErrFailedToGenRefreshToken):
		apiError = APIError{
			Message: err.Error(),
		}
		statusCode = http.StatusBadRequest

	case errors.Is(err, ErrFailedToHashPassword):
		apiError = APIError{
			Message: err.Error(),
		}
		statusCode = http.StatusBadRequest

	case errors.Is(err, ErrFailedToStoreToken):
		apiError = APIError{
			Message: err.Error(),
		}
		statusCode = http.StatusBadRequest

	case errors.Is(err, ErrFailedToUpdateUser):
		apiError = APIError{
			Message: err.Error(),
		}
		statusCode = http.StatusBadRequest

	case errors.Is(err, ErrPasswordsDoNotMatch):
		apiError = APIError{
			Message: err.Error(),
		}
		statusCode = http.StatusBadRequest

	case errors.Is(err, ErrPasswordCannotBeSame):
		apiError = APIError{
			Message: err.Error(),
		}
		statusCode = http.StatusBadRequest

	case errors.Is(err, ErrInvalidPassword):
		apiError = APIError{
			Message: err.Error(),
		}
		statusCode = http.StatusBadRequest

	case errors.Is(err, ErrFailedToChangePassword):
		apiError = APIError{
			Message: err.Error(),
		}
		statusCode = http.StatusBadRequest

	case errors.Is(err, ErrAuctionNotFound):
		apiError = APIError{
			Message: err.Error(),
		}
		statusCode = http.StatusNotFound

	case errors.Is(err, ErrInvalidAuctionDetails):
		apiError = APIError{
			Message: err.Error(),
		}
		statusCode = http.StatusBadRequest

	case errors.Is(err, ErrFailedToCreateAuction):
		apiError = APIError{
			Message: err.Error(),
		}
		statusCode = http.StatusBadRequest

	case errors.Is(err, ErrFailedToUpdateAuction):
		apiError = APIError{
			Message: err.Error(),
		}
		statusCode = http.StatusBadRequest

	case errors.Is(err, ErrFailedToDeleteAuction):
		apiError = APIError{
			Message: err.Error(),
		}
		statusCode = http.StatusBadRequest

	case errors.Is(err, ErrAuctionNotOpenForBids):
		apiError = APIError{
			Message: err.Error(),
		}
		statusCode = http.StatusBadRequest

	case errors.Is(err, ErrBidTooLow):
		apiError = APIError{
			Message: err.Error(),
		}
		statusCode = http.StatusBadRequest

	case errors.Is(err, ErrBidBySeller):
		apiError = APIError{
			Message: err.Error(),
		}
		statusCode = http.StatusBadRequest

	case errors.Is(err, ErrPermissionDenied):
		apiError = APIError{
			Message: err.Error(),
		}
		statusCode = http.StatusUnauthorized

	case errors.Is(err, ErrAuctionAlreadyClosed):
		apiError = APIError{
			Message: err.Error(),
		}
		statusCode = http.StatusNotFound

	case errors.Is(err, ErrBidNotFound):
		apiError = APIError{
			Message: err.Error(),
		}
		statusCode = http.StatusNotFound

	case errors.Is(err, ErrFailedToSaveBid):
		apiError = APIError{
			Message: err.Error(),
		}
		statusCode = http.StatusBadRequest

	default:
		apiError = APIError{
			Message: "an unexpected error occurred",
			Details: err.Error(),
		}
		statusCode = http.StatusInternalServerError
	}

	c.JSON(statusCode, apiError)
}
