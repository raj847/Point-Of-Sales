package service

import (
	"context"
	"errors"
	"net"
	"net/mail"
	"strings"
	"time"
	"unicode"
	"vandesar/entity"
	"vandesar/repository"
)

type UserService interface {
	Login(ctx context.Context, user *entity.User) (id int, err error)
	Register(ctx context.Context, user *entity.User) (entity.User, error)

	Delete(ctx context.Context, id int) error
}

type userService struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) UserService {
	return &userService{userRepository}
}

func (s *userService) Login(ctx context.Context, user *entity.User) (id int, err error) {
	//check email and password

	dbUser, err := s.userRepository.GetUserByEmail(ctx, user.Email)
	if err != nil {
		return 0, err
	}

	if dbUser.Email == "" || dbUser.ID == 0 {
		return 0, errors.New("user not found")
	}

	if user.Password != dbUser.Password {
		return 0, errors.New("wrong email or password")
	}

	return dbUser.ID, nil
}

func (s *userService) Register(ctx context.Context, user *entity.User) (entity.User, error) {
	dbUser, err := s.userRepository.GetUserByEmail(ctx, user.Email)
	if err != nil {
		return *user, err
	}

	_, err = mail.ParseAddress(dbUser.Email)

	if err != nil {
		return *user, errors.New("format email invalid")
	}

	domain := strings.Split(dbUser.Email, "@")

	_, err = net.LookupMX(domain[1])

	if err != nil {
		return *user, errors.New("your domain not found")
	}

	if dbUser.Email != "" || dbUser.ID != 0 {
		return *user, errors.New("email already exists")
	}

	user.CreatedAt = time.Now()

	newUser, err := s.userRepository.CreateUser(ctx, *user)
	if err != nil {
		return *user, err
	}

	isMoreThan8 := len(dbUser.Password) > 8

	var isLower, isUpper, isSymbol bool

	for _, char := range dbUser.Password {
		if !isLower && unicode.IsLower(char) {
			isLower = true
		}
		if !isUpper && unicode.IsUpper(char) {
			isUpper = true
		}
		if !isSymbol && (unicode.IsSymbol(char) || unicode.IsPunct(char)) {
			isSymbol = true
		}
	}

	isValid := isMoreThan8 && isLower && isUpper && isSymbol

	if !isValid {
		return *user, errors.New("password is not valid")
	}

	return newUser, nil
}

func (s *userService) Delete(ctx context.Context, id int) error {
	return s.userRepository.DeleteUser(ctx, id)
}
