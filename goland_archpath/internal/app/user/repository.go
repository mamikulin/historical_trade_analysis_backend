package user

import (
    "gorm.io/gorm"

	"gorm.io/driver/postgres"
)

type Repository struct {
    DB *gorm.DB
}

func NewRepository(dsn string) (*Repository, error) {
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        return nil, err
    }

    if err := db.AutoMigrate(&User{}); err != nil {
        return nil, err
    }

    return &Repository{DB: db}, nil
}

// Create a new user (registration)
func (r *Repository) CreateUser(user *User) error {
    return r.DB.Create(user).Error
}

// Get user by login
func (r *Repository) GetUserByLogin(login string) (*User, error) {
    var user User
    if err := r.DB.First(&user, "login = ?", login).Error; err != nil {
        return nil, err
    }
    return &user, nil
}

// Update user information
func (r *Repository) UpdateUser(id uint, user *User) error {
    return r.DB.Model(&User{}).Where("id = ?", id).Updates(user).Error
}
