package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/dustinxie/ecc"
	"github.com/google/uuid"
	"io"
	"log"
	"math/big"
	"net/http"
)

var (
	UserID     = "42bb"
	secpAppID  = "5dde4e1bdf9e4966b387ba58f4b3fdc3"
	deviceID   = ""
	signatures = ""
	//privateKey *ecdsa.PrivateKey
	pubKey = ""
)

func NewPrivateKey() (*ecdsa.PrivateKey, error) {
	p256k1 := ecc.P256k1()
	return ecdsa.GenerateKey(p256k1, rand.Reader)
}

func NewPrivateKeyFromHex(hex_ string) (*ecdsa.PrivateKey, error) {
	data, err := hex.DecodeString(hex_)
	if err != nil {
		return nil, err
	}
	return NewPrivateKeyFromBytes(data), nil

}

func NewPrivateKeyFromBytes(priv []byte) *ecdsa.PrivateKey {
	p256k1 := ecc.P256k1()
	x, y := p256k1.ScalarBaseMult(priv)
	return &ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{
			Curve: p256k1,
			X:     x,
			Y:     y,
		},
		D: new(big.Int).SetBytes(priv),
	}
}

func PrivateKeyToHex(private *ecdsa.PrivateKey) string {
	return hex.EncodeToString(PrivateKeyToBytes(private))
}

func PrivateKeyToBytes(private *ecdsa.PrivateKey) []byte {
	return private.D.Bytes()
}

func PublicKeyToHex(public *ecdsa.PublicKey) string {
	return hex.EncodeToString(PublicKeyToBytes(public))
}

func PublicKeyToBytes(public *ecdsa.PublicKey) []byte {
	x := public.X.Bytes()
	if len(x) < 32 {
		for i := 0; i < 32-len(x); i++ {
			x = append([]byte{0}, x...)
		}
	}

	y := public.Y.Bytes()
	if len(y) < 32 {
		for i := 0; i < 32-len(y); i++ {
			y = append([]byte{0}, y...)
		}
	}
	return append(x, y...)
}
func GetSHA256Encode(data string) string {
	h := sha256.New()
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}
func sign() {
	deviceID = GetSHA256Encode(UserID)
	privateKey, err := NewPrivateKeyFromHex(deviceID)
	if err != nil {
		log.Println(err, "err")
	}
	singdata := fmt.Sprintf("%s:%s:%s:%d", secpAppID, deviceID, UserID, 1)
	hash := sha256.Sum256([]byte(singdata))
	data, _ := ecc.SignBytes(privateKey, hash[:], ecc.RecID|ecc.LowerS)
	signatures = hex.EncodeToString(data) //strconv.Itoa(state.nonce)
	pubKey = PublicKeyToHex(&privateKey.PublicKey)
	log.Println(deviceID, "deviceID")
	log.Println(pubKey, "pubKey")
	log.Println(signatures, "signatures")
}
func alistSession() {
	api := "https://api.aliyundrive.com/users/v1/users/device/create_session"
	form := fmt.Sprintf(`{"deviceName": "samsung","modelName": "SM-G9810","nonce":0,"pubKey":"%s","refreshToken":"16d049976"}`, pubKey)
	req, err := http.NewRequest("POST", api, bytes.NewBufferString(form))
	if err != nil {
		log.Println(err, "NewRequest")
	}
	req.Header.Set("content-type", "application/json; charset=utf-8")
	req.Header.Set("authorization", "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpQzM2VjZFpXazhxZ1wJVkUuQUxMXCIsXCJGSUxFLkFMTFwiLFwiVklFVy5BTExcIixcIlNIQVJFLkFMTFwiLFwiU1RPUkFHRS5BTExcIixcIlNUT1JBR0VGSUxFLkxJU1RcIixcIlVTRVIuQUxMXCIsXCJCQVRDSFwiLFwiQUNDT1VOVC5BTExcIixcIklNQUdFLkFMTFwiLFwiSU5WSVRFLkFMTFwiLFwiU1lOQ01BUFBJTkcuTElTVFwiXSxcInJvbGVcIjpcInVzZXJcIixcInJlZlwiOlwiXCIsXCJkZXZpY2VfaWRcIjpcIjE2ZDA5MzRlNzYwZDRiOTZhOTMxNjQ0ZDFmM2M5OTc2XCJ9IiwiZXhwIjoxNjc5MjAyNTEzLCJpYXQiOjE2NzkxOTUyNTN9.pr5Prp1DlbrjBZgoT-INF65Bhcmxj3-8pbdRH-CGaOLjpeHL0x9E1cgn5O_8yHsV6q7GQlsYBQcnz1OksDakBfeJ7Tbai5u_2xupvotOYywjmQw2RAZ-Ib57roW9JfUziHCByeXfHOj3ubtl3oNBI9-VES2AtUyGwaoxjFKW9WU")
	req.Header.Add("user-agent", "AliApp(AYSD/4.3.0) com.alicloud.databox/28792670 Channel/36178427979800@rimet_android_4.3.0 language/zh-CN /Android Mobile/samsung SM-G9810")
	//req.Header.Add("origin", "https://aliyundrive.com")
	req.Header.Add("referer", "https://aliyundrive.com/")
	req.Header.Add("accept", "application/json")
	//req.Header.Add("accept-language", "zh-CN,zh;q=0.8,en-US;q=0.5,en;q=0.3")
	req.Header.Add("x-canary", "client=Android,app=adrive,version=v4.3.0")
	req.Header.Add("x-request-id", uuid.NewString())
	req.Header.Add("x-device-id", deviceID)   //
	req.Header.Add("x-signature", signatures) //
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err, "Do")

	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err, "ReadAll")

	}
	log.Println(string(body), "body")
}
func AlistDownload(fileId string) string {
	api := "https://api.aliyundrive.com/v2/file/get_download_url"
	form := fmt.Sprintf(`{"expire_sec": 14400, "drive_id": "1401739120","file_id": "6413d0c48eaae5f3ef754b9f8459870c5463ac97"}`)
	req, err := http.NewRequest("POST", api, bytes.NewBufferString(form))
	if err != nil {
		//return "" err
		log.Println(err, "NewRequest")
	}
	req.Header.Set("content-type", "application/json; charset=utf-8")
	req.Header.Set("authorization", "eyJhbbUpzb24iOiJ7XCJjbGllbnRJZFwiOlwiMjlcIixcInzYwZDLCJpdRH-CGaOLjakBfeJ7Tbai5u_2xupvotOYywjmQw2RAZ-Ib57roW9JfUziHCByeXfHOj3ubtl3oNBI9-VES2AtUyGwaoxjFKW9WU")
	req.Header.Add("user-agent", "AliApp(AYSD/4.3.0) com.alicloud.databox/28792670 Channel/36178427979800@rimet_android_4.3.0 language/zh-CN /Android Mobile/samsung SM-G9810")
	//req.Header.Add("origin", "https://aliyundrive.com")
	req.Header.Add("referer", "https://aliyundrive.com/")
	req.Header.Add("accept", "application/json")
	//req.Header.Add("accept-language", "zh-CN,zh;q=0.8,en-US;q=0.5,en;q=0.3")
	req.Header.Add("x-canary", "client=Android,app=adrive,version=v4.3.0")
	req.Header.Add("x-request-id", uuid.NewString())
	req.Header.Add("x-device-id", "") // 400
	req.Header.Add("x-signature", "")
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
