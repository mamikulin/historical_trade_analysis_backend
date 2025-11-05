package user

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) RegisterUser(user *User, password string) error {
	if user.Login == "" || password == "" {
		return fmt.Errorf("login and password are required")
	}

	// Проверяем, существует ли пользователь
	existingUser, _ := s.repo.GetByLogin(user.Login)
	if existingUser != nil {
		return fmt.Errorf("user with login %s already exists", user.Login)
	}

	// Хешируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	user.PasswordHash = string(hashedPassword)
	
	// Устанавливаем роль по умолчанию, если не указана
	if user.Role == "" {
		user.Role = "user"
	}

	return s.repo.Create(user)
}

func (s *Service) AuthenticateUser(login, password string) (*User, error) {
	user, err := s.repo.GetByLogin(login)
	if err != nil {
		return nil, fmt.Errorf("invalid login or password")
	}

	// Проверяем пароль
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("invalid login or password")
	}

	return user, nil
}

func (s *Service) GetUserByID(id uint) (*User, error) {
	return s.repo.GetByID(id)
}

func (s *Service) UpdateUser(id uint, user *User) error {
	return s.repo.Update(id, user)
}