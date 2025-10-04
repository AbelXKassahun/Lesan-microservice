package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

// type Claims struct {
// 	UserID string `json:"user_id"`
// 	Tier   string `json:"tier"`
// 	Type   string `json:"type"`
// 	jwt.RegisteredClaims
// }

type Claims struct {
	NameID string `json:"nameid"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// var JWTSecret = []byte(os.Getenv("JWT_SECRET_KEY"))
func getJWTSecret() string {
	myEnv, err := godotenv.Read("../../.env")
	if err != nil {
		log.Fatalln(err)
	}
	return myEnv["JWT_SECRET_KEY"]
}

var JWTSecret = []byte(getJWTSecret())

func Base64URLEncode(data []byte) string {
	return strings.TrimRight(base64.URLEncoding.EncodeToString(data), "=")
}

func Sign(message string, jwtSecret []byte) string {
	h := hmac.New(sha256.New, jwtSecret)
	h.Write([]byte(message))
	return Base64URLEncode(h.Sum(nil))
}

func VerifyJWT(w http.ResponseWriter, r *http.Request, isRefreshToken bool) (bool, *Claims) {
	var tokenString string
	if isRefreshToken {
		var requestBody struct {
			RefreshToken string `json:"refresh_token"`
		}
		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil || requestBody.RefreshToken == "" {
			http.Error(w, "Invalid request: missing refresh_token", http.StatusBadRequest)
			return false, nil
		}
		tokenString = requestBody.RefreshToken
	} else {
		tokenString = GetJWTFromRequest(w, r)
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return JWTSecret, nil
	})
	if err != nil {
		log.Printf("Token validation failed: %v", err)
		http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
		return false, nil
	}

	if !token.Valid {
		http.Error(w, "Token is not valid", http.StatusUnauthorized)
		return false, nil
	}

	// ✅ Verify issuer and audience (matches your C# code)
	if claims.Issuer != "users-service" {
		http.Error(w, "Invalid issuer", http.StatusUnauthorized)
		return false, nil
	}

	if claims.Audience[0] != "lesan-microservices" {
		http.Error(w, "Invalid audience", http.StatusUnauthorized)
		return false, nil
	}

	return true, claims
}

// func VerifyJWT(w http.ResponseWriter, r *http.Request, isRefreshToken bool) (bool, *Claims) {
// 	fmt.Println("@@@JWT secret", JWTSecret)
// 	var tokenString string
// 	if isRefreshToken {
// 		var requestBody struct {
// 			RefreshToken string `json:"refresh_token"`
// 		}
// 		err := json.NewDecoder(r.Body).Decode(&requestBody)
// 		if err != nil || requestBody.RefreshToken == "" {
// 			http.Error(w, "Invalid request: missing refresh_token", http.StatusBadRequest)
// 			return false, nil
// 		}
// 		tokenString = requestBody.RefreshToken
// 	} else {
// 		tokenString = GetJWTFromRequest(w, r)
// 	}

// 	// verify signiture and parse claims
// 	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
// 		// Validate the alg
// 		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
// 		}
// 		// Return the secret key used for signing tokens
// 		return JWTSecret, nil
// 	})

// 	if err != nil {
// 		log.Printf("Token validation failed: %v", err)
// 		if err == jwt.ErrSignatureInvalid {
// 			http.Error(w, "Invalid token signature", http.StatusUnauthorized)
// 			return false, nil
// 		}
// 		http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
// 		return false, nil
// 	}

// 	// verifying standard claims
// 	claims, ok := token.Claims.(*Claims)
// 	if !ok || !token.Valid {
// 		log.Printf("Token claims invalid or token is not valid")
// 		http.Error(w, "Invalid token claims", http.StatusUnauthorized)
// 		return false, nil
// 	}

// 	// [To Do] - uncomment this code if you want to check if the token is blacklisted

// 	jti := claims.ID
// 	log.Println(jti)

// 	return true, claims
// }

func GetJWTFromRequest(w http.ResponseWriter, r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		w.WriteHeader(498)
		w.Write([]byte("authorization header missing"))
		return ""
	}

	const prefix = "Bearer "
	if !strings.HasPrefix(authHeader, prefix) {
		w.WriteHeader(498)
		w.Write([]byte("authorization header format must be Bearer {token}"))
		return ""
	}

	token := strings.TrimPrefix(authHeader, prefix)
	return token
}

func GetClaimsFromTokenString(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		JWTSecret := []byte(strings.TrimSpace(os.Getenv("JWT_SECRET")))
		return JWTSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func GetClaimsFromToken(w http.ResponseWriter, r *http.Request) *Claims {
	claims := &Claims{}

	tokenString := GetJWTFromRequest(w, r)

	// Parse and validate
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// validate algorithm
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return JWTSecret, nil
	})

	if err != nil {
		log.Fatal("Error parsing token:", err)
	}

	if !token.Valid {
		log.Fatal("Invalid token")
	}

	return claims
}
