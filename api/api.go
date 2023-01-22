package api

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

var Cookies string

func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func getReq(data url.Values, getUrl string, cookies string) ([]byte, string, error) {
	u, err := url.ParseRequestURI(getUrl)
	if err != nil {
		return nil, "", err
	}
	u.RawQuery = data.Encode()
	client := http.Client{}
	req, err := http.NewRequest("GET", u.String(), nil)
	req.Header = http.Header{
		"accept":     {"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9"},
		"User-Agent": {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36"},
		"cookie":     {cookies},
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			log.Printf("Error orrured when closing the session: %v", err)
		}
	}(resp.Body)
	if resp.StatusCode != 200 {
		return nil, "", errors.New("status code error: " + strconv.Itoa(resp.StatusCode) + " " + resp.Status)
	}
	if resp.Header.Get("Set-Cookie") != "" {
		cookies = ""
		for _, v := range resp.Header["Set-Cookie"] {
			cookies += v + ";"
		}
	} else {
		cookies = ""
	}
	s, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}
	return s, cookies, nil
}

func downloadFile(URL string, fileName string) error {
	response, err := http.Get(URL)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("Error occured when closing the session: %v", err)
		}
	}(response.Body)
	if response.StatusCode != 200 {
		return errors.New("status code error: " + strconv.Itoa(response.StatusCode) + " " + response.Status)
	}
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Printf("Error occured when closing the file: %v", err)
		}
	}(file)
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}
	return nil
}

func SearchUser(KeyWord string, page int) ([]JsonSearchUser, int, error) {
	if KeyWord == "" {
		return nil, 0, errors.New("KeyWord is empty")
	}
	err := error(nil)
	if Cookies == "" {
		Cookies = GetCookies()
		if err != nil {
			return nil, 0, err
		}
	}
	getUrl := "https://api.bilibili.com/x/web-interface/search/type"
	data := url.Values{}
	data.Set("search_type", "bili_user")
	data.Set("keyword", KeyWord)
	data.Set("page", strconv.Itoa(page))
	s, _, err := getReq(data, getUrl, Cookies)
	if err != nil {
		return nil, 0, err
	}
	var jSR JsonSearchRes
	err = json.Unmarshal(s, &jSR)
	if err != nil {
		return nil, 0, err
	}
	return jSR.Data.Result, jSR.Data.NumPages, nil
}

func GetCookies() string {
	getUrl := "https://www.bilibili.com/"
	data := url.Values{}
	data.Set("spm_id_from", "333.999.0.0")
	_, cookies, err := getReq(data, getUrl, "")
	if err != nil {
		log.Printf("Error occured when sending GET request: %v", err)
		return ""
	}
	return cookies
}

func GetUserInfo(uid int64) JsonUserInfo {
	err := error(nil)
	if Cookies == "" {
		Cookies = GetCookies()
	}
	getUrl := "https://account.bilibili.com/api/member/getCardByMid"
	data := url.Values{}
	data.Set("mid", strconv.FormatInt(uid, 10))
	s, _, err := getReq(data, getUrl, Cookies)
	if err != nil {
		log.Printf("Error occured when sending GET request: %v", err)
		return JsonUserInfo{}
	}
	var jUI JsonUserInfo
	err = json.Unmarshal(s, &jUI)
	if err != nil {
		log.Printf("Error occured when unmarshal the json: %v", err)
		return JsonUserInfo{}
	}
	return jUI
}

func DownloadFace(FaceUrl string, uid int64) (string, error) {
	if FaceUrl == "" {
		return "", errors.New("FaceUrl is empty")
	}
	if !PathExists("data/face") {
		err := os.Mkdir("data/face", 0777)
		if err != nil {
			return "", err
		}
	}
	fileName := "data/face/" + strconv.FormatInt(uid, 10) + ".jpg"
	err := downloadFile(FaceUrl, fileName)
	if err != nil {
		log.Printf("Error occured when downloading the file: %v", err)
		return "", err
	}
	return fileName, nil
}

func GetLivingRoom() []int64 {
	getUrl := "https://api.vtbs.moe/v1/living"
	data := url.Values{}
	s, _, err := getReq(data, getUrl, "")
	if err != nil {
		fmt.Println(err)
		return []int64{}
	}
	var LR []int64
	err = json.Unmarshal(s, &LR)
	if err != nil {
		fmt.Println(err)
		return []int64{}
	}
	return LR
}

