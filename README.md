# IMAGE PROCESSING SERVICE

A backend service for processing and managing images with basic transformations. Designed to handle low to medium workloads, this service integrates seamlessly with AWS S3 and supports modern API practices.

## Features
* Transformations:
  * Resize
  * Crop
  * Rotate
  * Flip
  * Filters
    * Modulation
    * Gamma
    * Gaussian Blur
    * Grayscale
    * Sharpness
    * Format (in dev)
* Storage:
  Upload images to AWS S3 with optional metadata and manage storage efficiently.
* OpenAPI:
  * Fully documented APIs for easy integration.

## Tech Stack
* Backend
  * Go for image processing backend service
* Storage
  * AWS S3 for reliable and scalable file storage
* Image Processing
  * Leveraging the ```govips``` library for faster and memory-efficient operations
 
## API Documentation
OpenAPI documentation is available at ```/swagger/index.html``` once the service is running

## Setup Instructions
### Prerequisite
  * Go 1.18+ installed.
  * AWS account with S3 access

### Clone the Repository
```bash
git clone https://github.com/yourusername/image-processing-service.git
cd image-processing-service
```
### Environment Setup
 1. Create a ```.env``` file in the project root:
    ```bash
    AWS_ACCESS_KEY_ID=your-access-key
    AWS_SECRET_ACCESS_KEY=your-secret-key
    S3_BUCKET_NAME=your-bucket-name
    S3_REGION=your-region
    ```
## Usage
### Run the service
```bash
go run cmd/main.go
```
### Endpoints
#### Upload Image
* POST ```/images/upload```
  * Upload an image to AWS S3 bucket.
  * Request body: ```multipart/form-data```
#### Get Image
* GET ```images?key=image.jpg```
  * Fetch the required image if available
#### List Images
* GET ```images/all?page=2&limit=15```
  * Fetch a paginated list of images
### Transform Image
* POST ```/images/transform```
* Request body:
  ```json
  {
  "image": "image-key-in-s3",
  "output": "output-image-name",
  "transformation": {
    "resize": { "width": 200, "height": 200 },
    "crop": { "x": 10, "y": 10, "width": 100, "height": 100 },
    "rotate": { "angle": 90 },
    "filters": { "brightness": 1.2, "contrast": 1.0 }
    }
  }
  ```

## Future Scope
* Scalable Architecture:
  * Microservices-based architecture with a Gateway and OAuth2 for authentication.
* Performance Optimization
* User Account Management
* Front-End (e.g React, HTMX)/Mobile App (e.g iOS)
* GraphQL based requests
* Implement CI/CD pipeline
* Expand Cloud Support
* Advanced Transformations
