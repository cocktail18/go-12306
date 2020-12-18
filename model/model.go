package model

//预定
type OrderRequest struct {
	SecretStr     string `json:"secretStr"`     //车次ID
	Backtraindate string `json:"backtraindate"` //查票时间
	Traindate     string `json:"traindate"`     // 出发时间
	From          string `json:"from"`
	To            string `json:"to"`
}

type ReturnValue struct {
	Status     bool `json:"status"`
	Httpstatus int  `json:"httpstatus"`
	Data       struct {} `json:"data"`
}
