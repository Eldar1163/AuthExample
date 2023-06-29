package service

import (
	"TestTask/model"
	"TestTask/pkg/common"
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

func (s *AuthService) GenerateToken(userId int) (string, error) {
	token := generateTokenByTime(time.Now().Unix())
	err := s.repo.WriteUserToken(userId, token, time.Now().Add(EXPIRETIME))
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *AuthService) CheckUserCredentials(username, password string) (model.User, error) {
	user, err := s.repo.GetUser(username, generatePasswordHash(password))
	if err != nil {
		_ = s.repo.WriteEvent(username, common.WRONGPASSWORDEVENT)
		return user, err
	}
	return user, nil
}

func (s *AuthService) GetBadAuthAttemptsCnt(userId int) (int, error) {
	return s.repo.WrongPasswordEnterCnt(userId)
}

func (s *AuthService) BlockUser(username string) error {
	return s.repo.WriteEvent(username, common.BLOCKEVENT)
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
