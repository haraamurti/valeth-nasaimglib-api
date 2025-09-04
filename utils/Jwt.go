package Utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Get JWT secret from environment variable, fallback to default for development
func getJWTSecret() []byte {
    secret := os.Getenv("JWT_SECRET")
    if secret == "" {
        // For development only - in production, ALWAYS use environment variable
        secret = "nasa-image-library-secret-key-2025"
    }
    return []byte(secret)
}

// Claims struct defines what information we store in the JWT token
type Claims struct {
    UserID uint   `json:"user_id"`
    Email  string `json:"email"`
    jwt.RegisteredClaims
}

// GenerateToken creates a new JWT token for a user
func GenerateToken(userID uint, email string) (string, error) {
    // Create claims with user information and expiration
    claims := Claims{
        UserID: userID,
        Email:  email,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // Token expires in 24 hours
            IssuedAt:  jwt.NewNumericDate(time.Now()),                     // When token was created
            NotBefore: jwt.NewNumericDate(time.Now()),                     // Token valid from now
            Issuer:    "nasa-image-library",                               // Who issued this token
            Subject:   email,                                              // Who this token is for
        },
    }

    // Create token with claims
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    
    // Sign token with secret key
    return token.SignedString(getJWTSecret())
}

// ValidateToken validates and parses a JWT token
func ValidateToken(tokenString string) (*Claims, error) {
    // Parse token with claims
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        // Verify signing method
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, jwt.ErrSignatureInvalid
        }
        return getJWTSecret(), nil
    })

    if err != nil {
        return nil, err
    }

    // Extract and return claims if token is valid
    if claims, ok := token.Claims.(*Claims); ok && token.Valid {
        return claims, nil
    }

    return nil, jwt.ErrSignatureInvalid
}

// RefreshToken generates a new token for an existing valid token
func RefreshToken(tokenString string) (string, error) {
    // Validate current token
    claims, err := ValidateToken(tokenString)
    if err != nil {
        return "", err
    }

    // Generate new token with same user info
    return GenerateToken(claims.UserID, claims.Email)
}