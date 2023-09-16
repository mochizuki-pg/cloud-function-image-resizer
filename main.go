package resizeimage

import (
	"context"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"net/http"
	"os"
	"strconv"

	"cloud.google.com/go/storage"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"golang.org/x/image/draw"
)

func init() {
	functions.HTTP("ResizeImage", ResizeImage)
}

func ResizeImage(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	image, format, err := fetchImageFromStorage(ctx, r)
	if err != nil {
		http.Error(w, fmt.Sprintf("Fetch error: %v", err), http.StatusInternalServerError)
		return
	}

	resizedImage, err := resizeImage(image, r)
	if err != nil {
		http.Error(w, fmt.Sprintf("Resize error: %v", err), http.StatusInternalServerError)
		return
	}

	err = writeToResponse(w, resizedImage, format)
	if err != nil {
		http.Error(w, fmt.Sprintf("Write error: %v", err), http.StatusInternalServerError)
		return
	}
}

func fetchImageFromStorage(ctx context.Context, r *http.Request) (image.Image, string, error) {
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
	reader, err := object.NewReader(ctx)
	if err != nil {
		return nil, "", err
	}
	defer reader.Close()

	originalImage, format, err := image.Decode(reader)
	if err != nil {
		return nil, "", err
	}
	return originalImage, format, nil
}

func resizeImage(originalImage image.Image, r *http.Request) (image.Image, error) {
	widthParam := r.URL.Query().Get("w")
	heightParam := r.URL.Query().Get("h")

	if widthParam == "" && heightParam == "" {
		return originalImage, nil
	}

	var width, height int

	if widthParam != "" {
		intWidth, err := strconv.Atoi(widthParam)
		if err != nil || intWidth <= 0 {
			return nil, fmt.Errorf("Invalid width parameter")
		}
		width = intWidth
	}

	if heightParam != "" {
		intHeight, err := strconv.Atoi(heightParam)
		if err != nil || intHeight <= 0 {
			return nil, fmt.Errorf("Invalid height parameter")
		}
		height = intHeight
	}

	if width == 0 && height == 0 {
		return nil, fmt.Errorf("Either width or height must be specified")
	}

	srcBounds := originalImage.Bounds()
	srcW, srcH := srcBounds.Max.X, srcBounds.Max.Y

	if width == 0 {
		width = int(float64(srcW) / float64(srcH) * float64(height))
	} else if height == 0 {
		height = int(float64(srcH) / float64(srcW) * float64(width))
	}

	dst := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.ApproxBiLinear.Scale(dst, dst.Bounds(), originalImage, srcBounds, draw.Over, nil)

	return dst, nil
}

func encodeImage(w http.ResponseWriter, img image.Image, format string) error {
	switch format {
	case "png":
		return png.Encode(w, img)
	case "jpeg":
		return jpeg.Encode(w, img, nil)
	case "gif":
		return gif.Encode(w, img, nil)
	default:
		return fmt.Errorf("Unsupported image format: %s", format)
	}
}

func writeToResponse(w http.ResponseWriter, resizedImage image.Image, originalFormat string) error {
	w.Header().Set("Content-Type", "image/"+originalFormat)
	return encodeImage(w, resizedImage, originalFormat)
}
