package config

import (
	utilsEnv "ayana/utils/env"

	"github.com/cloudinary/cloudinary-go/v2"
)

func SetupCloudinary(config *utilsEnv.Config) (*cloudinary.Cloudinary, error) {
	cldSecret := config.CLOUDINARY_API_SECRET
	cldName := config.CLOUDINARY_CLOUD_NAME
	cldKey := config.CLOUDINARY_API_KEY

	cld, err := cloudinary.NewFromParams(cldName, cldKey, cldSecret)
	if err != nil {
		return nil, err
	}

	return cld, nil
}

func EnvCloudUploadFolderHome(config *utilsEnv.Config) string {
	return config.CLOUDINARY_HOME_FOLDER

}
func EnvCloudDeleteFolderHome(config *utilsEnv.Config) string {
	return config.CLOUDINARY_HOME_FOLDER

}
