package service

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/alfaa19/service-account-test/internal/models"
	"github.com/alfaa19/service-account-test/internal/models/dto"
	"github.com/alfaa19/service-account-test/internal/repository"
	logger "github.com/alfaa19/service-account-test/pkg/logrus"
)

type service struct {
	repo repository.Repository
	log  *logger.CustomLogger
}

type Service interface {
	GetAccountByNoRekening(ctx context.Context, noRekening string) (*models.Account, error)
	CreateAccount(ctx context.Context, reqAccount *dto.AccountRegistration) (*models.Account, error)
	UpdateBalanceWithdraw(ctx context.Context, accountNumber string, amount float64) error
	UpdateBalanceDeposit(ctx context.Context, accountNumber string, amount float64) error
}

func NewService(repo repository.Repository, log *logger.CustomLogger) Service {
	return &service{
		repo: repo,
		log:  log,
	}
}

func (s *service) GetAccountByNoRekening(ctx context.Context, noRekening string) (*models.Account, error) {
	account, err := s.repo.GetAccountByNoRekening(ctx, noRekening)
	if err != nil {
		s.log.LogOperation(ctx, "GetAccountByNoRekening", "error", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}

	s.log.LogOperation(ctx, "GetAccountByNoRekening", "success", map[string]interface{}{
		"type":       "service",
		"account_id": account.ID,
	})
	return account, nil
}

func (s *service) CreateAccount(ctx context.Context, reqAccount *dto.AccountRegistration) (*models.Account, error) {
	// Check if NIK already exists
	existingAccountByNIK, err := s.repo.NIKExist(ctx, reqAccount.NIK)
	if err != nil {
		s.log.LogOperation(ctx, "CreateAccount", "error", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}
	if existingAccountByNIK {

		s.log.LogOperation(ctx, "CreateAccount", "error", map[string]interface{}{
			"error": errors.New("NIK already exists"),
		})
		return nil, errors.New("NIK already exists")
	}

	// Check if phone number already exists
	existingAccountByPhone, err := s.repo.PhoneNumberExist(ctx, reqAccount.NoHP)
	if err != nil {
		s.log.LogOperation(ctx, "CreateAccount", "error", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}
	if existingAccountByPhone {
		s.log.LogOperation(ctx, "CreateAccount", "error", map[string]interface{}{
			"error": errors.New("phone Number already exists"),
		})
		return nil, errors.New("phone Number already exists")
	}

	account := &models.Account{
		AccountNumber: generateAccountNumber(),
		Name:          reqAccount.Nama,
		NIK:           reqAccount.NIK,
		PhoneNumber:   reqAccount.NoHP,
		Balance:       0.0,
	}

	// Create the account
	createdAccount, err := s.repo.CreateAccount(ctx, account)
	if err != nil {
		s.log.LogOperation(ctx, "CreateAccount", "error", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}

	s.log.LogOperation(ctx, "CreateAccount", "success", map[string]interface{}{
		"type": "service",
	})
	return createdAccount, nil
}
func (s *service) UpdateBalanceWithdraw(ctx context.Context, accountNumber string, amount float64) error {
	err := s.repo.UpdateBalanceWithdraw(ctx, accountNumber, amount)
	if err != nil {
		s.log.LogOperation(ctx, "UpdateBalanceWithdraw", "error", map[string]interface{}{
			"error": err.Error(),
		})
		return err
	}

	s.log.LogOperation(ctx, "UpdateBalanceWithdraw", "success", map[string]interface{}{
		"type":       "service",
		"account_id": accountNumber,
	})
	return nil
}
func (s *service) UpdateBalanceDeposit(ctx context.Context, accountNumber string, amount float64) error {
	err := s.repo.UpdateBalanceDeposit(ctx, accountNumber, amount)
	if err != nil {
		s.log.LogOperation(ctx, "UpdateBalanceDeposit", "error", map[string]interface{}{
			"error": err.Error(),
		})
		return err
	}

	s.log.LogOperation(ctx, "UpdateBalanceDeposit", "success", map[string]interface{}{
		"type":       "service",
		"account_id": accountNumber,
	})
	return nil
}

func generateAccountNumber() string {
	rand.Seed(time.Now().UnixNano())
	accountNumber := ""
	for i := 0; i < 10; i++ {
		accountNumber += fmt.Sprintf("%d", rand.Intn(10))
	}
	return accountNumber
}
