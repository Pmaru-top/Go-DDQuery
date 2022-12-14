package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

var Cookies string

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

func getReq(data url.Values, getUrl string, cookies string) ([]byte, string) {
	u, err := url.ParseRequestURI(getUrl)
	if err != nil {
		panic(err)
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
		panic(err)
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			log.Printf("Error orrured when closing the session: %v", err)
		}
	}(resp.Body)
	if resp.StatusCode != 200 {
		panic("error occurred when sending GET request")
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
		panic(err)
	}
	return s, cookies
}

func downloadFile(URL string, fileName string) {
	response, err := http.Get(URL)
	if err != nil {
		panic(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("Error occured when closing the session: %v", err)
		}
	}(response.Body)
	if response.StatusCode != 200 {
		panic("error occurred when downloading the file")
	}
	file, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Printf("Error occured when closing the file: %v", err)
		}
	}(file)
	_, err = io.Copy(file, response.Body)
	if err != nil {
		panic(err)
	}
}

func SearchUser(KeyWord string, page int) ([]JsonSearchUser, int) {
	if KeyWord == "" {
		panic("KeyWord is empty")
	}
	err := error(nil)
	if Cookies == "" {
		Cookies = GetCookies()
		if err != nil {
			panic(err)
		}
	}
	getUrl := "https://api.bilibili.com/x/web-interface/search/type"
	data := url.Values{}
	data.Set("search_type", "bili_user")
	data.Set("keyword", KeyWord)
	data.Set("page", strconv.Itoa(page))
	s, _ := getReq(data, getUrl, Cookies)
	if err != nil {
		panic(err)
	}
	var jSR JsonSearchRes
	err = json.Unmarshal(s, &jSR)
	if err != nil {
		panic(err)
	}
	return jSR.Data.Result, jSR.Data.NumPages
}

func GetCookies() string {
	getUrl := "https://www.bilibili.com/"
	data := url.Values{}
	data.Set("spm_id_from", "333.999.0.0")
	_, cookies := getReq(data, getUrl, "")
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
	s, _ := getReq(data, getUrl, Cookies)
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

func DownloadFace(FaceUrl string, uid int64) string {
	if FaceUrl == "" {
		panic("FaceUrl is empty")
	}
	if !pathExists("data/face") {
		err := os.Mkdir("data/face", 0777)
		if err != nil {
			return ""
		}
	}
	fileName := "data/face/" + strconv.FormatInt(uid, 10) + ".jpg"
	downloadFile(FaceUrl, fileName)
	return fileName
}

func DownloadVupJson() {
	getUrl := "https://cfapi.vtbs.moe/v1/short"
	fileName := "data/vup.json"
	downloadFile(getUrl, fileName)
}
