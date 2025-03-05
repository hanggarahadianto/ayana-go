package main

import (
	"ayana/db"
	"ayana/routes"
	utilsEnv "ayana/utils/env"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	log.Println("🔹 Starting Ayana Backend...")

	// 🔹 Setup logging untuk debugging
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetOutput(os.Stdout)

	// 🔹 Load environment variables dari .env
	log.Println("📂 Loading environment variables...")
	configure, err := utilsEnv.LoadConfig(".")
	if err != nil {
		log.Fatalf("❌ ERROR: Tidak bisa memuat environment variables: %v", err)
	}

	// 🔹 Initialize database
	log.Println("📦 Initializing database...")
	db.InitializeDb(&configure)
	log.Println("✅ Database initialized successfully!")

	// 🔹 Set Gin ke mode release (agar lebih cepat di production)
	gin.SetMode(gin.ReleaseMode)

	// 🔹 Buat router baru
	r := gin.Default()

	// 🔹 Middleware CORS untuk menangani request dari frontend
	log.Println("🌍 Setting up CORS middleware...")
	// r.Use(cors.New(cors.Config{
	// 	AllowOrigins:     []string{"*"}, // Bisa gunakan "*" jika ingin allow semua
	// 	AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	// 	AllowHeaders:     []string{"Content-Type", "Authorization"},
	// 	AllowCredentials: true,
	// 	MaxAge:           12 * time.Hour,
	// }))

	// Tambahkan middleware khusus untuk OPTIONS
	r.Use(func(c *gin.Context) {
		if c.Request.Method == "OPTIONS" {
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type")
			c.Header("Access-Control-Allow-Credentials", "true")
			c.AbortWithStatus(200) // Ubah dari 403 ke 200
			return
		}
		c.Next()
	})

	// 🔹 Debugging Middleware: Log setiap request yang masuk
	r.Use(func(c *gin.Context) {
		log.Printf("📥 Incoming Request: %s %s", c.Request.Method, c.Request.URL.Path)
		c.Next()
		log.Printf("📤 Response Status: %d", c.Writer.Status())
	})

	// 🔹 Setup routes untuk berbagai fitur aplikasi
	log.Println("📌 Setting up routes...")
	routes.SetupAuthRouter(r)
	routes.SetupHomeRouter(r)
	routes.SetupReservationRouter(r)
	routes.SetupInfoRouter(r)
	routes.SetupProjectRouter(r)
	routes.SetupWeeklyProgressRouter(r)
	routes.SetupCashFlowRouter(r)

	// 🔹 Route utama (tes apakah server berjalan)
	r.GET("/", func(c *gin.Context) {
		log.Println("🏠 Root endpoint accessed")
		c.JSON(200, gin.H{
			"message": "Welcome to Ayana Backend! 🚀",
		})
	})

	// 🔹 Jalankan server
	serverAddr := "0.0.0.0:" + configure.ServerPort
	log.Printf("🚀 Server berjalan di: http://localhost:%s", configure.ServerPort)
	log.Fatal(r.Run(serverAddr)) // Gunakan `r.Run()` agar CORS tetap bekerja
}
