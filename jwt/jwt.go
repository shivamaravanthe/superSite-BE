package jwt

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/shivamaravanthe/superSite-BE/utils"
)

var secretKey = []byte("supersite_BE_jwt")

func CreateToken(email string) (string, error) {
	token := &jwt.Token{Header: map[string]interface{}{
		"typ": "JWT", "alg": jwt.SigningMethod(jwt.SigningMethodHS256).Alg()},
		Claims: jwt.Claims(jwt.MapClaims{
			"email": email,
			"exp":   time.Now().Add(time.Minute * 500).Unix(),
		}), Method: jwt.SigningMethod(jwt.SigningMethodHS256)}

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func VerifyToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}

func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := VerifyToken(r.Header.Get("token")); err != nil {
			if strings.Contains(err.Error(), "token is expired") {
				utils.CreateErrorResponse(w, 401, "Session Expired", nil)
			} else {
				utils.CreateErrorResponse(w, 401, err.Error(), nil)
			}
			return
		}
		next.ServeHTTP(w, r)
	})
}
