package utils

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware validates JWT tokens
func AuthMiddleware(c *fiber.Ctx) error {
    // 1. Get Authorization header
    authHeader := c.Get("Authorization")
    if authHeader == "" {
        return c.Status(401).JSON(fiber.Map{
            "error": "Authorization required",
            "details": "Please provide Authorization header with Bearer token",
        })
    }

    // 2. Check if it starts with "Bearer "
    if !strings.HasPrefix(authHeader, "Bearer ") {
        return c.Status(401).JSON(fiber.Map{
            "error": "Invalid authorization format",
            "details": "Authorization header must be: Bearer <token>",
        })
    }

    // 3. Extract token (remove "Bearer " prefix)
    tokenString := strings.TrimPrefix(authHeader, "Bearer ")
    if tokenString == "" {
        return c.Status(401).JSON(fiber.Map{
            "error": "Missing token",
            "details": "No token provided after Bearer",
        })
    }

    // 4. Parse and validate token
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        // Validate signing method
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, jwt.ErrSignatureInvalid
        }
        return []byte(getJWTSecret()), nil // Use your same secret key from jwt.go
    })

    if err != nil || !token.Valid {
        return c.Status(401).JSON(fiber.Map{
            "error": "Invalid token",
            "details": "Token is expired, invalid, or malformed",
        })
    }

    // 5. ✅ FIXED: Properly extract claims
    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok {
        return c.Status(401).JSON(fiber.Map{
            "error": "Invalid token claims",
            "details": "Unable to parse token claims",
        })
    }

    // 6. ✅ FIXED: Safely access claims with type assertion
    userID, userIDOk := claims["user_id"]
    email, emailOk := claims["email"]
    
    if !userIDOk || !emailOk {
        return c.Status(401).JSON(fiber.Map{
            "error": "Missing token data",
            "details": "Token missing required user information",
        })
    }

    // 7. Add user info to context for use in handlers
    c.Locals("user_id", userID)
    c.Locals("email", email)

    // 8. Continue to next handler
    return c.Next()
}