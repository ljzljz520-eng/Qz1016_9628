package service

import (
	"labops/model"
	"labops/storage"
	"time"
)

type UserService struct{ Store *storage.Store }

func NewUserService(s *storage.Store) *UserService { return &UserService{Store: s} }
func (s *UserService) Create(u model.User) error {
	if u.CreatedAt.IsZero() {
		u.CreatedAt = time.Now()
	}
	if err := model.ValidateUser(u); err != nil {
		return err
	}
	return s.Store.SaveUser(u)
}
func (s *UserService) Find(id string) (*model.User, error) { return s.Store.GetUser(id) }
