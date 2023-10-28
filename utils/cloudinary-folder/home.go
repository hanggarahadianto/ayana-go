package utils

import (
	"ayana/config"
	"ayana/utils"
	"context"
	"log"
	"mime/multipart"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

func UploadtoHomeFolder(file multipart.File, filePath string) (string, error) {
	configure, err := utils.LoadConfig(".")
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
