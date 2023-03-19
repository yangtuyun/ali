package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/tendermint/tendermint/crypto/secp256k1"
)

var (
	appId         = "5dde4e1bdf9e4966b387ba58f4b3fdc3"  //
	deviceId      = "hbOcHN13FDoBASQdaddsaJigAmjZRH"    //
	userId        = "42fbdaa34dcf971198cdddasa33ecb8bb" //
	nonce         = 0
	publicKey     = ""
	signatureData = ""
)

func randomString(l int) []byte {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		rand.NewSource(time.Now().UnixNano())
		bytes[i] = byte(randInt(1, 2^256-1))
	}
	return bytes
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func AddHeader(req *http.Request) *http.Request {
	req.Header.Set("authorization", "Bearer eyJhbGciOc3RvbUpzb24iOiJ7XCJjtYWluSWRcIjpcImJqMjlcIixcInNjb3BlXCI6W1wiRFJJVkUuQUxMXCIsXCJTSEFSRS5BTExcIixcIkZJTEUuQUxMXCIsXCJVU0VSLkFMTFwiLFwiVklFVy5BTExcIixcIlNUT1JBR0UuQUxMXCIsXCJTVE9SQUdFRklMRS5MSVNUXCIsXCJCQVRDSFwiLFwiT0FVVEguQUxMXCIsXCJJTUFHRS5BTExcIixcIklOVklURS5BTExcIixcIkFDQ09VTlQuQUxMXCIsXCJTWU5DTUFQUElORy5MSVNUXCIsXCJTWU5DTUFQUElORy5ERUxFVEVcIl0sXCJyb2xlXCI6XCJ1c2VyXCIsXCJyZWZcIjpcImh0dHBzOi8vd3d3LmFsaXl1bmRyaXZlLmNvbS9cIixcImRldmljZV9pZFwiOlwiYWZkOTY4MTYxZjE3NDg0OWE4OTkxNjJlNjAyZWJmN2JcIn0iLCJleHAiOjE2NzkxNjQzNTYsImlhdCI6MTY3OTE1NzA5Nn0.JKk5a6-IJOEyN8tpEhlei84_z8mIWvIVwqVtoGIo_lIiSRZcN9q1QUIgYkGFypwLl3kh3mKRIjoLoVWC04clL4W5jeQqn72txzbgQ8V7MF_kxyYJ27UchAjZUsX-AluGXnP2cD7kJrmlmcmdVdHGgX_YYx3BecyY0RiH0hh9UY8")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36")
	req.Header.Add("origin", "https://aliyundrive.com")
	req.Header.Add("accept", "*/*")
	req.Header.Add("Accept-Language", "zh-CN,zh;q=0.8,en-US;q=0.5,en;q=0.3")
	req.Header.Add("x-canary", "client=web,app=adrive,version=v3.17.0")
	req.Header.Add("x-device-id", deviceId)
	req.Header.Add("x-signature", signatureData)
	return req
}

func CreateSession() error {
	api := "https://api.aliyundrive.com/users/v1/users/device/create_session"
	form := fmt.Sprintf(`{"deviceName": "Edge浏览器","modelName": "Windows网页版","pubKey":"%s"}`, publicKey)
	req, err := http.NewRequest("POST", api, bytes.NewBufferString(form))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req = AddHeader(req)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	log.Println(string(body), "body")
	if strings.Contains(string(body), "true") {
		return nil
	}
	return errors.New(string(body))
}
func RenewSession() error {
	api := "https://api.aliyundrive.com/users/v1/users/device/renew_session"
	//form := fmt.Sprintf(`{"deviceName": "Edge浏览器","modelName": "Windows网页版","pubKey":"%s"}`, publicKey)
	req, err := http.NewRequest("POST", api, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req = AddHeader(req)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	log.Println(string(body), "body")
	if strings.Contains(string(body), "true") {
		return nil
	}
	return errors.New(string(body))
}

func AliDownload(fileId string) string {
	api := "https://api.aliyundrive.com/v2/file/get_download_url"
	form := fmt.Sprintf(`{"expire_sec": 14400, "drive_id": "1739120","file_id": "6413d0c48459870c5463ac97"}`)
	req, err := http.NewRequest("POST", api, bytes.NewBufferString(form))
	if err != nil {
		//return "" err
		log.Println(err, "NewRequest")
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req = AddHeader(req)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err, "Do")

	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err, "ReadAll")

	}
	log.Println(string(body))
	var data map[string]interface{}
	json.Unmarshal(body, &data)
	url, ok := data["url"].(string)
	if ok {
		return url
	}
	return url
}

func InitAliKey() error {
	max := 32
	key := randomString(max)
	data := fmt.Sprintf("%s:%s:%s:%d", appId, deviceId, userId, nonce)
	var privKey = secp256k1.PrivKey(key)
	pubKeys := privKey.PubKey()
	publicKey = "04" + hex.EncodeToString(pubKeys.Bytes())
	signature, err := privKey.Sign([]byte(data))
	if err != nil {
		return err
	}
	signatureData = hex.EncodeToString(signature) + "01"
	log.Println(publicKey, "publicKey")
	log.Println(signatureData, "signatureData")
	return nil
}
func cc() {
	err := InitAliKey()
	if err != nil {
		log.Println(err)
	}

	err = CreateSession()
	if err != nil {
		log.Println(err)
	}
	//err = RenewSession()
	//if err != nil {
	//	log.Println(err)
	//}
	AliDownload("")

}
