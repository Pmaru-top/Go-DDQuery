package api

import (
	"encoding/json"
	"errors"
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
	if !pathExists("data/face") {
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

func DownloadVupJson() {
	getUrl := "https://cfapi.vtbs.moe/v1/short"
	fileName := "data/vup.json"
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
