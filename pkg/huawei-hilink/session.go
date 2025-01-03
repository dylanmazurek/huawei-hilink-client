package huaweihilink

type Session struct {
	Host      string
	Username  string
	LoggedIn  bool
	SessionId string

	Token  string
	Token2 string

	PublicKey PublicKey
}

type PublicKey struct {
	rsan string
	rsae string
}
