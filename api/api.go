package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

var Cookies string

func getReq(data url.Values, getUrl string, cookies string) ([]byte, string, error) {
	u, err := url.ParseRequestURI(getUrl)
	if err != nil {
		log.Printf("ParseRequestURI error: %v", err)
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
		log.Printf("Error Occured when sending GET request: %v", err)
		return nil, "", err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			log.Printf("Error orrured when closing the session: %v", err)
		}
	}(resp.Body)
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
		log.Printf("Error occured when reading the response body: %v", err)
		return nil, "", err
	}
	return s, cookies, nil
}

func SearchUser(KeyWord string, page int) ([]JsonSearchUser, int) {
	if KeyWord == "" {
		return nil, 0
	}
	if Cookies == "" {
		Cookies = GetCookies()
	}
	getUrl := "https://api.bilibili.com/x/web-interface/search/type"
	data := url.Values{}
	data.Set("search_type", "bili_user")
	data.Set("keyword", KeyWord)
	data.Set("page", strconv.Itoa(page))
	s, _, err := getReq(data, getUrl, Cookies)
	if err != nil {
		return nil, 0
	}
	var jSR JsonSearchRes
	err = json.Unmarshal(s, &jSR)
	if err != nil {
		log.Printf("Error occured when unmarshal the json: %v", err)
		return nil, 0
	}
	return jSR.Data.Result, jSR.Data.NumPages
}

func GetCookies() string {
	getUrl := "https://www.bilibili.com/"
	data := url.Values{}
	data.Set("spm_id_from", "333.999.0.0")
	_, cookies, err := getReq(data, getUrl, "")
	if err != nil {
		return ""
	}
	return cookies
}

func GetUserInfo(uid int64) JsonUserInfo {
	if Cookies == "" {
		Cookies = GetCookies()
	}
	getUrl := "https://account.bilibili.com/api/member/getCardByMid"
	data := url.Values{}
	data.Set("mid", strconv.FormatInt(uid, 10))
	s, _, err := getReq(data, getUrl, Cookies)
	if err != nil {
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
