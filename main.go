package main

import (
	"Go-DDQuery/account"
	"Go-DDQuery/picMaker"
	"fmt"
	"os/exec"
	"strconv"
)

func DDQuery(Name string, Uid int64) (string, error) {
	user := account.User{Name: Name, UID: Uid}
	err := user.GetUser()
	if err != nil {
		return "", err
	}
	route, err := picMaker.MkPic(user)
	if err != nil {
		return "", err
	}
	return route, nil
}

func main() {
	var input string
	fmt.Println("请输入要查询的昵称或UID (如: 陈睿 / UID: 1):")
	_, err := fmt.Scanln(&input)
	if err != nil {
		return
	}
	var uid int64
	var name string
	if (len(input) > 4 && (input[:4] == "UID:" || input[:4] == "uid:")) || (len(input) > 5 && (input[:4] == "UID：" || input[:4] == "uid：")) {
		input = input[4:]
		uid, err = strconv.ParseInt(input, 10, 64)
		if err != nil {
			fmt.Println("输入的UID格式错误")
			return
		}
	} else {
		name = input
	}
	route, err := DDQuery(name, uid)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(route)
	cmd := exec.Command("cmd", "/k", "start", route)
	err = cmd.Start()
	if err != nil {
		return
	}
}
