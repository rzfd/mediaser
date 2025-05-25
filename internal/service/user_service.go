package service

import (
	"errors"
	"strings"

	"github.com/rzfd/mediashar/internal/models"
	"github.com/rzfd/mediashar/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Create(user *models.User) error
	GetByID(id uint) (*models.User, error)
	GetByUsername(username string) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	Update(user *models.User) error
	Delete(id uint) error
	List(page, pageSize int) ([]*models.User, error)
	GetStreamers(page, pageSize int) ([]*models.User, error)
	Authenticate(email, password string) (*models.User, error)
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) Create(user *models.User) error {
	// Validate inputs
	if user.Username == "" || user.Email == "" || user.Password == "" {
		return errors.New("username, email and password are required")
	}

	// Check if email is valid
	if !strings.Contains(user.Email, "@") {
		return errors.New("invalid email format")
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	return s.userRepo.Create(user)
}

func (s *userService) GetByID(id uint) (*models.User, error) {
	return s.userRepo.GetByID(id)
}

func (s *userService) GetByUsername(username string) (*models.User, error) {
	return s.userRepo.GetByUsername(username)
}

func (s *userService) GetByEmail(email string) (*models.User, error) {
	return s.userRepo.GetByEmail(email)
}

func (s *userService) Update(user *models.User) error {
	// If password is being updated, hash it
	if user.Password != "" {
		existingUser, err := s.userRepo.GetByID(user.ID)
		if err != nil {
			return err
		}

		// Only hash if password has changed
		if user.Password != existingUser.Password {
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
			if err != nil {
				return err
			}
			user.Password = string(hashedPassword)
		}
	}

	return s.userRepo.Update(user)
}

func (s *userService) Delete(id uint) error {
	return s.userRepo.Delete(id)
}

func (s *userService) List(page, pageSize int) ([]*models.User, error) {
	offset := (page - 1) * pageSize
	return s.userRepo.List(offset, pageSize)
}

func (s *userService) GetStreamers(page, pageSize int) ([]*models.User, error) {
	offset := (page - 1) * pageSize
	return s.userRepo.GetStreamers(offset, pageSize)
}

func (s *userService) Authenticate(email, password string) (*models.User, error) {
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
} 