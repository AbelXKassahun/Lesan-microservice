package app

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"

	"lesson-service/internal/infra/storage"
)

type SignedURLResponse struct {
	URL string `json:"url"`
}

var minioClient = storage.InitMinio()

// Generate a signed URL for uploading a file
func UploadObjectService(bucket, objectName string) (string, error) { // objectName is basically the file name
	ctx := context.Background()
	// bucket := "user-submissions"

	// Example: get objectName from query (?file=myfile.mp3)
	// objectName := r.URL.Query().Get("file")
	if objectName == "" || bucket == "" {
		// http.Error(w, "file query param required", http.StatusBadRequest)
		log.Println("missing parametter, either bucket or object name")
		return "", fmt.Errorf("missing parametter, either bucket or object name")
	}

	// Ensure bucket exists
	err := minioClient.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
	if err != nil {
		resp := minio.ToErrorResponse(err)
		if resp.Code != "BucketAlreadyOwnedByYou" && resp.Code != "BucketAlreadyExists" {
			log.Println("bucket creation failed:", err)
			// http.Error(w, "could not create bucket", http.StatusInternalServerError)
			return "", fmt.Errorf("bucket creation failed: %v", err)
		}
	}

	// Generate presigned PUT URL (valid 5 minutes)
	presignedURL, err := minioClient.PresignedPutObject(ctx, bucket, objectName, 5*time.Minute)
	if err != nil {
		// http.Error(w, "could not generate signed url", http.StatusInternalServerError)
		log.Println("could not generate signed url")
		return "", fmt.Errorf("could not generate signed url: %v", err) 
	}

	// json.NewEncoder(w).Encode(SignedURLResponse{URL: presignedURL.String()})
	return presignedURL.String(), nil
}

// Generate a signed URL for fetching a file
func FetchObjectService(bucket, objectName string) (string, error) {
	ctx := context.Background()
	// bucket := "exercises-audio"

	// Example: get objectName from query (?file=lesson1.mp3)
	// objectName := r.URL.Query().Get("file")
	if objectName == "" || bucket == "" {
		// http.Error(w, "file query param required", http.StatusBadRequest)
		log.Println("missing parametter, either bucket or object name")
		return "", fmt.Errorf("missing parametter, either bucket or object name")
	}

	// Generate presigned GET URL (valid 5 minutes)
	reqParams := make(url.Values)
	presignedURL, err := minioClient.PresignedGetObject(ctx, bucket, objectName, 5*time.Minute, reqParams)
	if err != nil {
		// http.Error(w, "could not generate signed url", http.StatusInternalServerError)
		log.Println("could not generate signed url")
		return "", fmt.Errorf("could not generate signed url: %v", err)
	}

	// json.NewEncoder(w).Encode(SignedURLResponse{URL: presignedURL.String()})
	return presignedURL.String(), nil
}
