package service

import (
	"context"
	"errors"
	"net"
	"net/mail"
	"strings"
	"unicode"
	"vandesar/entity"
	"vandesar/utils"
)

type AdminRepository interface {
	GetAdminByID(ctx context.Context, id uint) (entity.Admin, error)
	GetAdminByEmail(ctx context.Context, email string) (entity.Admin, error)
	CreateAdmin(ctx context.Context, user entity.Admin) (entity.Admin, error)
	UpdateAdmin(ctx context.Context, user entity.Admin) (entity.Admin, error)
	DeleteAdmin(ctx context.Context, id uint) error
}

type CashierRepository interface {
	GetCashierByID(ctx context.Context, id uint) (entity.Cashier, error)
	GetCashierByUsername(ctx context.Context, username string) (entity.Cashier, error)
	CreateCashier(ctx context.Context, user entity.Cashier) (entity.Cashier, error)
	UpdateCashier(ctx context.Context, user entity.Cashier) (entity.Cashier, error)
	DeleteCashier(ctx context.Context, id uint) error
}

type UserRepository interface {
	AdminRepository
	CashierRepository
}

type UserService struct {
	adminRepository   AdminRepository
	cashierRepository CashierRepository
}

func NewUserService(userRepository UserRepository) *UserService {
	return &UserService{
		adminRepository:   userRepository,
		cashierRepository: userRepository,
	}
}

func (s *UserService) LoginAdmin(ctx context.Context, adminReq entity.AdminLogin) (id entity.Admin, err error) {
	existingAdmin, err := s.adminRepository.GetAdminByEmail(ctx, adminReq.Email)
	if err != nil {
		return entity.Admin{}, errors.New("user not found")
	}

	if existingAdmin.Email == "" || existingAdmin.ID == 0 {
		return entity.Admin{}, errors.New("user not found")
	}
	if existingAdmin.Email != adminReq.Email {
		return entity.Admin{}, errors.New("email not found")
	}

	if utils.CheckPassword(adminReq.Password, existingAdmin.Password) != nil {
		return entity.Admin{}, errors.New("password not match")
	}

	return existingAdmin, nil
}

func (s *UserService) LoginCashier(ctx context.Context, cashierReq entity.CashierLogin) (id entity.Cashier, err error) {
	existingCashier, err := s.cashierRepository.GetCashierByUsername(ctx, cashierReq.Username)
	if err != nil {
		return entity.Cashier{}, err
	}

	if existingCashier.Username == "" || existingCashier.ID == 0 {
		return entity.Cashier{}, errors.New("user not found")
	}
	if existingCashier.Username != cashierReq.Username {
		return entity.Cashier{}, errors.New("username not found")
	}

	if utils.CheckPassword(cashierReq.Password, existingCashier.Password) != nil {
		return entity.Cashier{}, errors.New("password not match")
	}

	return existingCashier, nil
}

func (s *UserService) RegisterAdmin(ctx context.Context, adminReq entity.AdminRegister) (entity.Admin, error) {
	existingAdmin, err := s.adminRepository.GetAdminByEmail(ctx, adminReq.Email)
	if err != nil {
		return entity.Admin{}, err
	}

	if existingAdmin.Email != "" || existingAdmin.ID != 0 {
		return entity.Admin{}, errors.New("email already exists")
	}

	_, err = mail.ParseAddress(adminReq.Email)
	if err != nil {
		return entity.Admin{}, errors.New("format email invalid")
	}

	domain := strings.Split(adminReq.Email, "@")
	_, err = net.LookupMX(domain[1])
	if err != nil {
		return entity.Admin{}, errors.New("your domain not found")
	}

	var lower, upper, symbol bool
	moreThan := len(adminReq.Password) > 8

	for _, char := range adminReq.Password {
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

	valid := moreThan && lower && upper && symbol
	if !valid {
		return entity.Admin{}, errors.New("password is not valid")
	}

	admin := entity.Admin{
		ShopName: adminReq.ShopName,
		Email:    adminReq.Email,
		Role:     "admin",
		Password: adminReq.Password,
	}

	hashedPassword, err := utils.HashPassword(admin.Password)
	if err != nil {
		return entity.Admin{}, err
	}
	admin.Password = hashedPassword

	newUser, err := s.adminRepository.CreateAdmin(ctx, admin)
	if err != nil {
		return entity.Admin{}, err
	}

	return newUser, nil
}

func (s *UserService) RegisterCashier(ctx context.Context, cashierReq entity.CashierRegister) (entity.Cashier, error) {
	existingAdmin, err := s.cashierRepository.GetCashierByUsername(ctx, cashierReq.Username)
	if err != nil {
		return entity.Cashier{}, err
	}

	if existingAdmin.Username != "" || existingAdmin.ID != 0 {
		return entity.Cashier{}, errors.New("email already exists")
	}

	var lower, upper, symbol bool
	moreThan := len(cashierReq.Password) > 8

	for _, char := range cashierReq.Password {
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

	valid := moreThan && lower && upper && symbol
	if !valid {
		return entity.Cashier{}, errors.New("password is not valid")
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

	newUser, err := s.cashierRepository.CreateCashier(ctx, cashier)
	if err != nil {
		return entity.Cashier{}, err
	}

	return newUser, nil
}
