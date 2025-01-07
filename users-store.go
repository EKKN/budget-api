package main

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"
)

type UsersStorage interface {
	GetDataByLogin(*Users) (*Users, error)
}
type UsersStore struct {
	mysql *MysqlDB
}

func NewUsersStorage(db *MysqlDB) *UsersStore {
	return &UsersStore{
		mysql: db,
	}
}

func (m *UsersStore) GetDataByLogin(users *Users) (*Users, error) {

	hashedPassword := MD5Hash(users.Password)

	AppLog(1)
	query := `SELECT id, userid, password FROM users WHERE userid = ? AND password = ?`
	row := m.mysql.db.QueryRow(query, users.UserID, hashedPassword)

	user := &Users{}
	err := row.Scan(&user.ID, &user.UserID, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by login: %w", err)
	}
	return user, nil
}

func MD5Hash(password string) string {
	hash := md5.Sum([]byte(password))
	return hex.EncodeToString(hash[:])
}
