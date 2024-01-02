package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"gopkg.in/gographics/imagick.v2/imagick"
)

func main() {
	imagick.Initialize()
	defer imagick.Terminate()

	http.HandleFunc("/resize", resizeHandler)
	fmt.Println("Server is running at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func errorResponse(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func resizeHandler(w http.ResponseWriter, r *http.Request) {
	widthStr := r.URL.Query().Get("w")
	heightStr := r.URL.Query().Get("h")

	var width, height int
	var err error

	if widthStr != "" {
		width, err = strconv.Atoi(widthStr)
		if err != nil {
			errorResponse(w, "Invalid width parameter", http.StatusBadRequest)
			return
		}
	}

	if heightStr != "" {
		height, err = strconv.Atoi(heightStr)
		if err != nil {
			errorResponse(w, "Invalid height parameter", http.StatusBadRequest)
			return
		}
	}

	imgData, err := os.ReadFile("test/test_image.jpg")
	if err != nil {
		errorResponse(w, "Failed to read image file", http.StatusInternalServerError)
		return
	}

	resizedImage, err := resizeImage(imgData, width, height)
	if err != nil {
		errorResponse(w, "Failed to resize image", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Cache-Control", "public, max-age=86400")
	w.Write(resizedImage)
}

func resizeImage(data []byte, width, height int) ([]byte, error) {
	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	if err := mw.ReadImageBlob(data); err != nil {
		return nil, err
	}

	if width == 0 && height == 0 {
		return data, nil
	}

	originalWidth := mw.GetImageWidth()
	originalHeight := mw.GetImageHeight()

	if width == 0 {
		aspectRatio := float64(originalWidth) / float64(originalHeight)
		width = int(float64(height) * aspectRatio)
	} else if height == 0 {
		aspectRatio := float64(originalHeight) / float64(originalWidth)
		height = int(float64(width) * aspectRatio)
	}

	if err := mw.ResizeImage(uint(width), uint(height), imagick.FILTER_LANCZOS, 1); err != nil {
		return nil, err
	}

	return mw.GetImageBlob(), nil
}
