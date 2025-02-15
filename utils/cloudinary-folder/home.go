package utils

import (
	"ayana/config"
	utilsEnv "ayana/utils/env"
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"strings"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

func UploadtoHomeFolder(file multipart.File, filePath string) (string, error) {
	configure, err := utilsEnv.LoadConfig(".")
	if err != nil {
		log.Fatal("🚀 Could not load environment variables ", err)
	}
	ctx := context.Background()
	cld, err := config.SetupCloudinary(&configure)
	if err != nil {
		return "no background context", err
	}

	// create upload params
	uploadParams := uploader.UploadParams{
		PublicID:     filePath,
		ResourceType: "image",
		Folder:       config.EnvCloudUploadFolderHome(&configure),
	}

	result, err := cld.Upload.Upload(ctx, file, uploadParams)
	if err != nil {
		return "", err
	}

	imageUrl := result.SecureURL
	return imageUrl, nil
}

func extractPublicID(imageURL string) string {
	parts := strings.Split(imageURL, "/")
	if len(parts) < 2 {
		return ""
	}

	// Get the last part and remove the file extension
	filename := parts[len(parts)-1]
	publicID := strings.TrimSuffix(filename, ".jpg") // Adjust based on your file types
	publicID = strings.TrimSuffix(publicID, ".png")
	publicID = strings.TrimSuffix(publicID, ".jpeg")

	// Get the full path from the base folder
	publicID = "ayana/home/" + publicID

	fmt.Println("🆔 Extracted Public ID:", publicID)
	return publicID
}

func DeleteFromCloudinary(image string) error {
	fmt.Println("🛑 Deleting image from Cloudinary:", image)

	// Extract public ID from URL
	publicID := extractPublicID(image)
	if publicID == "" {
		return fmt.Errorf("❌ Failed to extract public ID from URL")
	}

	fmt.Println("PublicID:", publicID)

	// Load Cloudinary configuration
	configure, err := utilsEnv.LoadConfig(".")
	if err != nil {
		return fmt.Errorf("🚀 Could not load environment variables: %v", err)
	}

	// Initialize Cloudinary
	cld, err := config.SetupCloudinary(&configure)
	if err != nil {
		return fmt.Errorf("❌ Cloudinary setup error: %v", err)
	}

	// Delete the image
	invalidate := true
	_, err = cld.Upload.Destroy(context.Background(), uploader.DestroyParams{
		PublicID:     publicID,
		ResourceType: "image",
		Invalidate:   &invalidate, // Pass a pointer to a boolean
	})
	if err != nil {
		return fmt.Errorf("❌ Failed to delete from Cloudinary: %v", err)
	}

	fmt.Println("✅ Successfully deleted image from Cloudinary:", publicID)
	return nil
}
