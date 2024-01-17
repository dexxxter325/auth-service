package service

import (
	"CRUD_API"
	"CRUD_API/pkg/repository"
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"time"
)

const (
	salt            = "fds78ufyhj248ifnm2v892eifjm20"
	signingKey      = "dsfsdf6sg7sdg7s7dg67s"
	accessTokenTTL  = 100 * time.Hour
	refreshTokenTTL = 24 * 30 * time.Hour
)

type TokenClaims struct {
	jwt.RegisteredClaims     //стандартный набор данных в полях claims
	UserId               int `json:"user_id"` //будем вшивать в наш токен
}

type ServiceAuthorization struct {
	repo repository.Authorization
}

func NewServiceAuthorization(repo repository.Authorization) *ServiceAuthorization {
	return &ServiceAuthorization{repo: repo}
}

func (s *ServiceAuthorization) CreateUser(user CRUD_API.User) (int, error) {
	user.Password = generateHashPassword(user.Password)
	return s.repo.CreateUser(user)
}
func (s *ServiceAuthorization) GenerateAccessToken(username, password string, hashPassword bool) (string, error) {
	var hashedPassword string
	if hashPassword {
		hashedPassword = generateHashPassword(password)
	} else {
		hashedPassword = password
	}

	user, err := s.repo.GetUser(username, hashedPassword)
	if err != nil {
		log.Println("err in generate AccessToken!!:", err)
		return "", err
	}
	tokenClaims := &TokenClaims{ //наши поля для вшития в токен(claims)
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(accessTokenTTL)), //когда истекает
			IssuedAt:  jwt.NewNumericDate(time.Now()),                     //момент выдачи
		},
		UserId: user.ID, //вшиваем id в наш токен
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims) //.NewWithClaims create new token с methodом подписи HS256 и вшитым id,ср.действия,моментом выдачи
	return token.SignedString([]byte(signingKey))                   //возвращаем наш токен с секретным ключом(для безоп-ости)
}

func (s *ServiceAuthorization) ParseAccessToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Println("err in  jwt.SigningMethodHMAC ")
			return nil, errors.New("invalid signin method((")
		}

		return []byte(signingKey), nil // ключ для проверки подписи
	})
	if err != nil {
		log.Println("err in  jwt.ParseWithClaims ")
		return 0, err
	}
	claims, ok := token.Claims.(*TokenClaims) //проверяем наши поля на подлинность
	if !ok {
		log.Println("err in  token.Claims ")
		return 0, errors.New("token claims aren't of type")
	}
	return claims.UserId, nil //достаем из токена наше зашитое в claims id
}
func (s *ServiceAuthorization) GenerateRefreshToken(userId int) (string, error) {
	// Создаем структуру для хранения данных в payload
	tokenClaims := &TokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(refreshTokenTTL)), // Срок действия в секундах (111 дней)
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		UserId: userId, // Вшиваем id в наш токен
	}

	// Создаем токен с использованием HS256 алгоритма подписи
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims)

	// Подписываем токен с использованием секретного ключа
	return token.SignedString([]byte(signingKey))
}

func (s *ServiceAuthorization) ParseRefreshToken(refreshToken string) (int, error) {
	// Парсим токен с использованием секретного ключа
	token, err := jwt.ParseWithClaims(refreshToken, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Проверяем, что метод подписи совпадает с ожидаемым (HS256)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(signingKey), nil
	})

	if err != nil {
		return 0, fmt.Errorf("failed to parse claims:%s", err)
	}

	// Извлекаем данные из токена
	claims, ok := token.Claims.(*TokenClaims)
	if !ok {
		return 0, fmt.Errorf("failed to compare claims:%s", err)
	}

	return claims.UserId, nil
}

// GenerateNewTokenPair обновляет пару access и refresh токенов
func (s *ServiceAuthorization) GenerateNewTokenPair(refreshToken string) (string, string, error) {
	// Парсим текущий refresh токен
	userID, err := s.ParseRefreshToken(refreshToken)
	if err != nil {
		log.Println("Error in parsing refresh token:", err)
		return "", "", err
	}
	user, err := s.repo.GetUserById(userID)
	if err != nil {
		log.Println("Error in GetUserById:", err)
		return "", "", err
	}

	// Создаем новую пару токенов
	newAccessToken, err := s.GenerateAccessToken(user.Username, user.Password, false)
	if err != nil {
		log.Println("Error in GenerateAccessToken:", err)
		return "", "", err
	}

	newRefreshToken, err := s.GenerateRefreshToken(userID)
	if err != nil {
		log.Println("Error in GenerateRefreshToken:", err)
		return "", "", err
	}

	return newAccessToken, newRefreshToken, nil
}

// GenerateHashPassword /*хэшируем наш пароль(для усложнения взлома инфы)*/
func generateHashPassword(password string) string {
	hash := sha1.New()                               // Создаем новый хеш(кодировка) SHA1
	hash.Write([]byte(password))                     // Записываем в него байтовое представление пароля
	return fmt.Sprintf("%x", hash.Sum([]byte(salt))) // добавляем соль в хэш-пароль и возвращаем его
}

/*
eXAmple of JWT TOKEN:
Header(ALGORITHM & TOKEN TYPE)Инфа о самом токене:
{
  "alg": "HS256",
  "typ": "JWT"
}

ClAIMS/PAYLOAD(DATA)ID токена,роль юзера и тд:
{
  "sub": "1234567890",
  "name": "John Doe",
  "iat": 1516239022
}

VERIFY SIGNATURE(подпись):
HMACSHA256(
  base64UrlEncode(header) + "." +
  base64UrlEncode(payload),
)
Из этих данных собирается токен:
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.-Header
eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.-Payload
SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c-Подпись
*/
