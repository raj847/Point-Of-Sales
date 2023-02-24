package service

import (
	"context"
	"errors"
	"net"
	"net/mail"
	"strings"
	"unicode"
	"vandesar/entity"
	"vandesar/repository"
	"vandesar/utils"
)

type UserService struct {
	userRepository *repository.UserRepository
}

func NewUserService(userRepository *repository.UserRepository) *UserService {
	return &UserService{
		userRepository: userRepository,
	}
}

var (
	ErrUserNotFound          = errors.New("user not found")
	ErrUserPasswordDontMatch = errors.New("password not match")
	ErrUserAlreadyExists     = errors.New("user already exists")
	ErrEmailInvalid          = errors.New("email invalid")
	ErrPasswordInvalid       = errors.New("password invalid")
)

func (s *UserService) LoginAdmin(ctx context.Context, adminReq entity.AdminLogin) (id entity.Admin, err error) {
	existingAdmin, err := s.userRepository.GetAdminByEmail(ctx, adminReq.Email)
	if err != nil {
		return entity.Admin{}, ErrUserNotFound
	}

	if utils.CheckPassword(adminReq.Password, existingAdmin.Password) != nil {
		return entity.Admin{}, ErrUserPasswordDontMatch
	}

	return existingAdmin, nil
}

func (s *UserService) ChangeAdminPassword(ctx context.Context, changePassReq entity.AdminChangePassword) (id entity.Admin, err error) {
	existingAdmin, err := s.userRepository.GetAdminByID(ctx, changePassReq.AdminID)
	if err != nil {
		return entity.Admin{}, ErrUserNotFound
	}

	if utils.CheckPassword(changePassReq.OldPassword, existingAdmin.Password) != nil {
		return entity.Admin{}, ErrUserPasswordDontMatch
	}

	hashedPassword, _ := utils.HashPassword(changePassReq.NewPassword)

	err = s.userRepository.ChangeAdminPassword(ctx, existingAdmin.ID, hashedPassword)
	if err != nil {
		return entity.Admin{}, err
	}

	existingAdmin.Password = hashedPassword
	return existingAdmin, nil
}

func (s *UserService) LoginCashier(ctx context.Context, cashierReq entity.CashierLogin) (id entity.Cashier, err error) {
	existingCashier, err := s.userRepository.GetCashierByUsername(ctx, cashierReq.Username)
	if err != nil {
		return entity.Cashier{}, ErrUserNotFound
	}

	if utils.CheckPassword(cashierReq.Password, existingCashier.Password) != nil {
		return entity.Cashier{}, ErrUserPasswordDontMatch
	}

	return existingCashier, nil
}

func (s *UserService) RegisterAdmin(ctx context.Context, adminReq entity.AdminRegister) (entity.Admin, error) {
	existingAdmin, err := s.userRepository.GetAdminByEmail(ctx, adminReq.Email)
	if err != nil {
		return entity.Admin{}, err
	}

	if existingAdmin.ID != 0 {
		return entity.Admin{}, ErrUserAlreadyExists
	}

	_, err = mail.ParseAddress(adminReq.Email)
	if err != nil {
		return entity.Admin{}, ErrEmailInvalid
	}

	domain := strings.Split(adminReq.Email, "@")
	_, err = net.LookupMX(domain[1])
	if err != nil {
		return entity.Admin{}, ErrEmailInvalid
	}

	validPassword := validatePassword(adminReq.Password)
	if !validPassword {
		return entity.Admin{}, ErrPasswordInvalid
	}

	admin := entity.Admin{
		ShopName: adminReq.ShopName,
		Email:    adminReq.Email,
		Role:     "admin",
		Password: adminReq.Password,
		PhotoURL: adminReq.PhotoURL,
	}

	hashedPassword, err := utils.HashPassword(admin.Password)
	if err != nil {
		return entity.Admin{}, err
	}

	admin.Password = hashedPassword

	newUser, err := s.userRepository.CreateAdmin(ctx, admin)
	if err != nil {
		return entity.Admin{}, err
	}

	return newUser, nil
}

func (s *UserService) RegisterCashier(ctx context.Context, cashierReq entity.CashierRegister) (entity.Cashier, error) {
	existingCashier, err := s.userRepository.GetCashierByUsername(ctx, cashierReq.Username)
	if err != nil {
		return entity.Cashier{}, err
	}

	if existingCashier.ID != 0 {
		return entity.Cashier{}, ErrUserAlreadyExists
	}

	validPassword := validatePassword(cashierReq.Password)
	if !validPassword {
		return entity.Cashier{}, ErrPasswordInvalid
	}

	cashier := entity.Cashier{
		AdminID:  cashierReq.AdminID,
		Username: cashierReq.Username,
		Role:     "cashier",
		Password: cashierReq.Password,
	}

	hashedPassword, err := utils.HashPassword(cashier.Password)
	if err != nil {
		return entity.Cashier{}, err
	}

	cashier.Password = hashedPassword

	newUser, err := s.userRepository.CreateCashier(ctx, cashier)
	if err != nil {
		return entity.Cashier{}, err
	}

	return newUser, nil
}

func validatePassword(password string) bool {
	var lower, upper, symbol bool
	moreThan := len(password) > 8

	for _, char := range password {
		if unicode.IsLower(char) {
			lower = true
			continue
		}

		if unicode.IsUpper(char) {
			upper = true
			continue
		}

		if unicode.IsSymbol(char) || unicode.IsPunct(char) {
			symbol = true
			continue
		}
	}

	return moreThan && lower && upper && symbol
}

func (s *UserService) CheckTokenAdmin(ctx context.Context, id uint, token entity.CheckToken) (entity.Admin, error) {
	existingAdmin, err := s.userRepository.GetAdminByID(ctx, token.UserId)
	if err != nil {
		return entity.Admin{}, ErrUserNotFound
	}
	return existingAdmin, nil
}

func (s *UserService) CheckTokenCashier(ctx context.Context, id uint, token entity.CheckToken) (entity.Cashier, error) {
	existingCashier, err := s.userRepository.GetCashierByID(ctx, token.UserId)
	if err != nil {
		return entity.Cashier{}, ErrUserNotFound
	}
	return existingCashier, nil
}
func (s *UserService) GetAllCashiers(ctx context.Context, id uint) ([]entity.Cashier, error) {
	return s.userRepository.GetCashierbyAdmin(ctx, id)
}
