package handlers

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/puremike/online_auction_api/internal/errs"
	"github.com/puremike/online_auction_api/internal/imagesuploader"
)

type ImageHandler struct {
	imageService imagesuploader.ImageServiceInterface
}

func NewImageHandler(imageService imagesuploader.ImageServiceInterface) *ImageHandler {
	return &ImageHandler{
		imageService: imageService,
	}
}

type imageRes struct {
	ImagePath string `json:"image_path"`
}

// UploadImage godoc
//
//	@Summary		Upload Image
//	@Description	Allows a user to upload an image to the server.
//	@Tags			Images
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			image	formData	file	true	"the image to upload"
//	@Success		200		{object}	gin.H	"image uploaded successfully"
//	@Failure		400		{object}	gin.H	"Bad Request - no file uploaded"
//	@Failure		500		{object}	gin.H	"Internal Server Error - failed to process file upload"
//	@Router			/auctions/image_upload [post]
func (i *ImageHandler) UploadImage(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		if errors.Is(err, http.ErrMissingFile) {
			errs.MapServiceErrors(c, errs.NewHTTPError("no file uploaded", http.StatusBadRequest))
		} else {
			log.Printf("failed to parse form file: %v", err)
			errs.MapServiceErrors(c, errs.NewHTTPError("failed to process file upload", http.StatusInternalServerError))
		}

		return
	}

	imagePath, err := i.imageService.UploadImage(c.Request.Context(), file)
	if err != nil {
		log.Printf("image service failed to upload file: %v", err)
		errs.MapServiceErrors(c, err)
		return
	}

	res := imageRes{
		ImagePath: imagePath,
	}
	c.JSON(http.StatusOK, res)

	log.Printf("image uploaded successfully: %s", imagePath)
}
