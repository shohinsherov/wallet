package main

import (
	//	"io"
	"log"
	//	"os"
	//	"strings"
	//"fmt"
	"github.com/shohinsherov/wallet/pkg/wallet"
)

func main() {
	svc := &wallet.Service{}

	svc.ImportFromFile("data/test1.txt")

	acc, err := svc.FindAccountByID(2)
	if err != nil {
		log.Print(err)
	}

	log.Print(acc)
	/*_, err := svc.RegisterAccount("+992000000001")
	if err != nil {
		log.Print(err)
		return
	}
	//fmt.Println(acc)
	_, err = svc.RegisterAccount("+992000000002")
	if err != nil {
		log.Print(err)
		return
	}

	err = svc.ExportToFile("data/test1.txt")
	if err != nil {
		log.Print(err)
		return
	}*/

	/*file, err := os.Open("data/test1.txt")
	if err != nil {
		log.Print(err)
		return
	}

	text := make([]byte, 0)
	buf := make([]byte, 4)
	for {
		read, err := file.Read(buf)
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Print(err)
			return
		}
		text = append(text, buf[:read]...)
	}

		log.Print(text)
		d := string(text)
		log.Print(d)
		accounts := strings.Split(string(text), "|")
		log.Print(accounts)

		accounts = accounts[:len(accounts)-1]
		log.Print(accounts)
		for _, account := range accounts {
			val := strings.Split(account, ";")
			log.Print(val)
		}
		n := len(accounts)
		log.Print(n)
	//data := string(content)
	//	s := strings.Split(data, "|")
	///log.Print(data)
	//log.Print(s)

	/*text := "1,2,3,4"

	fmt.Println(s)
	fmt.Printf("%q\n", strings.Split(text, ","))*/

	/*account.Balance = 200
	fmt.Println(account)
	payment, err := svc.Pay(1, 100, "car")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(payment)
	fmt.Println(account)

	errR := svc.Reject(payment.ID)
	if errR != nil {
		fmt.Println(errR)
		return
	}

	fmt.Println(account)

	err = svc.Deposit(account.ID, 10)
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
	//svc.RegisterAccount("+992000000002")

	file, err := os.Create("data/test.txt")
	if err != nil {
		log.Print(err)
		return
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Print(err)
		}
	}()

	//log.Printf("%#v", file)

	_, err = file.Write([]byte("Hello, test its my text"))
	if err != nil {
		log.Print(err)
		return
	}
	content := make([]byte, 0)
	buf := make([]byte, 4)
	for {
		read, err := file.Read(buf)
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Print(err)
			return
		}

		content = append(content, buf[:read]...)
	}

	data := string(content)
	log.Print(data)*/
}
