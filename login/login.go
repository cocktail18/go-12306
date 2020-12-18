package login

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var callback = "jQuery19106558532255409184_" + s
var s = strconv.FormatInt(time.Now().UnixNano()/1e6, 10)

type Captcha struct {
	Image         string `json:"image"`
	ResultMessage string `json:"result_message"`
}

func Logins() string {
	fmt.Println("请输入cookie")
	cookie := bufio.NewScanner(os.Stdin)
	cookie.Scan()
	cookieText := cookie.Text()
	//maps := [...]string{
	//	"37,45",
	//	"111,46",
	//	"188,48",
	//	"259,47",
	//	"43,117",
	//	"111,119",
	//	"191,119",
	//	"264,117",
	//}
	//GetCaptcha(maps)
	return cookieText
}

func GetCaptcha(maps [8]string) {
	var cp Captcha

	url := "https://kyfw.12306.cn/passport/captcha/captcha-image64?login_site=E&module=login&rand=sjrand&1607992773484&callback=" + callback + "&_=" + s + ""
	method := "GET"

	client := &http.Client{
	}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Cookie", "_passport_session=7f8afc6afda243d2bf126ed2f0a5ceba2157; _passport_ct=5fd87b08c1a742f48627336a31a3de4ct4529; route=c5c62a339e7744272a54643b3be5bf64; BIGipServerotn=686817802.24610.0000; BIGipServerpool_passport=98828810.50215.0000")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	arr := strings.Split(string(body), "/**/"+callback+"(")
	arrs := arr[1][:len(arr[1])-2]
	err = json.Unmarshal([]byte(string(arrs)), &cp)
	b, err := base64.StdEncoding.DecodeString(cp.Image)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile("./yazm.jpg", []byte(b), 0666)
	if err != nil {
		panic(err)
	}
	jg := YanZhengMaShiBie()
	flysnowRegexp := regexp.MustCompile(`<B>(.*?)</B>`)
	params := flysnowRegexp.FindStringSubmatch(jg)

	ress := []string{}
	imgarr := []string{}
	for _, param := range params {
		imgarr = strings.Split(param, " ")
	}
	for _, v := range imgarr {
		ress = append(ress, v)
	}
	answer := ""
	for _, v := range ress {

		int, _ := strconv.Atoi(v)
		arrs := strings.Split(maps[int-1], ",")
		fmt.Println(arrs)
		z, _ := strconv.Atoi(arrs[0])
		y, _ := strconv.Atoi(arrs[1])
		fmt.Println(z, ",", y)
		time.Sleep(time.Millisecond * 2000)
		answer += arrs[0] + "," + arrs[1] + ","
	}
	fmt.Println(answer[:len(answer)-1])
	resq := Check(answer[:len(answer)-1])
	fmt.Println(resq)
}
func YanZhengMaShiBie() string {
	url := "http://littlebigluo.qicp.net:47720/"
	method := "POST"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	file, errFile1 := os.Open("./yazm.jpg")
	defer file.Close()
	part1, errFile1 := writer.CreateFormFile("pic_xxfile", filepath.Base("./yazm.jpg"))
	_, errFile1 = io.Copy(part1, file)
	if errFile1 != nil {

		fmt.Println(errFile1)
	}
	err := writer.Close()
	if err != nil {
		fmt.Println(err)
	}

	client := &http.Client{
	}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	return string(body)
}

func Check(answer string) string {
	url := "https://kyfw.12306.cn/passport/captcha/captcha-check?callback=" + callback + "&answer=" + answer + "&rand=sjrand&login_site=E&_=" + s + ""
	method := "GET"

	client := &http.Client{
	}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return ""
	}
	req.Header.Add("Cookie", "_passport_session=7f8afc6afda243d2bf126ed2f0a5ceba2157; _passport_ct=5fd87b08c1a742f48627336a31a3de4ct4529; route=c5c62a339e7744272a54643b3be5bf64; BIGipServerotn=686817802.24610.0000; BIGipServerpool_passport=98828810.50215.0000")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(body)
}
