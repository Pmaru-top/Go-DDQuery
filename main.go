package main

import (
	"Go-DDQuery/account"
	"Go-DDQuery/picMaker"
	"fmt"
	"os/exec"
)

func main() {
	user := account.User{Name: "突进吧蔗糖"}
	err := user.GetUser()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(user)
	route := picMaker.MkPic(user)
	cmd := exec.Command("cmd", "/k", "start", route)
	err = cmd.Start()
	if err != nil {
		return
	}
}
