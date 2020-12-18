package initdo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
)

var (
	GlobalRepeatSubmitToken string
	Trainlocation           string
	Keycheckischange        string
	LeftTic                 string
)

type Passenger struct {
	Data `json:"data"`
}
type Data struct {
	NormalPassengers `json:"normal_passengers"`
}
type NormalPassengers []struct {
	PassengerName string `json:"passenger_name"`
	PassengerIdNo string `json:"passenger_id_no"`
	MobileNo      string `json:"mobile_no"`
	AllEncStr     string `json:"allEncStr"`
}

//获取initDc信息
func DcInIt(cookieText string) {
	url := "https://kyfw.12306.cn/otn/confirmPassenger/initDc?_json_att"
	method := "GET"

	client := &http.Client{
	}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("Cookie", cookieText)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}
	//globalRepeatSubmitToken
	flysnowRegexp := regexp.MustCompile(`var globalRepeatSubmitToken = '(.*?)';`)
	params := flysnowRegexp.FindStringSubmatch(string(body))
	//ticketInfoForPassengerForm
	ticket := regexp.MustCompile(`var ticketInfoForPassengerForm=(.*?);`)
	from := ticket.FindStringSubmatch(string(body))
	GlobalRepeatSubmitToken = params[1]
	//train_location
	train := regexp.MustCompile(`'train_location':'(.*?)'`)
	trainfrom := train.FindStringSubmatch(from[1])
	Trainlocation = trainfrom[1]
	//key_check_isChange
	isChange := regexp.MustCompile(`'key_check_isChange':'(.*?)'`)
	isChangefrom := isChange.FindStringSubmatch(from[1])
	Keycheckischange = isChangefrom[1]
	//leftTicketStr
	leftTicket := regexp.MustCompile(`'leftTicketStr':'(.*?)'`)
	leftTicketfrom := leftTicket.FindStringSubmatch(from[1])
	LeftTic = leftTicketfrom[1]
}

func GetPassengerDTO(cookieText string) Passenger {
	var d Passenger
	url := "https://kyfw.12306.cn/otn/confirmPassenger/getPassengerDTOs"
	method := "GET"

	client := &http.Client{
	}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return d
	}
	req.Header.Add("cookie", cookieText)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return d
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return d
	}
	err = json.Unmarshal([]byte(string(body)), &d)
	return d
}
