package main

import (
	"Go-DDQuery/account"
	"fmt"
)

func main() {
	user := account.User{Name: "陈睿"}
	err := user.GetUser()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(user)
}
