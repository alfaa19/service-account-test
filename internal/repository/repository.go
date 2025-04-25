package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/alfaa19/service-account-test/internal/models"
	logger "github.com/alfaa19/service-account-test/pkg/logrus"
)

type repository struct {
	DB  *sql.DB
	log *logger.CustomLogger
}

type Repository interface {
	GetAccountByNoRekening(ctx context.Context, noRekening string) (*models.Account, error)
	NIKExist(ctx context.Context, nik string) (bool, error)
	PhoneNumberExist(ctx context.Context, phoneNumber string) (bool, error)
	CreateAccount(ctx context.Context, account *models.Account) (*models.Account, error)
	UpdateBalanceWithdraw(ctx context.Context, accountNumber string, amount float64) error
	UpdateBalanceDeposit(ctx context.Context, accountNumber string, amount float64) error
}

func NewRepository(db *sql.DB, log *logger.CustomLogger) Repository {
	return &repository{
		DB:  db,
		log: log,
	}
}

func (r *repository) GetAccountByNoRekening(ctx context.Context, noRekening string) (*models.Account, error) {
	var account models.Account
	query := `SELECT id, account_number, name, nik, phone_number, balance, created_at, updated_at 
			 FROM accounts WHERE account_number = $1`

	r.log.LogOperation(ctx, "GetAccountByNoRekening", "start", map[string]interface{}{
		"type": "repository",
	})

	err := r.DB.QueryRowContext(ctx, query, noRekening).Scan(
		&account.ID,
		&account.AccountNumber,
		&account.Name,
		&account.NIK,
		&account.PhoneNumber,
		&account.Balance,
		&account.CreatedAt,
		&account.UpdatedAt,
	)
	if err != nil {
		r.log.LogOperation(ctx, "GetAccountByNoRekening", "error", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}

	r.log.LogOperation(ctx, "GetAccountByNoRekening", "success", map[string]interface{}{
		"account_id": account.ID,
	})
	return &account, nil
}

func (r *repository) NIKExist(ctx context.Context, nik string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM accounts WHERE nik = $1)`

	r.log.LogOperation(ctx, "GetAccountByNIK", "start", map[string]interface{}{})

	err := r.DB.QueryRowContext(ctx, query, nik).Scan(&exists)
	if err != nil {
		r.log.LogOperation(ctx, "GetAccountByNIK", "error", map[string]interface{}{
			"error": err.Error(),
		})
		return false, err
	}

	r.log.LogOperation(ctx, "GetAccountByNIK", "success", map[string]interface{}{})
	return exists, nil
}

func (r *repository) PhoneNumberExist(ctx context.Context, phoneNumber string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM accounts WHERE phone_number = $1)`

	r.log.LogOperation(ctx, "GetAccountByPhoneNumber", "start", map[string]interface{}{})

	err := r.DB.QueryRowContext(ctx, query, phoneNumber).Scan(&exists)
	if err != nil {
		r.log.LogOperation(ctx, "GetAccountByPhoneNumber", "error", map[string]interface{}{
			"error": err.Error(),
		})
		return false, err
	}

	r.log.LogOperation(ctx, "GetAccountByPhoneNumber", "success", map[string]interface{}{})
	return exists, nil
}

func (r *repository) CreateAccount(ctx context.Context, account *models.Account) (*models.Account, error) {
	query := `INSERT INTO accounts (account_number, name, nik, phone_number, balance, created_at, updated_at) 
			 VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`

	r.log.LogOperation(ctx, "CreateAccount", "start", map[string]interface{}{})

	err := r.DB.QueryRowContext(ctx, query,
		account.AccountNumber,
		account.Name,
		account.NIK,
		account.PhoneNumber,
		account.Balance,
		account.CreatedAt,
		account.UpdatedAt,
	).Scan(&account.ID)
	if err != nil {
		r.log.LogOperation(ctx, "CreateAccount", "error", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}

	r.log.LogOperation(ctx, "CreateAccount", "success", map[string]interface{}{
		"account_id": account.ID,
	})

	return account, nil
}

func (r *repository) UpdateBalanceWithdraw(ctx context.Context, accountNumber string, amount float64) error {
	query := `UPDATE accounts SET balance = balance - $1 WHERE account_number = $2 AND balance >= $1`

	r.log.LogOperation(ctx, "UpdateBalanceWithdraw", "start", map[string]interface{}{
		"account_id": accountNumber,
		"amount":     amount,
	})

	result, err := r.DB.ExecContext(ctx, query, amount, accountNumber)
	if err != nil {
		r.log.LogOperation(ctx, "UpdateBalanceWithdraw", "error", map[string]interface{}{
			"error": err.Error(),
		})
		return err
	}
	rowAffected, _ := result.RowsAffected()
	if rowAffected == 0 {
		r.log.LogOperation(ctx, "UpdateBalanceWithdraw", "error", map[string]interface{}{
			"error": "insufficient balance",
		})
		return errors.New("insufficient balance")

	}

	r.log.LogOperation(ctx, "UpdateBalanceWithdraw", "success", map[string]interface{}{
		"account_id": accountNumber,
	})
	return nil
}

func (r *repository) UpdateBalanceDeposit(ctx context.Context, accountNumber string, amount float64) error {
	query := `UPDATE accounts SET balance = balance + $1 WHERE  account_number = $2`

	r.log.LogOperation(ctx, "UpdateBalanceDeposit", "start", map[string]interface{}{
		"account_id": accountNumber,
		"amount":     amount,
	})

	result, err := r.DB.ExecContext(ctx, query, amount, accountNumber)
	if err != nil {
		r.log.LogOperation(ctx, "UpdateBalanceDeposit", "error", map[string]interface{}{
			"error": err.Error(),
		})
		return err
	}

	rowAffected, _ := result.RowsAffected()
	if rowAffected == 0 {
		r.log.LogOperation(ctx, "UpdateBalanceWithdraw", "error", map[string]interface{}{
			"error": "saldo not updated",
		})
		return errors.New("saldo not updated")

	}
	r.log.LogOperation(ctx, "UpdateBalanceDeposit", "success", map[string]interface{}{
		"account_id": accountNumber,
	})
	return nil
}
