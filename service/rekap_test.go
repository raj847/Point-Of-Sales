package service_test

import (
	"os"
	"testing"
	"vandesar/repository"
	"vandesar/service"
	"vandesar/utils"
)


func TestRekap(t *testing.T) {
	os.Setenv("DATABASE_URL", "postgres://root:secret@localhost:5432/pos")
	utils.ConnectDB()
	db := utils.GetDBConnection()
	userRepo := repository.NewUserRepository(db)
	rekapRepo := repository.NewRekapRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)
	s := service.NewRekapService(rekapRepo, transactionRepo, userRepo)
	s.Rekap()
}
