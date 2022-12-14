package main

import (
	"Go-DDQuery/account"
	"Go-DDQuery/picMaker"
	"fmt"
	"os/exec"
	"strconv"
)

func main() {
	var input string
	fmt.Println("请输入要查询的昵称或UID (如: 陈睿 / UID: 1):")
	_, err2 := fmt.Scanln(&input)
	if err2 != nil {
		return
	}
	user := account.User{}
	if input[:5] == "UID: " {
		input = input[5:]
		uid, err := strconv.ParseInt(input, 10, 64)
		if err != nil {
			fmt.Println("输入的UID格式错误")
			return
		}
		user = account.User{Name: "", UID: uid}
	} else {
		user = account.User{Name: input}
	}
	err := user.GetUser()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(user)
	route := picMaker.MkPic(user)
	fmt.Println(route)
	cmd := exec.Command("cmd", "/k", "start", route)
	err = cmd.Start()
	if err != nil {
		return
	}
}