func DownloadVupJson() {
	getUrl := "https://api.vtbs.moe/v1/short"
	fileName := "data/vup.json"
	err := downloadFile(getUrl, fileName)
	if err != nil {
		log.Printf("Error occured when downloading the file: %v", err)
	}
}

func DownloadVupInfoJson() {
	getUrl := "https://vdb.vtbs.moe/json/fs.json"
	fileName := "data/vup_info.json"
	err := downloadFile(getUrl, fileName)
	if err != nil {
		log.Printf("Error occured when downloading the file: %v", err)
	}
}

func DownloadFont() {
	getUrl := "https://git.fishze.top/https://raw.githubusercontent.com/adobe-fonts/source-han-sans/release/Variable/TTF/SourceHanSansSC-VF.ttf"
	fileName := "data/font/SourceHanSansSC-VF.ttf"
	err := downloadFile(getUrl, fileName)
	if err != nil {
		log.Printf("Error occured when downloading the file: %v", err)
	}
}

func GetUserGuardUpdateTime() (int64, error) {
	getUrl := "https://api.vtbs.moe/v1/guard/time"
	data := url.Values{}
	t, _, err := getReq(data, getUrl, "")
	if err != nil {
		return 0, err
	}
	return int64(binary.BigEndian.Uint64(t)), nil
}

func DownloadGuardJson() {
	getUrl := "https://api.vtbs.moe/v1/guard/all"
	fileName := "data/guard.json"
	err := downloadFile(getUrl, fileName)
	if err != nil {
		log.Printf("Error occured when downloading the file: %v", err)
	}
}

func guardIcon2Write(fileName string) {
	guardIcon, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Error occured when opening the file: ", err)
		return
	}
	defer func(pngImgFile *os.File) {
		err := pngImgFile.Close()
		if err != nil {
			fmt.Println("Error occured when closing the file: ", err)
		}
	}(guardIcon)
	rawImg, err := png.Decode(guardIcon)
	if err != nil {
		fmt.Println("Error occured when decoding the image: ", err)
		return
	}
	newImg := image.NewRGBA(rawImg.Bounds())
	draw.Draw(newImg, newImg.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)
	draw.Draw(newImg, newImg.Bounds(), rawImg, rawImg.Bounds().Min, draw.Over)
	newImgFile, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error occured when creating the file: ", err)
		return
	}
	defer func(jpgImgFile *os.File) {
		err := jpgImgFile.Close()
		if err != nil {
			fmt.Println("Error occured when closing the file: ", err)
		}
	}(newImgFile)
	err = png.Encode(newImgFile, newImg)
	if err != nil {
		fmt.Println("Error occured when encoding the image: ", err)
		return
	}
}

func DownloadGuardIcon() {
	getUrl := map[int]string{}
	getUrl[0] = "https://i0.hdslb.com/bfs/live/1d16bf0fcc3b1b768d1179d60f1fdbabe6ab4489.png"
	getUrl[1] = "https://i0.hdslb.com/bfs/live/98a201c14a64e860a758f089144dcf3f42e7038c.png"
	getUrl[2] = "https://i0.hdslb.com/bfs/live/143f5ec3003b4080d1b5f817a9efdca46d631945.png"
	for i, j := range getUrl {
		fileName := "data/icon/g" + strconv.Itoa(i) + ".png"
		err := downloadFile(j, fileName)
		if err != nil {
			log.Printf("Error occured when downloading the file: %v", err)
			continue
		}
		guardIcon2Write(fileName)
	}
}

func DownloadIcons() {
	getUrl := map[string]string{}
	for i := 0; i <= 6; i++ {
		getUrl["l"+strconv.Itoa(i)] = "https://git.fishze.top/https://raw.githubusercontent.com/FishZe/Go-DDQuery/master/data/icon/l" + strconv.Itoa(i) + ".png"
	}
	getUrl["living"] = "https://git.fishze.top/https://raw.githubusercontent.com/FishZe/Go-DDQuery/master/data/icon/living.png"
	for i, j := range getUrl {
		fileName := "data/icon/" + i + ".png"
		err := downloadFile(j, fileName)
		if err != nil {
			log.Printf("Error occured when downloading the file: %v", err)
			continue
		}
	}
}
