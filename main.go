package main

import (
	"github.com/FishZe/Go-DDQuery/account"
	"github.com/FishZe/Go-DDQuery/picMaker"
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
	DDQuery("夜然z", 0)
}
