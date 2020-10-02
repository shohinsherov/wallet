package wallet

import (
	"errors"

	"github.com/google/uuid"
	"github.com/shohinsherov/wallet/pkg/types"
)

// ErrPhoneRegistered ...
var ErrPhoneRegistered = errors.New("phone already registered")

// ErrAmountMustBePositive ...
var ErrAmountMustBePositive = errors.New("amount must be grated than zero")

// ErrAccountNotFound ...
var ErrAccountNotFound = errors.New("Account not found")

// ErrNotEnoughBalance ...
var ErrNotEnoughBalance = errors.New("Not enough balance")

// ErrPaymentNotFound ....
var ErrPaymentNotFound = errors.New("Payment not found")

// Service ....
type Service struct {
	nextAccountID int64
	accounts      []*types.Account
	payments      []*types.Payment
}

// Error ....
type Error string

func (e Error) Error() string {
	return string(e)
}

// RegisterAccount ...
func (s *Service) RegisterAccount(phone types.Phone) (*types.Account, error) {
	for _, account := range s.accounts {
		if account.Phone == phone {
			return nil, ErrPhoneRegistered
		}
	}
	s.nextAccountID++
	account := &types.Account{
		ID:      s.nextAccountID,
		Phone:   phone,
		Balance: 0,
	}
	s.accounts = append(s.accounts, account)

	return account, nil
}

// Deposit ...
func (s *Service) Deposit(accountID int64, amount types.Money) error {
	if amount <= 0 {
		return ErrAmountMustBePositive
	}

	var account *types.Account
	for _, acc := range s.accounts {
		if acc.ID == accountID {
			account = acc
			break
		}
	}

	if account == nil {
		return ErrAccountNotFound
	}

	account.Balance += amount
	return nil

}

// Pay ...
func (s *Service) Pay(accountID int64, amount types.Money, category types.PaymentCategory) (*types.Payment, error) {
	if amount <= 0 {
		return nil, ErrAmountMustBePositive
	}

	var account *types.Account
	for _, acc := range s.accounts {
		if acc.ID == accountID {
			account = acc
			break
		}
	}

	if account == nil {
		return nil, ErrAccountNotFound
	}

	if account.Balance < amount {
		return nil, ErrNotEnoughBalance
	}

	account.Balance -= amount

	paymentID := uuid.New().String()
	payment := &types.Payment{
		ID:        paymentID,
		AccountID: accountID,
		Amount:    amount,
		Category:  category,
		Status:    types.PaymentStatusInProgress,
	}

	s.payments = append(s.payments, payment)
	return payment, nil
}

// FindAccountByID ....
func (s *Service) FindAccountByID(accountID int64) (*types.Account, error) {
	for _, account := range s.accounts {
		if account.ID == accountID {
			return account, nil
		}
		break
	}
	return nil, ErrAccountNotFound
}

// FindPaymentByID ....
func (s *Service) FindPaymentByID(paymentID string) (*types.Payment, error) {
	for _, payment := range s.payments {
		if payment.ID == paymentID {
			return payment, nil
		}
		break
	}
	return nil, ErrPaymentNotFound
}

// Reject ...
func (s *Service) Reject(paymentID string) (*types.Payment, error) {
	findPayment, err := s.FindPaymentByID(paymentID)
	if err != nil {
		return nil, ErrPaymentNotFound
	}
	findPayment.Status = types.PaymentStatusFail
	findAccount, err := s.FindAccountByID(findPayment.AccountID)
	if err != nil {
		return nil, ErrAccountNotFound
	}
	findAccount.Balance += findPayment.Amount

	return findPayment, nil

}
