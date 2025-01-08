package main

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"
)

type UsersStorage interface {
	GetByLogin(*Users) (*Users, error)
}

type UsersStore struct {
	db *sql.DB
}

func NewUsersStorage(db *sql.DB) *UsersStore {
	return &UsersStore{
		db: db,
	}
}

func (s *UsersStore) GetByLogin(user *Users) (*Users, error) {
	hashedPassword := MD5Hash(user.Password)

	query := `SELECT id, userid, password FROM users WHERE userid = ? AND password = ?`
	row := s.db.QueryRow(query, user.UserID, hashedPassword)

	retrievedUser := &Users{}
	err := row.Scan(&retrievedUser.ID, &retrievedUser.UserID, &retrievedUser.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by login: %w", err)
	}
	return retrievedUser, nil
}

func MD5Hash(password string) string {
	hash := md5.Sum([]byte(password))
	return hex.EncodeToString(hash[:])
}
