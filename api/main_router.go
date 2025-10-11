package handler // PENTING: Diubah dari 'main'

import (
	"fmt"
	"net/http"

	// Pastikan PATH PACKAGE Anda (valeth/handler, valeth/model, valeth/utils) sudah benar
	valethhandler "valeth/handler"
	"valeth/model"
	"valeth/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor" // Wajib untuk Vercel
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Variabel global untuk instance Fiber dan koneksi DB.
var fiberApp *fiber.App
var db *gorm.DB

// init() menjalankan inisialisasi yang berat hanya saat cold start.
func init() {
    // Menghilangkan logika main() dan app.Listen()
    app := fiber.New()
    
    fmt.Println("Initializing Fiber application for Vercel...")

    // 1. Koneksi Database
    // CATATAN: PENTING UNTUK MENGGUNAKAN os.Getenv("NAMA_VARIABEL") DI LINGKUNGAN PRODUKSI!
    dsn := "postgresql://postgres.xrsnptveunsdsnfcvxjz:yipikaye2123@aws-1-ap-southeast-1.pooler.supabase.com:5432/postgres"
    var errConnectDb error

    db, errConnectDb = gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if errConnectDb != nil {
        fmt.Println("Failed to connect to database!")
        // Menggunakan panic untuk menghentikan inisialisasi jika DB gagal terhubung
        panic("Database connection error: " + errConnectDb.Error()) 
    }

    // 2. Auto Migrate & Assign DB
    db.AutoMigrate(&model.NasaolineImage{}, &model.User{})
    fmt.Println("Connected to database and ran migrations!")
    // Mengganti penugasan handler.DB = db; menjadi valethhandler.DB = db;
    valethhandler.DB = db 

    // 3. Setup Routing
    app.Get("/", valethhandler.HandlerHome)
    app.Get("/imgs/:id", valethhandler.HandlerImgDetails)
    app.Get("/imgs/:id/download", utils.AuthMiddleware, valethhandler.Handlerdownloadimg)
    app.Get("/imgs/:id/view", utils.AuthMiddleware, valethhandler.HandlerImgView)
    app.Post("/register", valethhandler.Register)
    app.Post("/login", valethhandler.Login)

    // Simpan instance Fiber yang sudah dikonfigurasi
    fiberApp = app
    fmt.Println("Fiber app fully configured.")
}

// Handler adalah fungsi wajib yang dipanggil oleh Vercel Go Runtime untuk setiap permintaan HTTP.
func Handler(w http.ResponseWriter, r *http.Request) {
    // PENTING: Menyesuaikan Request URI agar Fiber dapat memproses path yang benar
    r.RequestURI = r.URL.String()

    // Mengubah Fiber App menjadi http.HandlerFunc dan melayani request
    adaptor.FiberApp(fiberApp).ServeHTTP(w, r)
}