package account

import (
	"Go-DDQuery/api"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
)

type User struct {
	UID           int64   `json:"uid"`
	Name          string  `json:"username"`
	Rank          int     `json:"rank"`
	Sex           string  `json:"sex"`
	Face          string  `json:"face"`
	FaceFile      string  `json:"face_file"`
	Sign          string  `json:"sign"`
	Coins         float64 `json:"coins"`
	RegTime       int     `json:"regtime"`
	Fans          int     `json:"fans"`
	Attention     int     `json:"attention"`
	CurrentExp    int     `json:"current_exp"`
	Attentions    []int64 `json:"attentions"`
	VupAttentions []Vup   `json:"vup_attentions"`
}

type Vup struct {
	Mid    int64  `json:"mid,omitempty"`
	Uname  string `json:"uname,omitempty"`
	RoomId int    `json:"roomid,omitempty"`
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func getVup() (map[int64]Vup, error) {
	if !pathExists("data/vup.json") {
		if !pathExists("data") {
			err := os.Mkdir("data", 0777)
			if err != nil {
				return nil, err
			}
		}
		api.DownloadVupJson()
	}
	jsonFile, err := os.Open("data/vup.json")
	if err != nil {
		fmt.Println(err)
		return map[int64]Vup{}, err
	}
	defer func(jsonFile *os.File) {
		err := jsonFile.Close()
		if err != nil {
			fmt.Println("jsonFile.Close() Error!")
		}
	}(jsonFile)
	byteValue, _ := io.ReadAll(jsonFile)
	var result []Vup
	err = json.Unmarshal(byteValue, &result)
	if err != nil {
		return map[int64]Vup{}, err
	}
	var vupMap = make(map[int64]Vup)
	for _, vup := range result {
		vupMap[vup.Mid] = vup
	}
	return vupMap, nil
}

func SearchName(Name string) (int64, error) {
	if Name == "" {
		return 0, errors.New("name can't be empty")
	}
	Users, PageNums := api.SearchUser(Name, 1)
	var AllUsers []api.JsonSearchUser
	for i := 1; i <= PageNums; i++ {
		for _, User := range Users {
			if User.Uname == Name {
				return int64(User.Mid), nil
			}
			AllUsers = append(AllUsers, User)
		}
		Users, _ = api.SearchUser(Name, i+1)
	}
	return 0, nil
}

func (u *User) GetUser() error {
	err := errors.New("")
	if u.UID == 0 && u.Name == "" {
		return errors.New("UID and Name can't be empty at the same time")
	}
	if u.UID == 0 {
		u.UID, err = SearchName(u.Name)
		if err != nil {
			return err
		}
	}
	if u.UID == 0 {
		return errors.New("user not found")
	}
	UserInfo := api.GetUserInfo(u.UID)
	if err != nil {
		return err
	}
	u.Rank, err = strconv.Atoi(UserInfo.Card.Rank)
	if err != nil {
		return err
	}
	vups, err := getVup()
	if err != nil {
		return err
	}
	u.FaceFile = api.DownloadFace(UserInfo.Card.Face, u.UID)
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
	for _, vup := range u.Attentions {
		if _, ok := vups[vup]; ok {
			u.VupAttentions = append(u.VupAttentions, vups[vup])
		}
	}
	return nil
}
