package main

import (
	"fmt"

	"github.com/shohinsherov/wallet/pkg/wallet"
)

func main() {
	svc := &wallet.Service{}
	account, err := svc.RegisterAccount("+992000000001")
	if err != nil {
		fmt.Println(err)
		return
	}
	account.Balance = 200
	fmt.Println(account)
	payment, err := svc.Pay(1, 100, "car")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(payment)
	fmt.Println(account)

	rejact, err := svc.Reject(payment.ID)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(rejact)
	fmt.Println(account)

	/*err = svc.Deposit(account.ID, 10)
	if err != nil {
		switch err {
		case wallet.ErrAmountMustBePositive:
			fmt.Println("Сумма должна быть больше 0")
		case wallet.ErrAccountNotFound:
			fmt.Println("Аккаунт не найден")
		}
		return
	}
	fmt.Println(account.Balance)
	//svc.RegisterAccount("+992000000002")*/
}
