package storage

import (
	"errors"
	"geoservice/models"
	"sync"
)

type UserStorage struct {
	users map[string]models.User
	mu    sync.RWMutex
}

func NewUserStorage() *UserStorage {
	return &UserStorage{
		users: make(map[string]models.User),
	}
}

func (s *UserStorage) AddUser(user models.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.users[user.Username]; exists {
		return errors.New("пользователь уже существует")
	}

	s.users[user.Username] = user
	return nil
}

func (s *UserStorage) GetUser(username string) (models.User, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, exists := s.users[username]
	return user, exists
}
