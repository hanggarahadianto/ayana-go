package utils

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DBHost     string `mapstructure:"DB_HOST"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBName     string `mapstructure:"DB_NAME"`
	DBPort     string `mapstructure:"DB_PORT"`
	ServerPort string `mapstructure:"SERVER_PORT"`

	ClientOrigin string `mapstructure:"CLIENT_ORIGIN"`
	JWTSecret    string `mapstructure:"JWT_SECRET"`

	AccessTokenPrivateKey  string        `mapstructure:"ACCESS_TOKEN_PRIVATE_KEY"`
	AccessTokenPublicKey   string        `mapstructure:"ACCESS_TOKEN_PUBLIC_KEY"`
	RefreshTokenPrivateKey string        `mapstructure:"REFRESH_TOKEN_PRIVATE_KEY"`
	RefreshTokenPublicKey  string        `mapstructure:"REFRESH_TOKEN_PUBLIC_KEY"`
	AccessTokenExpiresIn   time.Duration `mapstructure:"ACCESS_TOKEN_EXPIRED_IN"`
	RefreshTokenExpiresIn  time.Duration `mapstructure:"REFRESH_TOKEN_EXPIRED_IN"`
	AccessTokenMaxAge      int           `mapstructure:"ACCESS_TOKEN_MAXAGE"`
	RefreshTokenMaxAge     int           `mapstructure:"REFRESH_TOKEN_MAXAGE"`

	CLOUDINARY_CLOUD_NAME string `mapstructure:"CLOUDINARY_CLOUD_NAME"`
	CLOUDINARY_API_KEY    string `mapstructure:"CLOUDINARY_API_KEY"`
	CLOUDINARY_API_SECRET string `mapstructure:"CLOUDINARY_API_SECRET"`

	CLOUDINARY_HOME_FOLDER string `mapstructure:"CLOUDINARY_HOME_FOLDER"`

	EmailFrom string `mapstructure:"EMAIL_FROM"`
	SMTPHost  string `mapstructure:"SMTP_HOST"`
	SMTPPass  string `mapstructure:"SMTP_PASS"`
	SMTPPort  int    `mapstructure:"SMTP_PORT"`
	SMTPUser  string `mapstructure:"SMTP_USER"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.SetConfigFile(fmt.Sprintf("%s/.env", path))
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("⚠️ Warning: No .env file found at %s/.env. Using default env variables.\n", path)
	} else {
		fmt.Println("✅ Environment variables loaded from .env")
	}

	config.AccessTokenPrivateKey = strings.ReplaceAll(viper.GetString("ACCESS_TOKEN_PRIVATE_KEY"), `\n`, "\n")
	config.AccessTokenPublicKey = strings.ReplaceAll(viper.GetString("ACCESS_TOKEN_PUBLIC_KEY"), `\n`, "\n")

	// ✅ Convert ACCESS_TOKEN_EXPIRED_IN (string) -> time.Duration
	accessTokenDuration, err := time.ParseDuration(viper.GetString("ACCESS_TOKEN_EXPIRED_IN"))
	if err != nil {
		fmt.Println("❌ Error parsing ACCESS_TOKEN_EXPIRED_IN:", err)
		return config, err
	}
	config.AccessTokenExpiresIn = accessTokenDuration

	// ✅ Convert REFRESH_TOKEN_EXPIRED_IN (string) -> time.Duration
	refreshTokenDuration, err := time.ParseDuration(viper.GetString("REFRESH_TOKEN_EXPIRED_IN"))
	if err != nil {
		fmt.Println("❌ Error parsing REFRESH_TOKEN_EXPIRED_IN:", err)
		return config, err
	}
	config.RefreshTokenExpiresIn = refreshTokenDuration

	// ✅ Use GetInt for integer values
	config.AccessTokenMaxAge = viper.GetInt("ACCESS_TOKEN_MAXAGE")
	config.RefreshTokenMaxAge = viper.GetInt("REFRESH_TOKEN_MAXAGE")

	// ✅ Unmarshal ke struct
	if err := viper.Unmarshal(&config); err != nil {
		fmt.Println("❌ Error parsing .env file:", err)
		return config, err
	}

	return config, nil
}
