package httpserver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

const DefaultExpiredTime = 30 * 24 * time.Hour // 30 days

type TokenPayload struct {
	UserID string `json:"user_id"`
}

func CreateAccessToken(ttl time.Duration, payload TokenPayload, secretJWTKey string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	now := time.Now().UTC()
	claims := token.Claims.(jwt.MapClaims)

	claims["sub"] = payload
	claims["exp"] = now.Add(ttl).Unix()

	tokenString, err := token.SignedString([]byte(secretJWTKey))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateToken(token string, secretJWTKey string) (jwt.MapClaims, error) {
	parseToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, echo.NewHTTPError(http.StatusForbidden, "Unexpected signing method: %v", token.Header["alg"])
		}
		signature := []byte(secretJWTKey)
		return signature, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := parseToken.Claims.(jwt.MapClaims)
	if ok && parseToken.Valid {
		return claims, nil
	} else {
		return nil, fmt.Errorf("invalid token")
	}
}

func DecodeToken(claims jwt.MapClaims) (*TokenPayload, error) {
	// get data from token
	sub, ok := claims["sub"]
	if !ok {
		return nil, fmt.Errorf("missing sub")
	}
	// Convert the map to JSON
	jsonData, err := json.Marshal(sub)
	if err != nil {
		return nil, err
	}
	// Convert the JSON to a struct
	var payload TokenPayload
	if err := json.Unmarshal(jsonData, &payload); err != nil {
		return nil, err
	}
	return &payload, nil
}
