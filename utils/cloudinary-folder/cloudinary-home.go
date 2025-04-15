package utils

import (
	"ayana/config"
	utilsEnv "ayana/utils/env"
	"context"
	"fmt"

	"log"
	"mime/multipart"

	// "net/url"
	// "strings"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

func UploadToCloudinary(file multipart.File, filePath string) (string, error) {
	configure, err := utilsEnv.LoadConfig(".")
	if err != nil {
		log.Fatal("🚀 Could not load environment variables ", err)
	}

	ctx := context.Background()
	cld, err := config.SetupCloudinary(&configure)
	if err != nil {
		return "", err
	}

	// Upload params tanpa membuat folder baru
	uploadParams := uploader.UploadParams{
		PublicID:     filePath, // Nama file custom
		ResourceType: "image",
		Folder:       config.EnvCloudUploadFolderHome(&configure),
	}

	result, err := cld.Upload.Upload(ctx, file, uploadParams)
	if err != nil {
		return "", err
	}

	return result.SecureURL, nil
}

// func ExtractPublicID(imageURL string) (string, error) {
// 	parsedURL, err := url.Parse(imageURL)
// 	if err != nil {
// 		return "", err
// 	}

// 	// Hapus '/upload/vXYZ/' → ambil bagian setelahnya
// 	segments := strings.Split(parsedURL.Path, "/upload/")
// 	if len(segments) < 2 {
// 		return "", fmt.Errorf("format URL tidak valid")
// 	}

// 	publicID := segments[1]
// 	publicID = strings.TrimPrefix(publicID, "v")   // kadang masih ada v123/
// 	publicID = strings.SplitN(publicID, "/", 2)[1] // buang version
// 	publicID, _ = url.QueryUnescape(publicID)      // decode %20 → spasi

// 	// Hapus ekstensi ganda .png.png jika perlu
// 	publicID = strings.ReplaceAll(publicID, ".png.png", ".png")

// 	return publicID, nil
// }

// type CloudinaryConfig struct {
// 	CloudName string
// 	ApiKey    string
// 	ApiSecret string
// }

func DeleteFromCloudinary(publicID string) error {
	// Load .env config
	env, err := utilsEnv.LoadConfig(".")
	if err != nil {
		return fmt.Errorf("❌ Gagal load .env: %v", err)
	}

	// Setup Cloudinary client
	cld, err := config.SetupCloudinary(&env)
	if err != nil {
		return fmt.Errorf("❌ Gagal setup Cloudinary: %v", err)
	}

	// Hapus gambar dari Cloudinary
	invalidate := true
	_, err = cld.Upload.Destroy(context.Background(), uploader.DestroyParams{
		PublicID:     publicID,
		ResourceType: "image",

		Invalidate: &invalidate,
	})
	if err != nil {
		return fmt.Errorf("❌ Gagal hapus gambar dari Cloudinary: %v", err)
	}

	fmt.Println("✅ Berhasil hapus dari Cloudinary:", publicID)
	return nil
}
