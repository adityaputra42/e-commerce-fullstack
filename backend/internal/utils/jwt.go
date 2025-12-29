package utils

import (
	"e-commerce/backend/internal/config"
	"e-commerce/backend/internal/models"
	"errors"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService struct {
	config *config.Config
}

func NewJWTService(cfg *config.Config) *JWTService {
	return &JWTService{config: cfg}
}

func (s *JWTService) GenerateAccessToken(user *models.User) (string, time.Time, error) {
	// 🔍 DEBUG: Validasi user sebelum generate token
	if user == nil {
		log.Println("❌ JWT ERROR: user is nil")
		return "", time.Time{}, errors.New("user cannot be nil")
	}

	if user.ID == 0 {
		log.Printf("❌ JWT ERROR: user.ID is 0 for email: %s", user.Email)
		return "", time.Time{}, errors.New("user ID cannot be zero")
	}

	expiresAt := time.Now().Add(s.config.JWT.AccessTokenExpiry)

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role_id": user.RoleID,
		"type":    "access",
		"exp":     expiresAt.Unix(),
		"iat":     time.Now().Unix(),
	}

	// 🔍 DEBUG: Log claims yang akan di-encode
	log.Printf("✅ Generating Access Token - UserID: %d, Email: %s, RoleID: %d",
		user.ID, user.Email, user.RoleID)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.config.JWT.Secret))
	if err != nil {
		log.Printf("❌ JWT ERROR: Failed to sign token: %v", err)
		return "", time.Time{}, err
	}

	log.Printf("✅ Access Token Generated Successfully for UserID: %d", user.ID)
	return tokenString, expiresAt, nil
}

func (s *JWTService) GenerateRefreshToken(user *models.User) (string, time.Time, error) {
	// 🔍 DEBUG: Validasi user sebelum generate token
	if user == nil {
		log.Println("❌ JWT ERROR: user is nil")
		return "", time.Time{}, errors.New("user cannot be nil")
	}

	if user.ID == 0 {
		log.Printf("❌ JWT ERROR: user.ID is 0 for email: %s", user.Email)
		return "", time.Time{}, errors.New("user ID cannot be zero")
	}

	expiresAt := time.Now().Add(s.config.JWT.RefreshTokenExpiry)

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role_id": user.RoleID,
		"type":    "refresh",
		"exp":     expiresAt.Unix(),
		"iat":     time.Now().Unix(),
	}

	// 🔍 DEBUG: Log claims
	log.Printf("✅ Generating Refresh Token - UserID: %d, Email: %s", user.ID, user.Email)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.config.JWT.RefreshSecret))
	if err != nil {
		log.Printf("❌ JWT ERROR: Failed to sign refresh token: %v", err)
		return "", time.Time{}, err
	}

	log.Printf("✅ Refresh Token Generated Successfully for UserID: %d", user.ID)
	return tokenString, expiresAt, nil
}

func (s *JWTService) ValidateAccessToken(tokenString string) (*models.JWTClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.config.JWT.Secret), nil
	})

	if err != nil {
		log.Printf("❌ JWT VALIDATION ERROR: %v", err)
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		tokenType, ok := claims["type"].(string)
		if !ok || tokenType != "access" {
			log.Println("❌ JWT ERROR: Invalid token type")
			return nil, errors.New("invalid token type")
		}

		userID, ok := claims["user_id"].(float64)
		if !ok {
			log.Printf("❌ JWT ERROR: Invalid user_id claim, got: %v (type: %T)", claims["user_id"], claims["user_id"])
			return nil, errors.New("invalid user_id claim")
		}

		// 🔍 DEBUG: Validasi user ID tidak 0
		if uint(userID) == 0 {
			log.Printf("❌ JWT ERROR: user_id is 0 in token claims")
			return nil, errors.New("user_id cannot be zero")
		}

		email, ok := claims["email"].(string)
		if !ok {
			log.Println("❌ JWT ERROR: Invalid email claim")
			return nil, errors.New("invalid email claim")
		}

		roleID, ok := claims["role_id"].(float64)
		if !ok {
			log.Println("❌ JWT ERROR: Invalid role_id claim")
			return nil, errors.New("invalid role_id claim")
		}

		jwtClaims := &models.JWTClaims{
			UserID: uint(userID),
			Email:  email,
			RoleID: uint(roleID),
			Type:   tokenType,
		}

		// 🔍 DEBUG: Log parsed claims
		log.Printf("✅ Token Validated - UserID: %d, Email: %s, RoleID: %d",
			jwtClaims.UserID, jwtClaims.Email, jwtClaims.RoleID)

		return jwtClaims, nil
	}

	log.Println("❌ JWT ERROR: Invalid token")
	return nil, errors.New("invalid token")
}

func (s *JWTService) ValidateRefreshToken(tokenString string) (*models.JWTClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.config.JWT.RefreshSecret), nil
	})

	if err != nil {
		log.Printf("❌ REFRESH TOKEN VALIDATION ERROR: %v", err)
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		tokenType, ok := claims["type"].(string)
		if !ok || tokenType != "refresh" {
			log.Println("❌ JWT ERROR: Invalid refresh token type")
			return nil, errors.New("invalid token type")
		}

		userID, ok := claims["user_id"].(float64)
		if !ok {
			log.Printf("❌ JWT ERROR: Invalid user_id in refresh token, got: %v", claims["user_id"])
			return nil, errors.New("invalid user_id claim")
		}

		// 🔍 DEBUG: Validasi user ID tidak 0
		if uint(userID) == 0 {
			log.Printf("❌ JWT ERROR: user_id is 0 in refresh token claims")
			return nil, errors.New("user_id cannot be zero")
		}

		email, ok := claims["email"].(string)
		if !ok {
			return nil, errors.New("invalid email claim")
		}

		roleID, ok := claims["role_id"].(float64)
		if !ok {
			return nil, errors.New("invalid role_id claim")
		}

		jwtClaims := &models.JWTClaims{
			UserID: uint(userID),
			Email:  email,
			RoleID: uint(roleID),
			Type:   tokenType,
		}

		// 🔍 DEBUG: Log parsed refresh token claims
		log.Printf("✅ Refresh Token Validated - UserID: %d, Email: %s",
			jwtClaims.UserID, jwtClaims.Email)

		return jwtClaims, nil
	}

	return nil, errors.New("invalid token")
}

// Helper function untuk debug token tanpa validasi signature
func (s *JWTService) DebugToken(tokenString string) {
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		log.Printf("❌ DEBUG TOKEN ERROR: %v", err)
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		log.Println("🔍 TOKEN DEBUG INFO:")
		log.Printf("   user_id: %v (type: %T)", claims["user_id"], claims["user_id"])
		log.Printf("   email: %v", claims["email"])
		log.Printf("   role_id: %v", claims["role_id"])
		log.Printf("   type: %v", claims["type"])
		log.Printf("   exp: %v", claims["exp"])
		log.Printf("   iat: %v", claims["iat"])
	}
}
