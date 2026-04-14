package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// https://leapcell.io/blog/secure-your-apis-with-jwt-authentication-in-gin-middleware

const UserIDKey = "userID" // Key used to store user ID in Gin context to minimize typos

func RequireAuth() gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		authHeader := ginContext.GetHeader("Authorization")
		if authHeader == "" {
			ginContext.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization header is required"})
			return
		}

		// split the header by space (with max 2 parts) to get the token part. The expected format is "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			ginContext.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization header must be in the format: Bearer <token>"})
			return
		}

		tokenString := parts[1]

		// Parse and validate the JWT token
		token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
			// Ensure the signing method is HMAC and return the secret key for validation
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(os.Getenv("JWT_SECRET_KEY")), nil
		})

		// If there was an error parsing the token or if the token is invalid, return an unauthorized error
		if err != nil || !token.Valid {
			ginContext.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}

		// Extract the user ID from the token claims and store it in the Gin context for handlers to access
		claims, ok := token.Claims.(*jwt.RegisteredClaims)
		if !ok || claims.Subject == "" {
			ginContext.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
			return
		}

		// Store the user ID in the Gin context using a constant key to avoid typos and make it easier for handlers to access
		ginContext.Set(UserIDKey, claims.Subject)
		// Call the next handler in the chain
		ginContext.Next()
	}
}
