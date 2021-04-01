package jwt

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
	"io/ioutil"
	"time"
)

type TokenClaims struct {
	ExpiresAt int64  `json:"exp,omitempty"`
	Id        string `json:"jti,omitempty"`
	IssuedAt  int64  `json:"iat,omitempty"`
	NotBefore int64  `json:"nbf,omitempty"`
	Subject   int64  `json:"sub,omitempty"`
	Issuer    string  `json:"iss,omitempty"`
	jwt.StandardClaims
}

func Decode(t string) (*TokenClaims, error) {
	if viper.GetString("rsa_key.public") == "" {
		return nil, fmt.Errorf("please specify your public key path")
	}
	publicKey, errKey := ioutil.ReadFile(viper.GetString("rsa_key.public"))
	if errKey != nil {
		return nil, fmt.Errorf("error while reading public key file : %v", errKey)
	}

	key, errParse := jwt.ParseRSAPublicKeyFromPEM(publicKey)
	if errParse != nil {
		return nil, fmt.Errorf("failed to parse public key : %v", errParse)
	}

	token, err := jwt.ParseWithClaims(t, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod("RS256") != token.Method {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return key, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token : %v", err)
	}

	if claims, ok := token.Claims.(*TokenClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}

func Encode(id int64, iss string) (string, error) {
	if viper.GetString("rsa_key.private") == "" {
		return "", fmt.Errorf("please specify your public key path")
	}
	hashQuery := md5.New()
	hashQuery.Write([]byte(fmt.Sprintf("secret123:%v", time.Now().Add(time.Hour*24*365).Unix())))

	jti := hex.EncodeToString(hashQuery.Sum(nil))

	sign := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iat": time.Now().Unix(),
		"nbf": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour * 24 * 365).Unix(),
		"jti": jti,
		"sub": id,
		"iss": iss,
	})

	privateKey, readErr := ioutil.ReadFile(viper.GetString("rsa_key.private"))
	if readErr != nil {
		return "", fmt.Errorf("therewas and error while trying to read private key file : %v", readErr)
	}

	key, parseErr := jwt.ParseRSAPrivateKeyFromPEM(privateKey)
	if parseErr != nil {
		return "", fmt.Errorf("therewas and error while trying to parse private key : %v", parseErr)
	}

	if token, err := sign.SignedString(key); err != nil {
		return "", fmt.Errorf("therewas and error while trying to create token : %v", parseErr)
	} else {
		return token, nil
	}
}
