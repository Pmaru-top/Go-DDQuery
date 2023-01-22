package Go_DDQuery

import (
	"github.com/FishZe/Go-DDQuery/account"
	"github.com/FishZe/Go-DDQuery/picMaker"
)

func DDQuery(Name string, Uid int64) (route string, bytes []byte, err error) {
	user := account.User{Name: Name, UID: Uid}
	err = user.GetUser()
	if err != nil {
		return
	}

	bytes = picMaker.MkPic(user)
	// route = path + "/" + strconv.FormatInt(user.UID, 10) + ".png"

	return
}

func SavePic(route string, bytes []byte) (err error) {
	err = picMaker.SavePic(route, bytes)
	if err != nil {
		return
	}
	return
}
