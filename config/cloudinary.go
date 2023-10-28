package config

import (
	"ayana/utils"

	"github.com/cloudinary/cloudinary-go/v2"
)

func SetupCloudinary(config *utils.Config) (*cloudinary.Cloudinary, error) {
	cldSecret := config.CLOUDINARY_API_SECRET
	cldName := config.CLOUDINARY_CLOUD_NAME
	cldKey := config.CLOUDINARY_API_KEY

	cld, err := cloudinary.NewFromParams(cldName, cldKey, cldSecret)
	if err != nil {
		return nil, err
	}

	return cld, nil
}

// func SetupCloudinary() (*cloudinary.Cloudinary, error) {
// 	cldSecret := os.Getenv("CLOUDINARY_API_SECRET")
// 	cldName := os.Getenv("CLOUDINARY_CLOUD_NAME")
// 	cldKey := os.Getenv("CLOUDINARY_API_KEY")

// 	cld, err := cloudinary.NewFromParams(cldName, cldKey, cldSecret)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return cld, nil
// }

func EnvCloudUploadFolderHome(config *utils.Config) string {
	return config.CLOUDINARY_HOME_FOLDER

	// return homeFolder, nil
	// // return os.Getenv("CLOUDINARY_HOME_FOLDER")
}

// func EnvCloudUploadFolderHome() string {
// 	err := godotenv.Load(".env")
// 	if err != nil {
// 		log.Fatal("Error loading app.env file")
// 	}
// 	return os.Getenv("CLOUDINARY_HOME_FOLDER")
// }
