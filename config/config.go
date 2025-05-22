package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config chứa các thiết lập cấu hình của ứng dụng
type Config struct {
	Port           string
	MongoURI       string
	MongoDatabase  string
	JWTSecret      string
	JWTExpiryHours int
}

// LoadConfig đọc cấu hình từ .env và trả về đối tượng Config
func LoadConfig() *Config {
	// Load file .env
	err := godotenv.Load()
	if err != nil {
		log.Println("Không tìm thấy file .env, sử dụng biến môi trường hệ thống")
	}

	// Đọc các biến môi trường
	port := getEnv("PORT", "")
	if port == "" {
		log.Fatal("port không được để trống")
	}
	mongoURI := getEnv("MONGO_URI", "mongodb://localhost:27017")
	mongoDatabase := getEnv("MONGO_DATABASE", "")
	if mongoDatabase == "" {
		log.Fatal("Database_name không được để trống")
	}
	jwtSecret := getEnv("JWT_SECRET", "")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET không được để trống")
	}
	jwtExpiryHours, err := strconv.Atoi(getEnv("JWT_EXPIRY_HOURS", "24"))
	if err != nil {
		log.Fatal("Lỗi khi parse JWT_EXPIRY_HOURS: ", err)
	}

	return &Config{
		Port:           port,
		MongoURI:       mongoURI,
		MongoDatabase:  mongoDatabase,
		JWTSecret:      jwtSecret,
		JWTExpiryHours: jwtExpiryHours,
	}
}

// getEnv lấy giá trị biến môi trường, trả về giá trị mặc định nếu không tồn tại
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
