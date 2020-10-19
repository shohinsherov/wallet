package main

import (
	//	"io/ioutil"
	"github.com/shohinsherov/wallet/pkg/types"
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
	svc.Pay(1,10,"1")
	svc.Pay(1,10,"2")
	svc.Pay(1,10,"3")
	svc.Pay(1,10,"4")
	svc.Pay(1,10,"5")
	svc.Pay(1,10,"6")
	svc.Pay(1,10,"7")
	svc.Pay(1,10,"8")
	svc.Pay(1,10,"9")
	svc.Pay(1,10,"10")
	svc.Pay(1,10,"11")
	svc.Pay(1,10,"12")
	svc.Pay(1,10,"13")
	svc.Pay(1,10,"14")
	svc.Pay(1,10,"15")
	svc.Pay(1,10,"16")
	svc.Pay(1,10,"17")
	svc.Pay(1,10,"18")
	svc.Pay(1,10,"19")
	svc.Pay(1,10,"20")
	svc.Pay(1,10,"21")
	svc.Pay(1,10,"22")
	svc.Pay(1,10,"23")
	svc.Pay(1,10,"24")
	svc.Pay(1,10,"25")
	svc.Pay(1,10,"26")
	svc.Pay(1,10,"27")
	svc.Pay(1,10,"28")
	svc.Pay(1,10,"29")
	svc.Pay(1,10,"30")
	//svc.Pay(1,10,"31")
	
	
	ss, _ := svc.FilterPaymentsByFn(func(payment types.Payment) bool {
		if payment.AccountID == 1 {
			return true
		}
		return false
	}, 1)

	log.Print(len(ss))

	

}
