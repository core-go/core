package jwt

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"reflect"
	"strings"
	"time"
)

func GenerateToken(payload interface{}, secret string, expiresIn int64) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	// Create the Claims

	//if payload is a map
	if value, ok := payload.(map[string]interface{}); ok {
		claims := token.Claims.(jwt.MapClaims)
		for k, v := range value {
			claims[k] = v
		}
		claims["exp"] = time.Now().Add(time.Millisecond * time.Duration(expiresIn)).Unix()
		claims["iat"] = time.Now().Unix()
		tokenString, err := token.SignedString([]byte(secret))
		return tokenString, err
	}

	s := reflect.ValueOf(payload)
	if s.Kind() == reflect.Ptr {
		s = reflect.Indirect(s)
	}
	typeOfPayload := s.Type()
	claims := token.Claims.(jwt.MapClaims)
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		tag := typeOfPayload.Field(i).Tag
		field := strings.Split(tag.Get("json"), ",")
		if f.IsZero() {
			continue
		}
		claims[field[0]] = f.Interface()
	}
	claims["exp"] = time.Now().Add(time.Millisecond * time.Duration(expiresIn)).Unix()
	claims["iat"] = time.Now().Unix()

	tokenString, err := token.SignedString([]byte(secret))
	return tokenString, err
}

func VerifyToken(tokenString string, secret string) (map[string]interface{}, jwt.StandardClaims, error) {
	//token, err := jwt.ParseWithClaims(tokenString, &commonClaims{}, func(tok *jwt.Token) (interface{}, error) {
	keyLookupFn := func(token *jwt.Token) (interface{}, error) {
		// Check for expected signing method.
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	}
	token, err := jwt.Parse(tokenString, keyLookupFn)
	if err != nil {
		return nil, jwt.StandardClaims{}, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		var exp float64
		var iat float64
		c := jwt.StandardClaims{}

		if x, found := claims["exp"]; found {
			if exp, ok = x.(float64); !ok {
				return nil, jwt.StandardClaims{}, errors.New("exp is invalid (not an integer)")
			}
			c.ExpiresAt = int64(exp)
		} else {
			return nil, jwt.StandardClaims{}, errors.New("'exp' not found")
		}
		if x, found := claims["iat"]; found {
			if iat, ok = x.(float64); !ok {
				return nil, jwt.StandardClaims{}, errors.New("iat is invalid (not an integer)")
			}
			c.IssuedAt = int64(iat)
		} else {
			return nil, jwt.StandardClaims{}, errors.New("'iat' not found")
		}
		delete(claims, "exp")
		delete(claims, "iat")
		result := make(map[string]interface{})
		for k, v := range claims {
			result[k] = v
		}
		return result, c, err
	}
	return nil, jwt.StandardClaims{}, errors.New("invalid token")
}
