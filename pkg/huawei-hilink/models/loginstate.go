package models

import "encoding/xml"

type LoginStateResponse struct {
	XMLName               xml.Name `xml:"response"`
	PasswordType          int      `xml:"password_type"`
	ExternPasswordType    int      `xml:"extern_password_type"`
	HistoryLoginFlag      int      `xml:"history_login_flag"`
	State                 int      `xml:"State"`
	LockStatus            int      `xml:"lockstatus"`
	AccountsNumber        int      `xml:"accounts_number"`
	RSAPadingType         int      `xml:"rsapadingtype"`
	RemainWaitTime        int      `xml:"remainwaittime"`
	WifiPWDSameWithWebPWD int      `xml:"wifipwdsamewithwebpwd"`
	Username              string   `xml:"username"`
	FirstLogin            int      `xml:"firstlogin"`
	UserLevel             string   `xml:"userlevel"`
}
