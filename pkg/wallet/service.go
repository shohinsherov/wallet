package wallet

import (
	"bufio"

	//"io/ioutil"
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

// ErrFavoriteNotFound ..
var ErrFavoriteNotFound = errors.New("Favorite not found")

// ErrNotPayments ...
var ErrNotPayments = errors.New("No payments to export")

// Service ....
type Service struct {
	nextAccountID int64
	accounts      []*types.Account
	payments      []*types.Payment
	favorites     []*types.Favorite
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
		if account.ID > s.nextAccountID {
			s.nextAccountID = account.ID
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

// Repeat ...
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

// FavoritePayment ..
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

// FindFavoriteByID ..
func (s *Service) FindFavoriteByID(favoriteID string) (*types.Favorite, error) {
	for _, favorite := range s.favorites {
		if favoriteID == favorite.ID {
			return favorite, nil
		}
	}
	return nil, ErrFavoriteNotFound
}

// PayFromFavorite ..
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

// ExportToFile ...
func (s *Service) ExportToFile(path string) error {
	file, err := os.Create(path)
	if err != nil {
		log.Print(err)
		return err
	}

	for _, acc := range s.accounts {
		text := []byte(strconv.FormatInt(int64(acc.ID), 10) + string(";") + string(acc.Phone) + string(";") + strconv.FormatInt(int64(acc.Balance), 10) + string("|"))

		_, err = file.Write(text)
		if err != nil {
			log.Print(err)
			return err
		}

	}
	return nil
}

// ImportFromFile ...
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

		s.accounts = append(s.accounts, &types.Account{
			ID:      id,
			Phone:   types.Phone(phone),
			Balance: types.Money(balance),
		})

	}

	return nil

}

// Export ...
func (s *Service) Export(dir string) error {
	// ---- accounts
	if s.accounts != nil {
		file, err := os.Create(dir + "/accounts.dump")
		if err != nil {
			log.Print(err)
			return err
		}
		for _, acc := range s.accounts {
			text := []byte(strconv.FormatInt(int64(acc.ID), 10) + string(";") + string(acc.Phone) + string(";") + strconv.FormatInt(int64(acc.Balance), 10) + string(";") + string('\n'))
			_, err = file.Write(text)
			if err != nil {
				return err
			}
		}
		log.Print("Аккаунты экпортированы")
	} else {
		log.Print("Нет аккаунтов для экпорта")
	}

	// ---- payments
	if s.payments != nil {
		payFile, err := os.Create(dir + "/payments.dump")
		if err != nil {
			log.Print(err)
			return err
		}

		for _, pay := range s.payments {
			text := []byte(string(pay.ID) + string(";") + strconv.FormatInt(int64(pay.AccountID), 10) + string(";") + strconv.FormatInt(int64(pay.Amount), 10) + string(";") + string(pay.Category) + string(";") + string(pay.Status) + string(";") + string('\n'))
			_, err = payFile.Write(text)
			if err != nil {
				log.Print(err)
				return err
			}
		}
		log.Print("Платежи экпортированы")
	} else {
		log.Print("Нет платежей для экспорта")
	}

	// ---- favorites
	if s.favorites != nil {
		favFile, err := os.Create(dir + "/favorites.dump")
		if err != nil {
			log.Print(err)
			return err
		}

		for _, fav := range s.favorites {
			text := []byte(fav.ID + ";" + strconv.FormatInt(int64(fav.AccountID), 10) + ";" + fav.Name + ";" + strconv.FormatInt(int64(fav.Amount), 10) + ";" + string(fav.Category) + ";" + string('\n'))
			_, err := favFile.Write(text)
			if err != nil {
				log.Print(err)
				return err
			}
		}
		log.Print("Избранные экпортированы")
	} else {
		log.Print("Нет избранных для экспорта")
	}

	return nil
}

// Import ...
func (s *Service) Import(dir string) error {
	// ---- accounts
	src, err := os.Open(dir + "/accounts.dump")
	if err != nil {
		log.Print("Dump аккаунтов не найден")
	} else {
		defer func() {
			if cerr := src.Close(); cerr != nil {
				log.Print(cerr)
			}
		}()

		reader := bufio.NewReader(src)
		for {
			line, err := reader.ReadString('\n')
			if err == io.EOF {
				log.Print(line)
				break
			}
			if err != nil {
				log.Print(err)
				return err
			}

			val := strings.Split(line, ";")

			id, err := strconv.ParseInt(val[0], 10, 64)
			if err != nil {
				log.Print(err)
				return err
			}

			phone := val[1]

			balance, err := strconv.ParseInt(val[2], 10, 64)
			if err != nil {
				log.Print(err)
				return err
			}

			findacc, _ := s.FindAccountByID(id)
			if findacc != nil {
				findacc.Phone = types.Phone(phone)
				findacc.Balance = types.Money(balance)
			} else {
				s.nextAccountID = id
				newAcc := &types.Account{
					ID:      s.nextAccountID,
					Phone:   types.Phone(phone),
					Balance: types.Money(balance),
				}

				s.accounts = append(s.accounts, newAcc)
			}
		}
		log.Print("Аккаунты импортированы")
		/*for _, a := range s.accounts {
			log.Print(a)
		}*/
	}

	// ---- payments
	paySrc, err := os.Open(dir + "/payments.dump")
	if err != nil {
		log.Print("Dump платежей не найден")
	} else {
		defer func() {
			if cerr := paySrc.Close(); cerr != nil {
				log.Print(cerr)
			}
		}()

		payReader := bufio.NewReader(paySrc)
		for {
			payLine, err := payReader.ReadString('\n')
			if err == io.EOF {
				log.Print(payLine)
				break
			}
			if err != nil {
				log.Print(err)
				return err
			}

			val := strings.Split(payLine, ";")

			id := string(val[0])
			accID, err := strconv.ParseInt(val[1], 10, 64)
			if err != nil {
				log.Print(err)
				return err
			}

			amount, err := strconv.ParseInt(val[2], 10, 64)
			if err != nil {
				log.Print(err)
				return err
			}

			category := val[3]

			status := val[4]

			findPay, _ := s.FindPaymentByID(id)
			if findPay != nil {
				findPay.AccountID = accID
				findPay.Amount = types.Money(amount)
				findPay.Category = types.PaymentCategory(category)
				findPay.Status = types.PaymentStatus(status)
			} else {
				newPay := &types.Payment{
					ID:        id,
					AccountID: accID,
					Amount:    types.Money(amount),
					Category:  types.PaymentCategory(category),
					Status:    types.PaymentStatus(status),
				}

				s.payments = append(s.payments, newPay)
			}
		}
		log.Print("Платежы импортированы")
		/*for _, p := range s.payments {
			log.Print(p)
		}*/
	}

	// ---- favorites
	favFile, err := os.Open(dir + "/favorites.dump")
	if err != nil {
		log.Print("Dump избранных не найден")
	} else {
		reader := bufio.NewReader(favFile)
		for {
			favLine, err := reader.ReadString('\n')
			if err == io.EOF {
				log.Print(favLine)
				break
			}
			if err != nil {
				log.Print(err)
				return err
			}

			val := strings.Split(favLine, ";")

			id := val[0]
			accID, err := strconv.ParseInt(val[1], 10, 64)
			if err != nil {
				log.Print(err)
				return err
			}
			name := val[2]
			amount, err := strconv.ParseInt(val[3], 10, 64)
			if err != nil {
				log.Print(err)
				return err
			}
			category := val[4]

			findFav, _ := s.FindFavoriteByID(id)
			if findFav != nil {
				findFav.AccountID = accID
				findFav.Amount = types.Money(amount)
				findFav.Name = name
				findFav.Category = types.PaymentCategory(category)
			} else {
				newFav := &types.Favorite{
					ID:        id,
					AccountID: accID,
					Name:      name,
					Amount:    types.Money(amount),
					Category:  types.PaymentCategory(category),
				}
				s.favorites = append(s.favorites, newFav)
			}
		}
		log.Print("Избранные ипортированы")

		for _, f := range s.favorites {
			log.Print(f)
		}
	}

	return nil
}

// ExportAccountHistory вытаскивает все платежи конкретного аккаунта
func (s *Service) ExportAccountHistory(accountID int64) ([]types.Payment, error) {
	_, err := s.FindAccountByID(accountID)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	payment := []types.Payment{}
	for _, pay := range s.payments {
		if pay.AccountID == accountID {
			payment = append(payment, types.Payment{
				ID:        pay.ID,
				AccountID: pay.AccountID,
				Amount:    pay.Amount,
				Category:  pay.Category,
				Status:    pay.Status})
		}
	}
	//a := len(payment)
	//log.Print(a)
	//log.Print(payment)
	if len(payment) == 0 {
		return nil, ErrNotPayments //fmt.Errorf("%v",ErrNotPayments)
	}

	return payment, nil

}

// HistoryToFiles ...
func (s *Service) HistoryToFiles(payments []types.Payment, dir string, records int) error {
	if payments == nil {
		log.Print(ErrNotPayments)
		return nil
	}
	if len(payments) < records {
		file, err := os.Create(dir + "/payments.dump")
		if err != nil {
			log.Print(err)
			return err
		}
		for _, payment := range payments {
			text := []byte(payment.ID + ";" + strconv.FormatInt(int64(payment.AccountID), 10) + ";" + strconv.FormatInt(int64(payment.Amount), 10) + ";" + string(payment.Category) + ";" + string(payment.Status) + ";" + string('\n'))
			_, err = file.Write(text)
			if err != nil {
				log.Print(err)
				return err
			}
		}
		return nil
	}
	/*for i := 0; i <len(payments); i++ {
		log.Print(payments[1].ID)
	}*/

	n := len(payments) / records

	for k := 0; k < n; k++ {
		v := strconv.FormatInt(int64(k+1), 10)
		d := dir + "/payments" + string(v) + ".dump"
		file, err := os.Create(d)
		if err != nil {
			log.Print(err)
			return err
		}

		for i := records * k; i < records*(k+1); i++ {
			text := []byte(payments[i].ID + ";" + strconv.FormatInt(int64(payments[i].AccountID), 10) + ";" + strconv.FormatInt(int64(payments[i].Amount), 10) + ";" + string(payments[i].Category) + ";" + string(payments[i].Status) + ";" + string('\n'))
			_, err = file.Write(text)
			if err != nil {
				log.Print(err)
				return err
			}

		}

	}
	if n*records < len(payments) {
		v := strconv.FormatInt(int64(n+1), 10)
		d := dir + "/payments" + string(v) + ".dump"
		file, err := os.Create(d)
		if err != nil {
			log.Print(err)
			return err
		}
		for i := records * n; i < len(payments); i++ {
			text := []byte(payments[i].ID + ";" + strconv.FormatInt(int64(payments[i].AccountID), 10) + ";" + strconv.FormatInt(int64(payments[i].Amount), 10) + ";" + string(payments[i].Category) + ";" + string(payments[i].Status) + ";" + string('\n'))
			_, err = file.Write(text)
			if err != nil {
				log.Print(err)
				return err
			}

		}

	}

	return nil
}
