package imagesuploader

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/puremike/online_auction_api/internal/errs"
	"github.com/puremike/online_auction_api/pkg"
)

type ImageServiceInterface interface {
	UploadImage(ctx context.Context, file *multipart.FileHeader) (string, error)
}

type ImageService struct {
	bucket *string
	client *s3.Client
	region string
}

func NewImageService(bucket string) *ImageService {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	client := s3.NewFromConfig(cfg)

	return &ImageService{
		bucket: aws.String(bucket),
		client: client,
		region: cfg.Region,
	}
}

const (
	QueryDefaultContext = 5 * time.Second
	PresignURLExpires   = 1 * time.Hour
)

func (i *ImageService) UploadImage(ctx context.Context, file *multipart.FileHeader) (string, error) {

	ctx, cancel := context.WithTimeout(ctx, QueryDefaultContext)
	defer cancel()

	// basic file size validation
	maxUploadFileSize := 3 * 1024 * 1024 //3mb

	if file.Size > int64(maxUploadFileSize) {
		log.Printf("file size exceeds limit of %d MB", maxUploadFileSize/1024/1024)
		return "", errs.NewHTTPError(fmt.Sprintf("file size exceeds limit of %d MB", maxUploadFileSize/1024/1024), http.StatusBadRequest)
	}

	allowedContentType := map[string]string{
		"image/jpeg": ".jpg",
		"image/png":  ".png",
	}

	// save the file
	src, err := file.Open()
	if err != nil {
		log.Printf("failed to open file: %v", err)
		return "", errs.NewHTTPError("failed to open file", http.StatusInternalServerError)
	}

	defer src.Close()

	buffer := make([]byte, 512)
	if _, err := src.Read(buffer); err != nil {
		log.Printf("failed to read file: %v", err)
		return "", errs.NewHTTPError("failed to read file for MIME type detection", http.StatusInternalServerError)
	}

	// reset reader
	if _, err := src.Seek(0, io.SeekStart); err != nil {
		log.Printf("failed to reset reader: %v", err)
		return "", errs.NewHTTPError("failed to reset reader", http.StatusInternalServerError)
	}

	mimeType := http.DetectContentType(buffer)
	ext, allowed := allowedContentType[mimeType]
	if !allowed {
		log.Printf("unsupported file type: %s. Only JPEG and PNG are allowed.", mimeType)
		return "", errs.NewHTTPError(fmt.Sprintf("unsupported file type: %s. Only JPEG and PNG are allowed.", mimeType), http.StatusBadRequest)
	}

	key := "uploads/" + pkg.GenerateUniqueFileName() + strings.ToLower(ext)

	// Decode Image
	var img image.Image
	switch mimeType {
	case "image/jpeg":
		img, err = jpeg.Decode(src)
	case "image/png":
		img, err = png.Decode(src)
	default:
		log.Printf("unsupported file type: %s. Only JPEG, PNG are allowed.", mimeType)
		return "", errs.NewHTTPError(fmt.Sprintf("unsupported file type: %s. Only JPEG, PNG, GIF are allowed.", mimeType), http.StatusBadRequest)
	}
	if err != nil {
		log.Printf("failed to decode file: %v", err)
		return "", errs.NewHTTPError("failed to decode file", http.StatusInternalServerError)
	}

	// image compression
	var compressedBuffer bytes.Buffer
	switch mimeType {
	case "image/jpeg":
		err = jpeg.Encode(&compressedBuffer, img, &jpeg.Options{
			Quality: 80,
		})
	case "image/png":
		err = png.Encode(&compressedBuffer, img)
	}
	if err != nil {
		log.Printf("failed to compress file: %v", err)
		return "", errs.NewHTTPError("failed to compress file", http.StatusInternalServerError)
	}

	// upload to s3
	if _, err := i.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      i.bucket,
		Key:         &key,
		Body:        bytes.NewReader(compressedBuffer.Bytes()),
		ContentType: &mimeType,
	}); err != nil {
		log.Printf("failed to upload file to S3: %v", err)
		return "", errs.NewHTTPError("failed to upload file to S3", http.StatusInternalServerError)
	}

	// generate presigned URL

	// presignClient := s3.NewPresignClient(i.client)
	// presignedURL, err := presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
	// 	Bucket: i.bucket,
	// 	Key:    &key,
	// }, s3.WithPresignExpires(PresignURLExpires))
	// if err != nil {
	// 	log.Printf("failed to generate presigned URL: %v", err)
	// 	return "", errs.NewHTTPError("failed to generate presigned URL", http.StatusInternalServerError)
	// }

	// return presignedURL.URL, nil

	publicUR := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", *i.bucket, i.region, key)
	return publicUR, nil
}

