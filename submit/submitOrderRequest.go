package submit

import (
	"12306/city"
	"12306/config"
	"12306/initdo"
	"12306/model"
	"12306/notice"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	us "net/url"
	"os"
	"strings"
	"time"
)

//预定车票
func SubmitOrderRequest(cookieText string, secretStr string, traindate string, from string, to string) {
	now := time.Now()
	url := "https://kyfw.12306.cn/otn/leftTicket/submitOrderRequest"
	method := "POST"
	//enEscapeUrl, _ := u.QueryUnescape(secretStr)
	backtraindate := now.Format("2006-01-02")
	payload := strings.NewReader("secretStr=" + secretStr + "&train_date=" + traindate + "&back_train_date=" + backtraindate + "&tour_flag=dc&purpose_codes=ADULT&query_from_station_name=" + from + "&query_to_station_name=" + to)

	client := &http.Client{
	}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Cookie", cookieText)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

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
	fmt.Println(string(body))
}
func CheckOrderInfo(cookieText string, zw string, ccr []string) {
	d := initdo.GetPassengerDTO(cookieText)
	url := "https://kyfw.12306.cn/otn/confirmPassenger/checkOrderInfo"
	method := "POST"
	u := ""
	o := ""
	s := city.Seat(zw)
	for _, v := range d.Data.NormalPassengers {
		for _, r := range ccr {
			if v.PassengerName == r {
				u += "" + s + ",0,1," + v.PassengerName + ",1," + v.PassengerIdNo + "," + v.MobileNo + ",N," + v.AllEncStr + "_"
				o += "" + v.PassengerName + ",1," + v.PassengerIdNo + ",1" + "_"
			}
		}
	}
	fmt.Println(u)
	fmt.Println(o)
	payload := strings.NewReader("cancel_flag=2&bed_level_order_num=000000000000000000000000000000&passengerTicketStr=" + u + "&oldPassengerStr=" + o + "&tour_flag=dc&randCode=&whatsSelect=1&sessionId=&sig=&scene=nc_login&_json_att=&REPEAT_SUBMIT_TOKEN=" + initdo.GlobalRepeatSubmitToken)

	client := &http.Client{
	}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Cookie", cookieText)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

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
	fmt.Println("提交:", string(body))
	time.Sleep(time.Millisecond * 2000)
	ConfirmSingleForQueue(cookieText, u, o)
}

//提交订单
func ConfirmSingleForQueue(cookieText string, u string, o string) {
	var c config.Config
	c.GetConf()
	t := time.Now()
	var rq model.ReturnValue
	url := "https://kyfw.12306.cn/otn/confirmPassenger/confirmSingleForQueue"
	method := "POST"
	initdo.LeftTic = us.QueryEscape(initdo.LeftTic)
	reqs := "passengerTicketStr=" + u + "&oldPassengerStr=" + o + "&purpose_codes=00" + "&key_check_isChange=" + initdo.Keycheckischange + "&leftTicketStr=" + initdo.LeftTic + "&train_location=" + initdo.Trainlocation + "&seatDetailType=000&whatsSelect=1&roomType=00&dwAll=N&REPEAT_SUBMIT_TOKEN=" + initdo.GlobalRepeatSubmitToken

	payload := strings.NewReader(reqs)

	client := &http.Client{
	}
	req, err := http.NewRequest(method, url, payload)
	fmt.Println(payload)
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Cookie", cookieText)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	err = json.Unmarshal([]byte(string(body)), &rq)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("抢票结果:", string(body))
	if rq.Status {
		fmt.Println("抢到了快去付款吧！")

		_ = notice.SendMail(c.Mail, "恭喜抢票成功！", "抢票成功快登录12306付款30分钟内有效")

		fmt.Println("以邮箱通知！")
		os.Exit(2)
	}
	elapsed := time.Since(t)
	fmt.Println("耗时:", elapsed)
}
