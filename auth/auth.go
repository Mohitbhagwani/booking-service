// auth.go

package auth

import (
    "net/http"
    "github.com/dgrijalva/jwt-go"
    "github.com/google/uuid"
    "time"
)

// Define your JWT secret key
var jwtSecret = []byte("booking")

// GenerateJWT generates a new JWT token with user claims and a specified expiration time (in seconds)
func GenerateJWT(userID uuid.UUID, userRoles []string, expirationSeconds int64) (string, error) {
    // Calculate the expiration time
    expirationTime := time.Now().Add(time.Second * time.Duration(expirationSeconds))

    // Create a new JWT token
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id": userID.String(), // Convert uuid.UUID to string
        "roles":   userRoles,
        "exp":     expirationTime.Unix(),
    })

    // Sign the token with the secret key
    tokenString, err := token.SignedString(jwtSecret)
    if err != nil {
        return "", err
    }

    return tokenString, nil
}

// ValidateTokenMiddleware is middleware for validating JWT tokens
func ValidateTokenMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Extract the JWT token from the Authorization header
        tokenString := extractTokenFromRequest(r)

        if tokenString == "" {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        // Parse the JWT token
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            return jwtSecret, nil
        })

        if err != nil || !token.Valid {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
		claims, ok := token.Claims.(jwt.MapClaims)
        if !ok {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

		exp, ok := claims["exp"].(float64)
        if !ok {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        // Convert expiration time to Unix timestamp
        expirationTime := time.Unix(int64(exp), 0)

        // Check if the token has expired
        if time.Now().After(expirationTime) {
            http.Error(w, "Token has expired", http.StatusUnauthorized)
            return
        }

        // Token is valid; proceed to the next handler
        next.ServeHTTP(w, r)
    })
}

func extractTokenFromRequest(r *http.Request) string {
    // Retrieve the token from the Authorization header
    token := r.Header.Get("Authorization")
    if token != "" {
        // Check if the header has the "Bearer " prefix and remove it
        if len(token) > 7 && token[:7] == "Bearer " {
            return token[7:]
        }
    }
    return ""
}
