package imageresizer

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"gopkg.in/gographics/imagick.v2/imagick"

	"cloud.google.com/go/storage"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

func init() {
	imagick.Initialize()
	functions.HTTP("ResizeImage", ResizeImage)
}

func errorResponse(w http.ResponseWriter, message string, errType string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(code)
	errorResponse := map[string]interface{}{
		"error": map[string]interface{}{
			"message": message,
			"type":    errType,
			"code":    code,
		},
	}
	json.NewEncoder(w).Encode(errorResponse)
}

func ResizeImage(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	originalImageBlob, format, err := fetchImageFromStorage(ctx, r)
	if err != nil {
		errorResponse(w, fmt.Sprintf("Fetch error: %v", err), "FetchError", http.StatusInternalServerError)
		return
	}

	resizedImageBlob, err := resizeImage(originalImageBlob, r)
	if err != nil {
		errorResponse(w, fmt.Sprintf("Resize error: %v", err), "ResizeError", http.StatusInternalServerError)
		return
	}

	err = writeToResponse(w, resizedImageBlob, format)
	if err != nil {
		errorResponse(w, fmt.Sprintf("Write error: %v", err), "WriteError", http.StatusInternalServerError)
		return
	}
}

func fetchImageFromStorage(ctx context.Context, r *http.Request) ([]byte, string, error) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, "", err
	}

	imageName := r.URL.Query().Get("image_name")
	bucketName := os.Getenv("GCS_BUCKET_NAME")
	if bucketName == "" {
		return nil, "", fmt.Errorf("GCS_BUCKET_NAME environment variable is not set")
	}

	bucket := client.Bucket(bucketName)
	object := bucket.Object(imageName)

	attrs, err := object.Attrs(ctx)
	if err != nil {
		return nil, "", err
	}

	reader, err := object.NewReader(ctx)
	if err != nil {
		return nil, "", err
	}
	defer reader.Close()

	buf, err := io.ReadAll(reader)
	if err != nil {
		return nil, "", err
	}

	return buf, attrs.ContentType, nil
}

func resizeImage(originalImageBlob []byte, r *http.Request) ([]byte, error) {
	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	err := mw.ReadImageBlob(originalImageBlob)
	if err != nil {
		return nil, err
	}

	originalWidth := mw.GetImageWidth()
	originalHeight := mw.GetImageHeight()

	widthParam := r.URL.Query().Get("w")
	heightParam := r.URL.Query().Get("h")

	width := uint(0)
	height := uint(0)

	if widthParam == "" && heightParam == "" {
		return originalImageBlob, nil
	}

	if widthParam != "" {
		intWidth, err := strconv.Atoi(widthParam)
		if err != nil || intWidth <= 0 {
			return nil, fmt.Errorf("Invalid width parameter")
		}
		width = uint(intWidth)
	}

	if heightParam != "" {
		intHeight, err := strconv.Atoi(heightParam)
		if err != nil || intHeight <= 0 {
			return nil, fmt.Errorf("Invalid height parameter")
		}
		height = uint(intHeight)
	}

	if width == 0 && height != 0 {
		aspectRatio := float64(originalWidth) / float64(originalHeight)
		width = uint(float64(height) * aspectRatio)
	}

	if height == 0 && width != 0 {
		aspectRatio := float64(originalHeight) / float64(originalWidth)
		height = uint(float64(width) * aspectRatio)
	}

	if err := mw.ResizeImage(width, height, imagick.FILTER_LANCZOS, 1.0); err != nil {
		return nil, err
	}

	return mw.GetImageBlob(), nil
}

func writeToResponse(w http.ResponseWriter, resizedImageBlob []byte, originalFormat string) error {
	w.Header().Set("Content-Type", originalFormat)
	w.Header().Set("Cache-Control", "public, max-age=86400")
	_, err := w.Write(resizedImageBlob)
	return err
}
