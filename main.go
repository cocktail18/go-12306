package main

import (
	"12306/city"
	"12306/config"
	"12306/initdo"
	"12306/model"
	"12306/submit"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"
)

var wg sync.WaitGroup
var once sync.Once
var loginis bool

type resYuPiao struct {
	Data `json:"data"`
}
type Data struct {
	Result []string `json:"result"`
	Flag   bool     `json:"flag"`
}
type CheckUser struct {
	Data `json:"data"`
}

func main() {
	//初始化配置文件
	var c config.Config
	c.GetConf()
	//fmt.Println(c.Train)
	wg.Add(2)
	go GetHuiJia(c)
	go IsLogin(c)
	wg.Wait()
}

func GetHuiJia(c config.Config) {
	defer wg.Done()
	for {
		var yp resYuPiao
		//var yd *model.OrderRequest

		yd := &model.OrderRequest{
			Traindate: c.Time,
			From:      city.GetCity(c.Form),
			To:        city.GetCity(c.To),
		}
		url := "https://kyfw.12306.cn/otn/leftTicket/query?leftTicketDTO.train_date=" + yd.Traindate + "&leftTicketDTO.from_station=" + yd.From + "&leftTicketDTO.to_station=" + yd.To + "&purpose_codes=ADULT"
		method := "GET"
		client := &http.Client{}
		req, err := http.NewRequest(method, url, nil)

		if err != nil {
			fmt.Println(err)
			return
		}
		req.Header.Add("cookie", c.Cookie)
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
		err = json.Unmarshal([]byte(string(body)), &yp)
		for _, v := range yp.Data.Result {
			arr := strings.Split(v, "|")
			name := [...]string{
				"车次", "出发车站", "到达车站", "出发时间", "到达时间", "历时", "商务座", "一等座", "二等座", "高级软卧", "软卧", "动卧", "硬卧", "软座", "硬座", "无座", "其他", "备注",
			}

			data := map[string]string{
				"车次":   "",
				"出发车站": "",
				"到达车站": "",
				"出发时间": "",
				"到达时间": "",
				"历时":   "",
				"商务座":  "",
				"一等座":  "",
				"二等座":  "",
				"高级软卧": "",
				"软卧":   "",
				"动卧":   "",
				"硬卧":   "",
				"软座":   "",
				"硬座":   "",
				"无座":   "",
				"其他":   "",
				"备注":   "",
			}
			data["车次"] = arr[3]                 // 获取车次信息，在3号位置
			data["出发车站"] = city.GetCity(arr[6]) // 始发站信息在6号位置
			data["到达车站"] = city.GetCity(arr[7]) // 终点站信息在7号位置
			data["出发时间"] = arr[8]               // 出发时间在8号位置
			data["到达时间"] = arr[9]               // 抵达时间在9号位置
			data["历时"] = arr[10]                // 经历时间在10号位置
			data["商务座"] = arr[32]               // 特别注意，商务座在32或25位置
			data["一等座"] = arr[31]               // 一等座信息在31号位置
			data["二等座"] = arr[30]               // 二等座信息在30号位置
			data["高级软卧"] = arr[21]              // 高级软卧信息在21号位置
			data["软卧"] = arr[23]                // 软卧信息在23号位置
			data["动卧"] = arr[27]                // 动卧信息在27号位置
			data["硬卧"] = arr[28]                // 硬卧信息在28号位置
			data["软座"] = arr[24]                // 软座信息在24号位置
			data["硬座"] = arr[29]                // 硬座信息在29号位置
			data["无座"] = arr[26]                // 无座信息在26号位置
			data["其他"] = arr[22]                // 其他信息在22号位置
			data["备注"] = arr[1]                 // 备注信息在1号位置
			for _, v := range name {
				if data[v] == "" {
					data[v] = "-"
				}
			}

			if len(c.Train) == 0 {
				fmt.Println(fmt.Sprintf("车次:%s,出发车站:%s,到达车站:%s,出发时间:%s,到达时间:%s,商务座:%s,一等座:%s,二等座:%s,高级软卧:%s,软卧:%s,动卧:%s,硬卧:%s,软座:%s,硬座:%s,无座:%s,备注:%s",
					data["车次"], data["出发车站"], data["到达车站"], data["出发时间"], data["到达时间"], data["商务座"], data["一等座"], data["二等座"], data["高级软卧"], data["软卧"], data["动卧"],
					data["硬卧"], data["软座"], data["硬座"], data["无座"], data["备注"]))
				time.Sleep(time.Millisecond * 500)
				for _, v := range c.Seats {
					if data[v] != "无" && data[v] != "-" {
						fmt.Println(fmt.Sprintf("%s", "有余票了！"))
						fmt.Println(fmt.Sprintf("%s", "去提交订单了！"))
						fmt.Println(fmt.Sprintf("车次：%s,ID：%s", data["车次"], arr[0]))
						if !loginis {
							fmt.Println(fmt.Sprintf("%s", "cookie失效了，请重新登录！"))
						} else {
							submit.SubmitOrderRequest(c.Cookie, arr[0], yd.Traindate, c.Form, c.To)

							initdo.DcInIt(c.Cookie)
							time.Sleep(time.Millisecond)

							submit.CheckOrderInfo(c.Cookie, v, c.Passengers)
						}

					}
				}
			} else {
				for _, v := range c.Train {
					if v == data["车次"] {
						fmt.Println(fmt.Sprintf("车次:%s,出发车站:%s,到达车站:%s,出发时间:%s,到达时间:%s,商务座:%s,一等座:%s,二等座:%s,高级软卧:%s,软卧:%s,动卧:%s,硬卧:%s,软座:%s,硬座:%s,无座:%s,备注:%s",
							data["车次"], data["出发车站"], data["到达车站"], data["出发时间"], data["到达时间"], data["商务座"], data["一等座"], data["二等座"], data["高级软卧"], data["软卧"], data["动卧"],
							data["硬卧"], data["软座"], data["硬座"], data["无座"], data["备注"]))
						time.Sleep(time.Millisecond * 500)
						for _, v := range c.Seats {
							if data[v] != "无" && data[v] != "-" {
								fmt.Println(fmt.Sprintf("%s", "有余票了！"))
								fmt.Println(fmt.Sprintf("%s", "去提交订单了！"))
								fmt.Println(fmt.Sprintf("车次：%s,ID：%s", data["车次"], arr[0]))

								if !loginis {
									fmt.Println(fmt.Sprintf("%s", "cookie失效了，请重新登录！"))
								} else {
									submit.SubmitOrderRequest(c.Cookie, arr[0], yd.Traindate, c.Form, c.To)

									initdo.DcInIt(c.Cookie)
									time.Sleep(time.Millisecond)

									submit.CheckOrderInfo(c.Cookie, v, c.Passengers)

								}

							}
						}
					}
				}
			}
		}
	}

}

func IsLogin(c config.Config) {
	defer wg.Done()
	for {
		var login CheckUser
		url := "https://kyfw.12306.cn/otn/login/checkUser"
		method := "GET"

		client := &http.Client{
		}
		req, err := http.NewRequest(method, url, nil)

		if err != nil {
			fmt.Println(err)
			return
		}
		req.Header.Add("Cookie", c.Cookie)

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
		err = json.Unmarshal([]byte(string(body)), &login)
		fmt.Println(login.Data.Flag)
		loginis = login.Data.Flag
		//每一分钟去查登录是否失效
		time.Sleep(time.Minute * 1)

	}

}
