# Cloud Function Image Resizer

## Description

This project provides a Google Cloud Function to resize images stored in Google Cloud Storage. The function is written in Go and uses the Imagick library for image manipulation. Users can specify the desired width and height as URL query parameters to dynamically resize images.

## Deployment

To deploy this function to Google Cloud, you can use the `gcloud` command-line tool. Below is the sample command to deploy the function:

```bash
gcloud functions deploy image-resizer --entry-point ResizeImage --runtime go118 --trigger-http --allow-unauthenticated --project your-project-id --region your-region --set-env-vars GCS_BUCKET_NAME="your-bucket-name"
```

Replace `your-project-id` with the Google Cloud Project ID, `your-region` with the Google Cloud region, and `your-bucket-name` with the name of the Google Cloud Storage bucket where your images are stored.

## Usage

### With both width and height specified:

To request a resized image with both the width and height specified, you can add `w` and `h` query parameters to the URL:

```
https://REGION-PROJECT_ID.cloudfunctions.net/image-resizer?image_name=image.jpg&w=200&h=300
```

Replace `REGION` with your Google Cloud region (e.g., asia-northeast1) and `PROJECT_ID` with your Google Cloud project ID.

### With either width or height specified:

If you only specify either the width (`w`) or the height (`h`), the function will automatically calculate the other dimension while maintaining the original aspect ratio:

```
// To specify only width:
https://REGION-PROJECT_ID.cloudfunctions.net/image-resizer?image_name=image.jpg&w=200

// To specify only height:
https://REGION-PROJECT_ID.cloudfunctions.net/image-resizer?image_name=image.jpg&h=300
```

Replace `REGION` with your Google Cloud region (e.g., asia-northeast1) and `PROJECT_ID` with your Google Cloud project ID.

### Error Handling

In case of an error, the function will return a JSON response with the following structure:

```json
{
  "error": {
    "code": 500,
    "message": "Fetch error: storage: object doesn't exist",
    "type": "FetchError"
  }
}
```

- `code`: HTTP status code representing the error.
- `message`: A descriptive error message.
- `type`: Type of the error, useful for categorizing errors on the client side.

For example, if the specified image does not exist in the Google Cloud Storage bucket, the function will return a `FetchError` with a 500 status code and a message indicating that the object doesn't exist.
