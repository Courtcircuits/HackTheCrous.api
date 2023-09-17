package util

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
)

func Hash(s string) string {
	hash := sha256.Sum256([]byte(s))
	return hex.EncodeToString(hash[:])
}

// returns a hash of a random string + parameter s and the second parameter
func HashAndSalted(s string) (string, string) {
	hash_round, err := strconv.Atoi(Get("HASH_ROUND"))
	if err != nil {
		panic("HASH_ROUND is not a number")
	}
	salt := GenNonce(hash_round)
	return Hash(salt + s), salt
}

func CompareHash(hash string, clear string, salt string) bool {
	return hash == Hash(salt+clear)
}

func GenNonce(length int) string {
	chars := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZa"
	result := ""
	for i := 0; i < length; i++ {
		result += string(chars[rand.Intn(len(chars))])
	}
	return result
}

func GenKeysJWT() []byte {

	_, private_key, err := ed25519.GenerateKey(nil)

	if err != nil {
		panic(err)
	}

	return []byte(private_key)
}

func GenActivationCode() string {
	chars := "0123456789"
	result := ""
	for i := 0; i < 6; i++ {
		result += string(chars[rand.Intn(len(chars))])
	}
	return result
}

func GenJWT(expiration time.Time, payload map[string]any) string {
	secret := []byte(Get("JWT_SECRET"))

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = expiration
	for key, val := range payload {
		claims[key] = val
	}

	tokenString, err := token.SignedString(secret)

	if err != nil {
		log.Printf("Error while signing JWT : %q", err)
		return ""
	}

	return tokenString
}

var ErrInvalidToken = errors.New("jwt token is invalid for some reason")

func VerifyJWT(raw_token string) (map[string]any, error) {
	token, err := jwt.Parse(raw_token, func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return map[string]any{
				"error": "Unauthorized",
			}, errors.New("Unauthorized")
		}
		return []byte(Get("JWT_SECRET")), nil
	})

	if err != nil {
		return map[string]any{}, err
	}

	var toReturn map[string]any = make(map[string]any)

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		for key, val := range claims {
			toReturn[key] = val
		}
		return toReturn, nil
	}
	return map[string]any{}, ErrInvalidToken
}

func GenRefreshToken() string {
	expiration := time.Now().Add(time.Hour * 24 * 7 * 3) // in a week
	refresh_payload := map[string]any{
		"nonce": GenNonce(30),
	}
	return GenJWT(expiration, refresh_payload)
}
