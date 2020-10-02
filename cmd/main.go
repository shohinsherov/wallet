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
	fmt.Println(account)
	account, err = svc.FindAccountByID(2)
	//if err != nil {
		fmt.Println(err)
	//}
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
