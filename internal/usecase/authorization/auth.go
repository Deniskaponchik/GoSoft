package authorization

import (
	"errors"
	"github.com/deniskaponchik/GoSoft/internal/entity"
	//"crypto/sha1"
	"github.com/golang-jwt/jwt"
	"time"
)

type Jwt struct {
	signingKey string
	//salt string  //для хэширования пароля. пока что не использую
	tokenTTL time.Duration
}

func NewAuthJwt(k string) *Jwt {
	return &Jwt{
		signingKey: k,
		tokenTTL:   15 * time.Minute,
		//salt
	}
}

type tokenClaims struct {
	jwt.StandardClaims
	//UserId int `json:"user_id"`
	UserLogin string `json:"user_login"`
	//UserGivenName string `json:"user_given_name"`
}

func (j *Jwt) GenerateToken(user *entity.User) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(j.tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		//user.Id,
		user.Login,
		//user.GivenName,
	})

	return token.SignedString([]byte(j.signingKey))
}

/*
func (s *AuthService) GenerateToken(username, password string) (string, error) {
	user, err := s.repo.GetUser(username, generatePasswordHash(password))
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user.Id,
	})

	return token.SignedString([]byte(signingKey))
}*/

func (j *Jwt) ParseToken(accessToken string) (string, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(j.signingKey), nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return "", errors.New("token claims are not of type *tokenClaims")
	}

	//return claims.UserId, nil
	return claims.UserLogin, nil
	//return claims.UserGivenName, nil
}

/*
func (s *AuthService) ParseToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(signingKey), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, errors.New("token claims are not of type *tokenClaims")
	}

	return claims.UserId, nil
}*/

/* ORIGINAL
//https://github.com/zhashkevych/todo-app/blob/master/pkg/service/auth.go
const (
	//salt       = "hjqrhjqw124617ajfhajs"  //func generatePasswordHash(password string)
	signingKey = "qrkjk#4#%35FSFJlja#4353KSFjH"
	tokenTTL   = 12 * time.Hour
)

type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateUser(user todo.User) (int, error) {
	user.Password = generatePasswordHash(user.Password)
	return s.repo.CreateUser(user)
}

func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}*/
