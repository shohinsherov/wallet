package types

// Money Денежная сумма
type Money int64

// Currency валюта
type Currency string

// Коды валют
const (
	TJS Currency = "TJS"
	RUB Currency = "RUB"
	USD Currency = "USD"
)

// PAN номер карты
type PAN string

// Card платежная карта
type Card struct {
	ID         int
	PAN        PAN
	Balance    Money
	MinBalance Money
	Currency   Currency
	Color      string
	Name       string
	Active     bool
}

// PaymentCategory категория
type PaymentCategory string

// PaymentStatus статус платежа
type PaymentStatus string

// статусы платежа
const (
	PaymentStatusOk         PaymentStatus = "OK"
	PaymentStatusFail       PaymentStatus = "FAIL"
	PaymentStatusInProgress PaymentStatus = "INPROGRESS"
)

// Payment информация о платеже
type Payment struct {
	ID        string
	AccountID int64
	Amount    Money
	Category  PaymentCategory
	Status    PaymentStatus
}

// Phone телефон
type Phone string

// Account счет пользователья
type Account struct {
	ID      int64
	Phone   Phone
	Balance Money
}
