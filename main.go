package main

import (
	"ayana/db"
	"ayana/routes"
	utilsEnv "ayana/utils/env"
	"log"
	"os"

	"github.com/gin-contrib/cors"
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

	r.Use(cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			return true // Mengizinkan semua domain yang mengirim request
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Authorization"},
		AllowCredentials: true, // Harus true agar bisa mengirim token/cookie
	}))

	// 🔹 Setup routes untuk berbagai fitur aplikasi
	log.Println("📌 Setting up routes...")

	routes.SetupAuthRouter(r)
	routes.SetupHomeRouter(r)
	routes.SetupReservationRouter(r)
	routes.SetupInfoRouter(r)
	routes.SetupProjectRouter(r)
	routes.SetupWeeklyProgressRouter(r)
	routes.SetupCashFlowRouter(r)
	routes.SetupPayoutRouter(r)
	routes.SetupCompanyRouter(r)
	routes.SetupGoodRouter(r)
	routes.SetupAccountRouter(r)
	routes.SetupPayinRouter(r)
	routes.SetupTransactionCategoryRouter(r)
	routes.SetupJournalEntryRouter(r)

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
