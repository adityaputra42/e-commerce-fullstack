package config

import (
	"log"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
	JWT      JWTConfig
	CORS     CORSConfig
	SMTP     SMTPConfig
	System   SystemConfig
	Supabase SupabaseConfig
	// RunMigrationsOnBoot mengontrol apakah database.Migrate()/SeedDatabase()
	// dijalankan otomatis setiap kali aplikasi start. Set
	// RUN_MIGRATIONS_ON_BOOT=false di production dan jalankan
	// `go run cmd/migrate/main.go` sebagai deployment step eksplisit.
	RunMigrationsOnBoot bool
}

type SupabaseConfig struct {
	Url string
	Key string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
	// TimeZone dipakai untuk koneksi DB (mis. "Asia/Jakarta"). Kosong berarti
	// pakai default Asia/Jakarta (lihat database.Connect).
	TimeZone string
}

type ServerConfig struct {
	Port string
	Host string
	Env  string
}

type JWTConfig struct {
	Secret             string
	RefreshSecret      string
	AccessTokenExpiry  time.Duration
	RefreshTokenExpiry time.Duration
}

type CORSConfig struct {
	AllowedOrigins []string
}

type SMTPConfig struct {
	Host      string
	Port      string
	Username  string
	Password  string
	FromEmail string
}

type SystemConfig struct {
	DefaultAdminEmail    string
	DefaultAdminPassword string
	// FrontendResetPasswordURL adalah halaman di frontend yang menerima
	// ?token=... dan menampilkan form set-password-baru. Dipakai untuk
	// menyusun link di email forgot-password.
	FrontendResetPasswordURL string
}

func Load() *Config {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Warning: Could not load .env file: %v", err)
	}

	accessTokenExpiry, _ := time.ParseDuration(viper.GetString("JWT_ACCESS_TOKEN_EXPIRY"))
	refreshTokenExpiry, _ := time.ParseDuration(viper.GetString("JWT_REFRESH_TOKEN_EXPIRY"))

	allowedOrigins := strings.Split(viper.GetString("CORS_ALLOWED_ORIGINS"), ",")
	for i := range allowedOrigins {
		allowedOrigins[i] = strings.TrimSpace(allowedOrigins[i])
	}

	// Default true kalau env var belum diset sama sekali, supaya alur
	// development lokal (go run cmd/main.go) tetap jalan seperti sebelumnya.
	runMigrationsOnBoot := true
	if viper.IsSet("RUN_MIGRATIONS_ON_BOOT") {
		runMigrationsOnBoot = viper.GetBool("RUN_MIGRATIONS_ON_BOOT")
	}

	return &Config{
		RunMigrationsOnBoot: runMigrationsOnBoot,
		Database: DatabaseConfig{
			Host:     viper.GetString("DB_HOST"),
			Port:     viper.GetString("DB_PORT"),
			User:     viper.GetString("DB_USER"),
			Password: viper.GetString("DB_PASSWORD"),
			Name:     viper.GetString("DB_NAME"),
			SSLMode:  viper.GetString("DB_SSL_MODE"),
			TimeZone: viper.GetString("DB_TIMEZONE"),
		},
		Server: ServerConfig{
			Port: viper.GetString("PORT"),
			Host: viper.GetString("HOST"),
			Env:  viper.GetString("ENV"),
		},
		JWT: JWTConfig{
			Secret:             viper.GetString("JWT_SECRET"),
			RefreshSecret:      viper.GetString("JWT_REFRESH_SECRET"),
			AccessTokenExpiry:  accessTokenExpiry,
			RefreshTokenExpiry: refreshTokenExpiry,
		},
		CORS: CORSConfig{
			AllowedOrigins: allowedOrigins,
		},
		SMTP: SMTPConfig{
			Host:      viper.GetString("SMTP_HOST"),
			Port:      viper.GetString("SMTP_PORT"),
			Username:  viper.GetString("SMTP_USERNAME"),
			Password:  viper.GetString("SMTP_PASSWORD"),
			FromEmail: viper.GetString("SMTP_FROM_EMAIL"),
		},
		System: SystemConfig{
			DefaultAdminEmail:        viper.GetString("DEFAULT_ADMIN_EMAIL"),
			DefaultAdminPassword:     viper.GetString("DEFAULT_ADMIN_PASSWORD"),
			FrontendResetPasswordURL: viper.GetString("FRONTEND_RESET_PASSWORD_URL"),
		},
		Supabase: SupabaseConfig{
			Url: viper.GetString("SUPABASE_URL"),
			Key: viper.GetString("SUPABASE_KEY"),
		},
	}
}
