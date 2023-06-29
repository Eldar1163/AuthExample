package repository

import (
	"TestTask/model"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"time"
)

type ApiPostgres struct {
	db *sqlx.DB
}

func NewApiPostgres(db *sqlx.DB) *ApiPostgres {
	return &ApiPostgres{db: db}
}

type sessionStruct struct {
	UserId     int       `db:"user_id"`
	Token      string    `db:"token"`
	ExpireDate time.Time `db:"expire_date"`
}

func (a *ApiPostgres) CheckToken(token string) (int, error) {
	var s sessionStruct
	query := fmt.Sprintf("SELECT user_id, token, expire_date FROM %s WHERE token = $1", sessionTable)
	err := a.db.Get(&s, query, token)
	if err != nil {
		return 0, err
	}
	if s.ExpireDate.Before(time.Now()) {
		return 0, errors.New("token is expired")
	}
	return s.UserId, nil
}

func (a *ApiPostgres) ClearAudit(userId int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE user_id = $1", authAuditTable)
	_, err := a.db.Exec(query, userId)
	if err != nil {
		return err
	}
	return nil
}

func (a *ApiPostgres) GetUserEvents(token string) ([]model.Audit, error) {
	var events []model.Audit
	query := fmt.Sprintf("SELECT event_date, CASE WHEN event = 0 THEN 'WRONG PASSWORD' WHEN event = 1 THEN 'BLOCKED' END AS event FROM %s WHERE user_id = (SELECT user_id FROM %s WHERE token = $1)", authAuditTable, sessionTable)
	err := a.db.Select(&events, query, token)
	if err != nil {
		return nil, err
	}
	return events, nil
}
