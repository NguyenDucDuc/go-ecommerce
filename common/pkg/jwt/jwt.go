package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Roles []string `json:"roles"`
	jwt.RegisteredClaims
}

type JwtService struct {
	jwtSecret  string
	jwtAccessExp int
	jwtRefreshExp int
	jwtIssuer string
}

func NewJWTService(jwtSecret string, jwtAccessExp int, jwtRefreshExp int, jwtIssuer string) IJwtService {
	return &JwtService{
		jwtSecret:  jwtSecret,
		jwtAccessExp: jwtAccessExp,
		jwtRefreshExp: jwtRefreshExp,
		jwtIssuer: jwtIssuer,
	}
}
// 1. Generate Access Token (Thường sống ngắn: 15-60 phút)
func (s *JwtService) GenerateAccessToken(userID string, email string, roles []string) (string, error) {
	expirationTime := time.Now().Add(time.Duration(s.jwtAccessExp) * 24 * time.Hour)
	claims := &CustomClaims{
		UserID: userID,
		Email:   email,
		Roles: roles,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    s.jwtIssuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

// 2. Generate Refresh Token (Thường sống dài: 7-30 ngày)
// Refresh Token thường ít thông tin hơn để bảo mật
func (s *JwtService) GenerateRefreshToken(userID string) (string, error) {
	expirationTime := time.Now().Add(time.Duration(s.jwtRefreshExp) * 24 * time.Hour)
	claims := &jwt.RegisteredClaims{
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(expirationTime),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Issuer:    s.jwtIssuer,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

// 3. Hàm Validate Token (Dùng cho Middleware)
func (s *JwtService) ValidateToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Kiểm tra thuật toán ký
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}