package account

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/FishZe/Go-DDQuery/api"
	"io"
	"log"
	"os"
	"time"
)

type User struct {
	UID           int64     `json:"uid"`
	Name          string    `json:"username"`
	Rank          int       `json:"rank"`
	Sex           string    `json:"sex"`
	Face          string    `json:"face"`
	FaceFile      string    `json:"face_file"`
	Sign          string    `json:"sign"`
	Coins         float64   `json:"coins"`
	RegTime       int       `json:"regtime"`
	Fans          int       `json:"fans"`
	Attention     int       `json:"attention"`
	CurrentExp    int       `json:"current_exp"`
	Attentions    []int64   `json:"attentions"`
	VupAttentions []Vup     `json:"vup_attentions"`
	UserGuard     UserGuard `json:"user_guard"`
}

type Vup struct {
	Mid    int64  `json:"mid,omitempty"`
	Uname  string `json:"uname,omitempty"`
	RoomId int64  `json:"roomid,omitempty"`
	Living bool   `json:"living,omitempty"`
}

type UserGuard struct {
	Uname string    `json:"uname"`
	Face  string    `json:"face"`
	Mid   int64     `json:"mid"`
	Dd    [][]int64 `json:"dd"`
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

func checkDataUpdate(fileName string, download func()) error {
	if !pathExists(fileName) {
		if !pathExists("data") {
			err := os.Mkdir("data", 0777)
			if err != nil {
				log.Printf("os.Mkdir Error: %v", err)
				return err
			}
		}
		download()
		return nil
	}
	fileInfo, err := os.Stat(fileName)
	if err != nil {
		return err
	}
	if time.Now().Unix()-fileInfo.ModTime().Unix() > 24*60*60 {
		download()
	}
	return nil
}

func getVup() (map[int64]Vup, error) {
	err := checkDataUpdate("data/vup.json", api.DownloadVupJson)
	if err != nil {
		return nil, err
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

func getGuard(user *User) (UserGuard, error) {
	err := checkDataUpdate("data/guard.json", api.DownloadGuardJson)
	if err != nil {
		log.Printf("checkDataUpdate Error: %v", err)
	}
	updateTime, err := api.GetUserGuardUpdateTime()
	if err != nil {
		updateTime = 0
		log.Printf("GetUserGuardUpdateTime Error: %v", err)
	}
	if time.Now().Unix()-updateTime > 3*24*60*60 {
		return UserGuard{Uname: user.Name, Face: "", Mid: user.UID, Dd: make([][]int64, 3)}, nil
	}
	jsonFile, err := os.Open("data/guard.json")
	if err != nil {
		fmt.Println(err)
		return UserGuard{Uname: user.Name, Face: "", Mid: user.UID, Dd: make([][]int64, 3)}, err
	}
	defer func(jsonFile *os.File) {
		err := jsonFile.Close()
		if err != nil {
			fmt.Println("jsonFile.Close() Error!")
		}
	}(jsonFile)
	byteValue, _ := io.ReadAll(jsonFile)
	var result map[int64]UserGuard
	err = json.Unmarshal(byteValue, &result)
	if err != nil {
		return UserGuard{Uname: user.Name, Face: "", Mid: user.UID, Dd: make([][]int64, 3)}, err
	}
	if _, ok := result[user.UID]; ok {
		return result[user.UID], nil
	}
	return UserGuard{Uname: user.Name, Face: "", Mid: user.UID, Dd: make([][]int64, 3)}, nil
}

func SearchName(Name string) (int64, error) {
	if Name == "" {
		return 0, errors.New("name can't be empty")
	}
	Users, PageNums, err := api.SearchUser(Name, 1)
	if err != nil {
		return 0, err
	}
	var AllUsers []api.JsonSearchUser
	for i := 1; i <= PageNums; i++ {
		for _, User := range Users {
			if User.Uname == Name {
				return int64(User.Mid), nil
			}
			AllUsers = append(AllUsers, User)
		}
		Users, _, err = api.SearchUser(Name, i+1)
		if err != nil {
			log.Printf("SearchUser Error: %v", err)
			continue
		}
	}
	return 0, nil
}

func (u *User) GetUser() error {
	if u.UID == 0 && u.Name == "" {
		return errors.New("UID and Name can't be empty at the same time")
	}
	var err error
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
	if len(UserInfo.Card.Attentions) == 0 {
		return errors.New("user has no attention")
	}
	vups, err := getVup()
	if err != nil {
		return err
	}
	u.FaceFile, err = api.DownloadFace(UserInfo.Card.Face, u.UID)
	if err != nil {
		log.Printf("DownloadFace Error: %v", err)
	}
	if err != nil {
		return err
	}
	u.Rank = UserInfo.Card.LevelInfo.CurrentLevel
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
	if len(u.VupAttentions) == 0 {
		return errors.New("user has no vup attention")
	}
	LivingRoom := api.GetLivingRoom()
	for i, vup := range u.VupAttentions {
		u.VupAttentions[i].Living = false
		for _, room := range LivingRoom {
			if vup.RoomId == room {
				u.VupAttentions[i].Living = true
				break
			}
		}
	}
	u.UserGuard, err = getGuard(u)
	if err != nil {
		log.Printf("getGuard Error: %v", err)
	}
	return nil
}
