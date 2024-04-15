package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/charmbracelet/huh"
	"io"
	"log"
	"math/rand/v2"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"now/sign/chatlist"
	"now/sign/userdata"
	"strings"
	"time"
)

var missionIDs = []int{2, 4, 6, 9, 10}
var client http.Client
var missionDailyMap = map[string]string{
	"cMissionDaily1":  "發送一則訊息",
	"cMissionDaily2":  "查看頻道列表",
	"cMissionDaily3":  "加入任一頻道",
	"cMissionDaily4":  "使用5分鐘",
	"cMissionDaily5":  "使用60分鐘",
	"cMissionDaily6":  "點擊巴哈相關連結",
	"cMissionDaily7":  "午間登入",
	"cMissionDaily8":  "晚間登入",
	"cMissionDaily9":  "查看別人小卡",
	"cMissionDaily10": "每日充能",
}

func init() {

}

type ErrorPayload struct {
	Error *struct {
		Code    int           `json:"code"`
		Message string        `json:"message"`
		Status  string        `json:"status"`
		Details []interface{} `json:"details"`
	} `json:"error,omitempty"`
}

func main() {
	jar, _ := cookiejar.New(nil)
	client = http.Client{Jar: jar}

	token := ""
	var resultPayload []byte
	var userData userdata.UserData
	var username, password, otpcode string
	err := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Username").
				Value(&username),
			huh.NewInput().
				Title("Password").
				Value(&password),
			huh.NewInput().
				Title("OTP Code").
				Value(&otpcode),
		)).Run()
	if err != nil {
		return
	}
	log.Println(err)

	token = getCsrfToken()
	param := url.Values{}
	param.Add("uid", username)
	param.Add("passwd", password)
	param.Add("bahamutCsrfToken", token)
	if otpcode != "" {
		param.Add("twoStepAuth", otpcode)
	}

	req, _ := http.NewRequest("POST", "https://api.gamer.com.tw/mobile_app/user/v4/do_login.php", strings.NewReader(param.Encode()))

	req.Header.Add("User-Agent", "Bahadroid (https://www.gamer.com.tw/)")
	req.Header.Add("X-Bahamut-App-Android", "tw.com.gamer.now")
	req.Header.Add("X-Bahamut-App-Version", "230")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Cookie", "ckBahamutCsrfToken="+token)

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	err = json.NewDecoder(resp.Body).Decode(&userData)
	if err != nil {
		log.Println(err)
		return
	}

	var result *http.Response
	result, _ = client.Get("https://api.gamer.com.tw/lite/v1/chat_list.php")
	resultPayload, _ = io.ReadAll(result.Body)
	var errorPayload ErrorPayload
	json.Unmarshal(resultPayload, &errorPayload)
	if errorPayload.Error != nil {
		return
	}
	var chatList chatlist.ChatList
	err = json.Unmarshal(resultPayload, &chatList)
	if err != nil {
		log.Println(err)
	}
	log.Printf("使用者資料: %s(%s)", userData.Data.Nickname, userData.Data.Userid)
	log.Printf("準備執行釘選的頻道簽到(前十個):\n%s", strings.Join(chatList.Pin, "\n"))
	for i, s := range chatList.Pin {
		for _, missionID := range missionIDs {
			missionStr := missionDailyMap[fmt.Sprintf("cMissionDaily%d", missionID)]
			log.Printf("執行第%d個頻道的第%d個任務 - %s", i+1, missionID, missionStr)
			formData := url.Values{}

			formData.Set("bahamutCsrfToken", token)
			unixSec := time.Now().Unix()

			hash := sha256.New()
			hash.Write([]byte(fmt.Sprintf("%s_Z8jDeK3Y6S_cMissionDaily%d_%d_%s", s, missionID, unixSec, userData.Data.Userid)))
			hashedData := hash.Sum(nil)
			formData.Add("checkHash", hex.EncodeToString(hashedData))
			formData.Add("jid", s)
			formData.Add("ts", fmt.Sprint(unixSec))
			formData.Add("type", fmt.Sprintf("cMissionDaily%d", missionID))

			request, err := http.NewRequest("POST", fmt.Sprintf("https://api.gamer.com.tw/lite/v2/conference_detall.php?jid=%s&type=cMissionDaily%d", s, missionID), strings.NewReader(formData.Encode()))
			if err != nil {
				log.Println(err)
			}
			request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
			request.Header.Add("User-Agent", "BahaLite/1.4.1 (tw.com.gamer.now; build: 49; Android 8.8.8) okHttp/4.4.0")

			request.Header.Add("Cookie", "ckBahamutCsrfToken="+token)
			resp, err := client.Do(request)
			if err != nil {
				log.Println(err)
			}
			body, _ := io.ReadAll(resp.Body)
			var errorPayload ErrorPayload
			json.Unmarshal(body, &errorPayload)

			if errorPayload.Error == nil {
				log.Println("執行成功")
				log.Println(string(body))
			} else {
				log.Println("執行失敗")
				log.Println(string(body))
			}

			time.Sleep(1 * time.Second)
		}

		if i == 10 {
			log.Println("超過十個，僅執行前面十個，程式結束。")
			break
		}
	}
}

func getCsrfToken() string {
	str := ""
	for i := 0; i < 8; i++ {
		hStr := hex.EncodeToString([]byte{byte(rand.IntN(256))})
		if len(hStr) == 1 {
			hStr = "0" + hStr
		}
		str += hStr
	}
	return str
}
