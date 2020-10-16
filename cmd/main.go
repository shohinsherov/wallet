package main

import (
	//	"io/ioutil"
	//"github.com/shohinsherov/wallet/pkg/types"
	"log"
	//"path/filepath"
	//"os"
	//"strings"
	//"fmt"
	"github.com/shohinsherov/wallet/pkg/wallet"
)

func main() {

	svc := &wallet.Service{}

	svc.RegisterAccount("001")
	svc.Deposit(1, 1000)
/*	svc.Pay(1,1,"asd")
	svc.Pay(1,2,"asd")
	svc.Pay(1,3,"asd")
	svc.Pay(1,4,"asd")
	svc.Pay(1,5,"asd")
	svc.Pay(1,6,"asd")
	svc.Pay(1,7,"asd")
	svc.Pay(1,8,"asd")*/

	pay, _ := svc.ExportAccountHistory(1)
	if pay != nil {
		log.Print(pay)
	}
	svc.HistoryToFiles(pay, "data", 3)

}
