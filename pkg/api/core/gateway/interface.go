package gateway

type Result struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
}

type ResultError struct {
	Status int   `json:"status"`
	Error  error `json:"error"`
}

type Info struct {
	GroupID    uint   `json:"group_id"`   //コントローラへの通知先で判別するためのGroupID ID=0の場合は管理者
	UUID       string `json:"uuid"`       //コントローラへの通知先で判別するためのUUID
	Controller string `json:"controller"` //コントローラへの通知先のIPアドレスとPort番号

}
