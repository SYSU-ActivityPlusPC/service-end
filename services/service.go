package services

import (
	"crypto/md5"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/sethvargo/go-password/password"
	"github.com/streadway/amqp"
	"github.com/sysu-activitypluspc/service-end/model"

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
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing my secret
		return hmacSampleSecret, nil
	})
	if err != nil {
		return 0, ""
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		expTime := claims["exp"]
		openID := claims["sub"]

		if (int64)(expTime.(float64)) <= time.Now().Unix() {
			return 1, openID.(string)
		}
		return 2, openID.(string)
	}
	return 0, ""
}

// GenerateJWT Generate jwt with openid(sub), issuance time(iat) and expiration time(exp)
func GenerateJWT(account string) (string, error) {
	// expire in two weeks
	var exp = time.Hour * 24 * 14
	var hmacSampleSecret = []byte(secret)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": account,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(exp).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(hmacSampleSecret)
	return tokenString, err
}

func GetPassword(id string, raw string) string {
	key := strings.Join([]string{raw, raw, id}, "@")
	return GetMd5([]byte(key))
}

// GetMd5 return md5 of given content
func GetMd5(content []byte) string {
	ret := md5.Sum(content)
	return fmt.Sprintf("%x", ret)
}

// CheckIsAdmin check if the given user is admin
func CheckIsAdmin(username string) bool {
	adminAccount := "sysuactivity2018"
	if username == adminAccount {
		return true
	}
	return false
}

// CheckEmail check if the user email exists in the db
func CheckIfEmailExist(email string) bool {
	user := model.GetUserByEmail(email)
	if user == nil {
		return true
	}
	if user.Email == "" {
		return false
	}
	return true
}

// GeneratePassword generate password
func GeneratePassword(length int) string {
	digitNum := 1 + rand.Int()%(length/2)
	res, err := password.Generate(length, digitNum, 0, false, false)
	if err != nil {
		fmt.Println(err)
		return "password"
	}
	return res
}

// WriteMessageQueue write content to the named queue
func WriteMessageQueue(name string, content []byte) bool {
	// Get mq detaild message
	addr := os.Getenv("MQ_ADDRESS")
	if len(addr) == 0 {
		addr = "localhost"
	}
	port := os.Getenv("MQ_PORT")
	if len(port) == 0 {
		port = "5672"
	}
	user := os.Getenv("MQ_USER")
	pass := os.Getenv("MQ_PASSWORD")
	url := fmt.Sprintf("amqp://%s:%s@%s:%s/", user, pass, addr, port)

	// Connect to mq
	conn, err := amqp.Dial(url)
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		name,  // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	// Send
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         content,
		})
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}
