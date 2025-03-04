package main

import (
	"ayana/db"
	"ayana/routes"
	utilsEnv "ayana/utils/env"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rs/cors"
)

func main() {
	// 🔹 Setup logging untuk debugging
	log.SetFlags(0)
	log.SetOutput(os.Stdout)

	// 🔹 Load environment variables dari .env
	configure, err := utilsEnv.LoadConfig(".")
	if err != nil {
		log.Fatal("🚀 ERROR: Tidak bisa memuat environment variables:", err)
	}

	// 🔹 Pastikan CLIENT_ORIGIN terbaca
	clientOrigin := configure.ClientOrigin
	if clientOrigin == "" {
		clientOrigin = os.Getenv("CLIENT_ORIGIN") // Coba baca dari ENV
	}
	if clientOrigin == "" {
		log.Fatal("❌ ERROR: CLIENT_ORIGIN tidak diset. Periksa file .env atau environment variables.")
	}

	log.Println("✅ CLIENT_ORIGIN berhasil dimuat:", clientOrigin)

	// 🔹 Initialize database
	db.InitializeDb(&configure)

	// 🔹 Set Gin ke mode release (agar lebih cepat di production)
	gin.SetMode(gin.ReleaseMode)

	// 🔹 Buat router baru
	r := gin.Default()

	// 🔹 Middleware CORS untuk menangani request dari frontend
	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{clientOrigin}, // Gunakan nilai dari environment
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	// 🔹 Debugging: Cek header response yang dikirim
	r.Use(func(c *gin.Context) {
		corsMiddleware.HandlerFunc(c.Writer, c.Request)
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	r.Use(func(c *gin.Context) {
		c.Next()
		log.Println("🚀 Response Headers:", c.Writer.Header())
	})

	// 🔹 Setup routes untuk berbagai fitur aplikasi
	routes.SetupAuthRouter(r)
	routes.SetupHomeRouter(r)
	routes.SetupReservationRouter(r)
	routes.SetupInfoRouter(r)
	routes.SetupProjectRouter(r)
	routes.SetupWeeklyProgressRouter(r)
	routes.SetupCashFlowRouter(r)

	// 🔹 Route utama (tes apakah server berjalan)
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to my Ayana application! 🚀",
		})
		log.Println("Welcome to my Ayana application! 🚀")
	})

	// 🔹 Jalankan server
	serverAddr := "0.0.0.0:" + configure.ServerPort
	log.Println("🚀 Server berjalan di port:", configure.ServerPort)
	log.Fatal(r.Run(serverAddr)) // Gunakan `r.Run()` agar CORS tetap bekerja
}
