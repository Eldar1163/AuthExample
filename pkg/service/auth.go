package service

import (
	"TestTask/model"
	"TestTask/pkg/repository"
	"crypto/sha1"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"time"
)

const (
	SALT       = "kjsanfjksdfsdf12354dsvsujdkhnv"
	EXPIRETIME = 24 * time.Hour
)

type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateUser(user model.User) (int, error) {
	user.Password = generatePasswordHash(user.Password)
	return s.repo.CreateUser(user)
}

func (s *AuthService) GenerateToken(username, password string) (string, error) {
	user, err := s.repo.GetUser(username, generatePasswordHash(password))
	if err != nil {
		return "", err
	}

	token := generateTokenByTime(time.Now().Unix())
	err = s.repo.WriteUserToken(user.Id, token, time.Now().Add(EXPIRETIME))
	if err != nil {
		return "", err
	}
	return token, nil
}

func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(SALT)))
}

func generateTokenByTime(unixTime int64) string {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, uint64(unixTime))

	hash := sha1.New()
	hash.Write(buf)
	bs := hash.Sum(nil)

	return hex.EncodeToString(bs)
}
