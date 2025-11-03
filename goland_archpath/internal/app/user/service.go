package user

import (

    "golang.org/x/crypto/bcrypt"
    "fmt"
)

type Service struct {
    repo *Repository
}

func NewService(repo *Repository) *Service {
    return &Service{repo: repo}
}

// HashPassword hashes the password before storing it in the database
func (s *Service) HashPassword(password string) (string, error) {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return "", err
    }
    return string(hashedPassword), nil
}

// ComparePasswords checks if the provided password matches the stored hash
func (s *Service) ComparePasswords(hashedPassword, password string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
    return err == nil
}

// Register a new user
func (s *Service) RegisterUser(user *User, password string) error {
    // Hash the password
    hashedPassword, err := s.HashPassword(password)
    if err != nil {
        return fmt.Errorf("could not hash password: %w", err)
    }
    user.PasswordHash = hashedPassword

    // Save the user to the database
    return s.repo.CreateUser(user)
}

// Authenticate user by login and password
func (s *Service) AuthenticateUser(login, password string) (*User, error) {
    user, err := s.repo.GetUserByLogin(login)
    if err != nil {
        return nil, fmt.Errorf("user not found: %w", err)
    }

    // Check if the password matches
    if !s.ComparePasswords(user.PasswordHash, password) {
        return nil, fmt.Errorf("incorrect password")
    }

    return user, nil
}

// Update user information
func (s *Service) UpdateUser(id uint, user *User) error {
    return s.repo.UpdateUser(id, user)
}
