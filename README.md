# Cloud Function Image Resizer

## Description

This project provides a Google Cloud Function to resize images stored in Google Cloud Storage. The function is written in Go and uses the Imagick library for image manipulation. Users can specify the desired width and height as URL query parameters to dynamically resize images.

## Deployment

To deploy this function to Google Cloud, you can use the `gcloud` command-line tool. Below is the sample command to deploy the function:

```bash
gcloud functions deploy contents-cdn --entry-point ResizeImage --runtime go118 --trigger-http --allow-unauthenticated --project your-project-id --region your-region --set-env-vars GCS_BUCKET_NAME="your-bucket-name"
```

Replace `your-project-id` with the Google Cloud Project ID, `your-region` with the Google Cloud region, and `your-bucket-name` with the name of the Google Cloud Storage bucket where your images are stored.

## Usage

After deploying the function, you can access it by sending HTTP GET requests. Include the image name and optionally the desired width and height as query parameters.

Example:

```
https://REGION-PROJECT_ID.cloudfunctions.net/contents-cdn?image_name=image.jpg&w=200&h=300
```

Replace `REGION` with your Google Cloud region (e.g., `asia-northeast1`) and `PROJECT_ID` with your Google Cloud project ID.
