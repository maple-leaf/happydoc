package services

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"

	"github.com/dgrijalva/jwt-go"
)

type AuthClaims struct {
	Username string
	jwt.StandardClaims
}

func GenerateJWTRS(claims AuthClaims) (token string, publicPEMBase64 string, err error) {
	// https://stackoverflow.com/a/45354453/2326199
	privateKey, err := rsa.GenerateKey(rand.Reader, 1024)
	publicKey := privateKey.PublicKey
	bytes, _ := x509.MarshalPKIXPublicKey(&publicKey)
	publicPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: bytes,
	})
	publicPEMBase64 = base64.StdEncoding.EncodeToString(publicPEM)
	rs256JWT := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token, err = rs256JWT.SignedString(privateKey)

	return
}

func ValidateJWTRS(tokenStr string, publicPEMBase64 string, claims AuthClaims) (token *jwt.Token, err error) {
	publicPEM, err := base64.StdEncoding.DecodeString(publicPEMBase64)
	if err != nil {
		return
	}
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicPEM)
	if err != nil {
		return
	}
	token, err = jwt.ParseWithClaims(tokenStr, &AuthClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return publicKey, nil
	})

	err = validateClaims(token, claims)

	return
}

func validateClaims(token *jwt.Token, claims AuthClaims) (err error) {
	parsedClaims, ok := token.Claims.(*AuthClaims)
	if !ok {
		err = throwClaimInvalidError(*parsedClaims)
	}
	if claims.Username != parsedClaims.Username {
		err = throwClaimInvalidError(*parsedClaims)
	}

	return
}

func throwClaimInvalidError(claims AuthClaims) error {
	return fmt.Errorf("claims invalid: %v", claims)
}
