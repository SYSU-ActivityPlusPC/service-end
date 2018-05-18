package controller

import (
	"io"
	"crypto/md5"
	"fmt"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

const secret = "sysu_activity_2018_activity_sysu"

// GetPoster judge if the poster and returns accurate one with given type
func GetPoster(raw string, actType int) string {
	if len(raw) == 0 {
		switch actType {
		// physics
		case 0:
			return "b6f487c6d08921463a6ebc0612d9fe1f.gif"
		// volunteer
		case 1:
			return "ccc55f553829fabb7c15227d79450dae.gif"
		// match
		case 2:
			return "2bee829b10b0a84002cf5cb5c4a3c8f3.gif"
		// show
		case 3:
			return "68dac067d05a98995a353ad8265b1f09.png"
		// speech
		case 4:
			return "a90dc26fbd5299e4053a3bbc39b5afc8.gif"
		// outdoor
		case 5:
			return "e8ae3078dfa14c62ff1e71104ec0b11f.png"
		// relax
		case 6:
			return "b2b71f5f39d3a4389d34ce1b248e9fee.png"
		}
	}
	return raw
}

// CheckToken Check token and return token status code with openId
// status code: 0 -> check error; 1 -> timeout; 2 -> ok
func CheckToken(tokenString string) (int, string) {
	var hmacSampleSecret = []byte(secret)
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing my secret
		return hmacSampleSecret, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		expTime := claims["exp"]
		openID := claims["sub"]

		if (int64)(expTime.(float64)) <= time.Now().Unix() {
			return 1, openID.(string)
		}
		return 2, openID.(string)
	} else {
		return 0, ""
	}
}

// GenerateJWT Generate jwt with openid(sub), issuance time(iat) and expiration time(exp)
func GenerateJWT(openId string) (string, error) {
	// expire in two weeks
	var exp = time.Hour * 24 * 14
	var hmacSampleSecret = []byte(secret)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": openId,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(exp).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(hmacSampleSecret)
	return tokenString, err
}

func getPassword(id string, raw string) string {
	key := strings.Join([]string{raw, raw, id}, "@")
	return GetMd5([]byte(key))
}

// GetMd5 return md5 of given content
func GetMd5(content []byte) string {
	ret := md5.Sum(content)
	return string(ret[:])
}

// GetFileMd5 get md5 of file
func GetFileMd5(f io.Reader) string{
	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		fmt.Println(err)
		return ""
	}
	ret := h.Sum(nil)
	return string(ret[:])
}