//  LOCAL STORAGE

// type ImageService struct {
// 	uploadDir string
// }

// func NewImageService(uploadDir string) *ImageService {
// 	// check if upload directory exist
// 	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
// 		if err := os.MkdirAll(uploadDir, 0755); err != nil {
// 			panic(fmt.Sprintf("Failed to create upload directory %s: %v", uploadDir, err))
// 		}
// 	}

// 	return &ImageService{
// 		uploadDir: uploadDir,
// 	}
// }

// func (i *ImageService) UploadImage(ctx context.Context, file *multipart.FileHeader) (string, error) {

// 	// basic file size validation
// 	maxUploadFileSize := 3 * 1024 * 1024 //3mb

// 	if file.Size > int64(maxUploadFileSize) {
// 		log.Printf("file size exceeds limit of %d MB", maxUploadFileSize/1024/1024)
// 		return "", errs.NewHTTPError(fmt.Sprintf("file size exceeds limit of %d MB", maxUploadFileSize/1024/1024), http.StatusBadRequest)
// 	}

// 	allowedContentType := map[string]string{
// 		"image/jpeg": ".jpg",
// 		"image/png":  ".png",
// 		"image/gif":  ".gif",
// 	}

// 	// save the file
// 	src, err := file.Open()
// 	if err != nil {
// 		log.Printf("failed to open file: %v", err)
// 		return "", errs.NewHTTPError("failed to open file", http.StatusInternalServerError)
// 	}

// 	buffer := make([]byte, 512)
// 	if _, err := src.Read(buffer); err != nil {
// 		log.Printf("failed to read file: %v", err)
// 		return "", errs.NewHTTPError("failed to read file for MIME type detection", http.StatusInternalServerError)
// 	}

// 	src.Seek(0, io.SeekStart) // reset reader

// 	mimeType := http.DetectContentType(buffer)
// 	ext, allowed := allowedContentType[mimeType]
// 	if !allowed {
// 		log.Printf("unsupported file type: %s. Only JPEG, PNG, GIF are allowed.", mimeType)
// 		return "", errs.NewHTTPError(fmt.Sprintf("unsupported file type: %s. Only JPEG, PNG, GIF are allowed.", mimeType), http.StatusBadRequest)
// 	}

// 	defer src.Close()

// 	uniqueFileName := pkg.GenerateUniqueFileName() + ext
// 	filePath := filepath.Join(i.uploadDir, uniqueFileName)

// 	dst, err := os.Create(filePath)
// 	if err != nil {
// 		log.Printf("failed to create file on server: %v", err)
// 		return "", errs.NewHTTPError("failed to create file on server", http.StatusInternalServerError)
// 	}

// 	defer dst.Close()

// 	if _, err := io.Copy(dst, src); err != nil {
// 		log.Printf("failed to save file to disk: %v", err)
// 		return "", errs.NewHTTPError("failed to save file to disk", http.StatusInternalServerError)
// 	}

// 	publicPath := fmt.Sprintf("/uploads/%s", uniqueFileName)
// 	return publicPath, nil
// }
