package wallet

import (
	"fmt"
	
	"testing"
)

func TestSerivce_FindAccountByID_success(t *testing.T) {
	svc := &Service{}
	svc.RegisterAccount("+992000000001")

	account, err := svc.FindAccountByID(1)

	fmt.Println(account)
	fmt.Println(err)
	
	// Output:
	// &{1 +992000000001 0} 
	// <nil>	
}

func TestSerivce_FindAccountByID_notfound(t *testing.T) {
	svc := &Service{}
	svc.RegisterAccount("+992000000001")

	account, err := svc.FindAccountByID(2)

	fmt.Println(account)
	fmt.Println(err)
	
	// Output:
	// <nil>
	// Account not found	
	
}