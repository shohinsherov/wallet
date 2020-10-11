package wallet

import (
	"errors"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

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

var ErrFavoriteNotFound = errors.New("Favorite not found")

// Service ....
type Service struct {
	nextAccountID int64
	accounts      []*types.Account
	payments      []*types.Payment
	favorites     []*types.Favorite
}

//// ---------------------------------------------------

//////////////-------------------------------------------
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

	}
	return nil, ErrPaymentNotFound
}

// Reject ...
func (s *Service) Reject(paymentID string) error {
	findPayment, err := s.FindPaymentByID(paymentID)
	if err != nil {
		return err
	}
	findPayment.Status = types.PaymentStatusFail
	findAccount, err := s.FindAccountByID(findPayment.AccountID)
	if err != nil {
		return err
	}
	findAccount.Balance += findPayment.Amount

	return nil

}

func (s *Service) Repeat(paymentID string) (*types.Payment, error) {
	findPayment, err := s.FindPaymentByID(paymentID)
	if err != nil {
		return nil, err
	}
	newPayment, err := s.Pay(findPayment.AccountID, findPayment.Amount, findPayment.Category)
	if err != nil {
		return nil, err
	}
	return newPayment, nil
}

func (s *Service) FavoritePayment(paymentID string, name string) (*types.Favorite, error) {
	findPayment, err := s.FindPaymentByID(paymentID)
	if err != nil {
		return nil, err
	}
	favID := uuid.New().String()
	favorite := &types.Favorite{
		ID:        favID,
		AccountID: findPayment.AccountID,
		Name:      name,
		Amount:    findPayment.Amount,
		Category:  findPayment.Category,
	}
	s.favorites = append(s.favorites, favorite)

	return favorite, nil
}

func (s *Service) FindFavoriteByID(favoriteID string) (*types.Favorite, error) {
	for _, favorite := range s.favorites {
		if favoriteID == favorite.ID {
			return favorite, nil
		}
	}
	return nil, ErrFavoriteNotFound
}

func (s *Service) PayFromFavorite(favoriteID string) (*types.Payment, error) {
	findFavorite, err := s.FindFavoriteByID(favoriteID)
	if err != nil {
		return nil, err
	}

	payFavorite, err := s.Pay(findFavorite.AccountID, findFavorite.Amount, findFavorite.Category)
	if err != nil {
		return nil, err
	}

	return payFavorite, nil
}

func (s *Service) ExportToFile(path string) error {
	file, err := os.Create(path)
	if err != nil {
		log.Print(err)
		return err
	}

	for _, acc := range s.accounts {
		_, err = file.Write([]byte(strconv.FormatInt(int64(acc.ID), 10)))
		if err != nil {
			log.Print(err)
			return err
		}

		defer func() {
			if cerr := file.Close(); cerr != nil {
				log.Print(err)
				//return cerr
			}
		}()
		_, err = file.Write([]byte(string(";")))
		if err != nil {
			log.Print(err)
			return err
		}
		_, err = file.Write([]byte(string(acc.Phone)))
		if err != nil {
			log.Print(err)
			return err
		}
		_, err = file.Write([]byte(string(";")))
		if err != nil {
			log.Print(err)
			return err
		}
		_, err = file.Write([]byte(strconv.FormatInt(int64(acc.Balance), 10)))
		if err != nil {
			log.Print(err)
			return err
		}
		_, err = file.Write([]byte(string("|")))
		if err != nil {
			log.Print(err)
			return err
		}
	}
	return nil
}

func (s *Service) ImportFromFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		log.Print(err)
		return err
	}

	defer func() {
		oerr := file.Close()
		if oerr != nil {
			log.Print(oerr)
		}
	}()

	text := make([]byte, 0)
	buf := make([]byte, 4)
	for {
		read, err := file.Read(buf)
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Print(err)
			return err
		}
		text = append(text, buf[:read]...)
	}

	accounts := strings.Split(string(text), "|")
	accounts = accounts[:len(accounts)-1]

	for _, acc := range accounts {
		values := strings.Split(acc, ";")

		id, err := strconv.ParseInt(values[0], 10, 64)
		if err != nil {
			log.Print(err)
			return err
		}

		phone := values[1]

		balance, err := strconv.ParseInt(values[2], 10, 64)
		if err != nil {
			log.Print(err)
			return err
		}
		account := &types.Account{
			ID:      id,
			Phone:   types.Phone(phone),
			Balance: types.Money(balance),
		}

		s.accounts = append(s.accounts, account)
	}

	return nil

}
