package jwt

type IJwtService interface {
	GenerateAccessToken(userID string, email string, roles []string) (string, error)
	GenerateRefreshToken(userID string) (string, error)
	ValidateToken(tokenString string) (*CustomClaims, error)
}