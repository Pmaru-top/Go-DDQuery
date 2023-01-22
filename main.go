package Go_DDQuery

import (
	"fmt"

	"github.com/Pmaru-top/Go-DDQuery/account"
	"github.com/Pmaru-top/Go-DDQuery/picMaker"
)

var uid int64

func DDQuery(Name string, Uid int64) (bytes []byte, err error) {
	user := account.User{Name: Name, UID: Uid}
	err = user.GetUser()
	if err != nil {
		return
	}

	uid = Uid
	bytes,err = picMaker.MkPic(user)

	if err != nil{
		fmt.Println("獲取失敗:\n",err)
	}

	return
}

func SavePic(path string, bytes []byte) (directory string,err error) {
	directory = fmt.Sprint(path, "/", uid, ".png")
	err = picMaker.SavePic(directory, bytes)
	return
}
