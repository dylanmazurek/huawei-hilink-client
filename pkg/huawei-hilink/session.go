package huaweihilink

type Session struct {
	Host      string `json:"host"`
	Username  string `json:"username"`
	LoggedIn  bool   `json:"-"`
	SessionId string `json:"session_id"`

	Token  string `json:"token"`
	Token2 string `json:"token2"`

	PublicKey PublicKey `json:"public_key"`
}

type PublicKey struct {
	Rsan string `json:"rsan"`
	Rsae string `json:"rsae"`
}
