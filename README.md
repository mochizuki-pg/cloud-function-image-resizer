# Cloud Function Image Resizer

## Description

This project provides a Google Cloud Function to resize images stored in Google Cloud Storage. The function is written in Go and uses the Imagick library for image manipulation. Users can specify the desired width and height as URL query parameters to dynamically resize images.

## Deployment

To deploy this function to Google Cloud, you can use the `gcloud` command-line tool. Below is the sample command to deploy the function:

```bash
gcloud functions deploy image-resizer --gen2 --entry-point ResizeImage --runtime go121 --trigger-http --allow-unauthenticated --project your-project-id --region your-region --set-env-vars GCS_BUCKET_NAME="your-bucket-name"
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


### Caching

The Cloud Function has been designed to support caching by default. When an image is successfully resized and served, the HTTP response will include a `Cache-Control` header. This header indicates that the resource can be publicly cached and specifies the duration for which the resource is considered fresh.

By default, the cache duration is set to 24 hours (`max-age=86400`). This means that once the image is fetched and resized, CDNs, browsers, or any intermediate cache servers can store and reuse the resized image for up to 24 hours without re-fetching it from the Cloud Function.

### Local Testing
To test the image resizing functionality locally, you can run the local server using the go run command. This will start a local HTTP server, allowing you to test the resizing features without deploying to Google Cloud.

Run the local server using the following command:

```bash
go run test/local_server.go
```

Once the server is running, you can test the image resizing by accessing the following URL endpoints with your web browser or HTTP client:

Local Server Endpoints
Resize with both width and height specified:

```bash
http://localhost:8080/resize?w=200&h=300
```
<br> 
<p align="center">
  <img src="https://github.com/mochizuki-pg/cloud-function-image-resizer/assets/38402160/cd47aaf9-5847-4433-bc24-eb73cddb844c" width=400px>
</p>
<br>
The server will resize the image located at `test/test_image.jpg` according to the specified w and h query parameters and return the resized image in the response.
