package repository

import (
	"TestTask/model"
	"TestTask/pkg/common"
	"fmt"
	"github.com/jmoiron/sqlx"
	"time"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) CreateUser(user model.User) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (username, password_hash) VALUES ($1, $2) RETURNING id", userTable)

	row := r.db.QueryRow(query, user.Username, user.Password)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *AuthPostgres) GetUser(username, passwordHash string) (model.User, error) {
	var user model.User
	query := fmt.Sprintf("SELECT id FROM %s WHERE username=$1 AND password_hash=$2", userTable)
	err := r.db.Get(&user, query, username, passwordHash)

	return user, err
}

func (r *AuthPostgres) GetUserByUsername(username string) (model.User, error) {
	var user model.User
	query := fmt.Sprintf("SELECT id FROM %s WHERE username=$1", userTable)
	err := r.db.Get(&user, query, username)

	return user, err
}

func (r *AuthPostgres) WriteUserToken(userId int, token string, expireDate time.Time) error {
	query := fmt.Sprintf("INSERT INTO %s (user_id, token, expire_date) VALUES($1, $2, $3)", sessionTable)
	if _, err := r.db.Exec(query, userId, token, expireDate); err != nil {
		return err
	}
	return nil
}

func (r *AuthPostgres) WrongPasswordEnterCnt(userId int) (int, error) {
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE user_id = $1 AND event = $2", authAuditTable)
	var wrongPasswordCnt int
	err := r.db.Get(&wrongPasswordCnt, query, userId, common.WRONGPASSWORDEVENT)
	if err != nil {
		return wrongPasswordCnt, err
	}
	return wrongPasswordCnt, nil
}

func (r *AuthPostgres) WriteEvent(username string, event int) error {
	user, err := r.GetUserByUsername(username)
	if err != nil {
		return err
	}
	query := fmt.Sprintf("INSERT INTO %s (user_id, event_date, event) VALUES($1, $2, $3)", authAuditTable)
	if _, err := r.db.Exec(query, user.Id, time.Now(), event); err != nil {
		return err
	}
	return nil
}
