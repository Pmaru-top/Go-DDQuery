package account

import (
	"Go-DDQuery/api"
	"errors"
	"strconv"
)

type User struct {
	UID        int64   `json:"uid"`
	Name       string  `json:"username"`
	Rank       int     `json:"rank"`
	Sex        string  `json:"sex"`
	Face       string  `json:"face"`
	Sign       string  `json:"sign"`
	Coins      float64 `json:"coins"`
	RegTime    int     `json:"regtime"`
	Fans       int     `json:"fans"`
	Attention  int     `json:"attention"`
	CurrentExp int     `json:"current_exp"`
	Attentions []int   `json:"attentions"`
}

func SearchName(Name string) int64 {
	if Name == "" {
		return 0
	}
	Users, PageNums := api.SearchUser(Name, 1)
	var AllUsers []api.JsonSearchUser
	for i := 1; i <= PageNums; i++ {
		for _, User := range Users {
			if User.Uname == Name {
				return int64(User.Mid)
			}
			AllUsers = append(AllUsers, User)
		}
		Users, _ = api.SearchUser(Name, i+1)
	}
	return 0
}

func (u *User) GetUser() error {
	err := errors.New("")
	if u.UID == 0 && u.Name == "" {
		return errors.New("UID and Name can't be empty at the same time")
	}
	if u.UID == 0 {
		u.UID = SearchName(u.Name)
	}
	if u.UID == 0 {
		return errors.New("user not found")
	}
	UserInfo := api.GetUserInfo(u.UID)
	u.Rank, err = strconv.Atoi(UserInfo.Card.Rank)
	if err != nil {
		return err
	}
	u.Name = UserInfo.Card.Name
	u.Sex = UserInfo.Card.Sex
	u.Face = UserInfo.Card.Face
	u.Sign = UserInfo.Card.Sign
	u.Coins = UserInfo.Card.Coins
	u.RegTime = UserInfo.Card.Regtime
	u.Fans = UserInfo.Card.Fans
	u.Attention = UserInfo.Card.Attention
	u.CurrentExp = UserInfo.Card.LevelInfo.CurrentExp
	u.Attentions = UserInfo.Card.Attentions
	return nil
}
