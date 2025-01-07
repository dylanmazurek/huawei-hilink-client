package models

import "encoding/xml"

type SMSCountResp struct {
	XMLName      xml.Name `xml:"response"`
	LocalUnread  int      `xml:"LocalUnread"`
	LocalInbox   int      `xml:"LocalInbox"`
	LocalOutbox  int      `xml:"LocalOutbox"`
	LocalDraft   int      `xml:"LocalDraft"`
	LocalDeleted int      `xml:"LocalDeleted"`
	SimUnread    int      `xml:"SimUnread"`
	SimInbox     int      `xml:"SimInbox"`
	SimOutbox    int      `xml:"SimOutbox"`
	SimDraft     int      `xml:"SimDraft"`
	LocalMax     int      `xml:"LocalMax"`
	SimMax       int      `xml:"SimMax"`
	SimUsed      int      `xml:"SimUsed"`
	NewMsg       int      `xml:"NewMsg"`
}

type SMSListReq struct {
	XMLName   xml.Name `xml:"request"`
	Phone     string   `xml:"phone"`
	PageIndex int      `xml:"pageindex"`
	ReadCount int      `xml:"readcount"`
	Nonce     string   `xml:"nonce"`
}

type SMSListResp struct {
	XMLName xml.Name `xml:"response"`
	Pwd     string   `xml:"pwd"`
	Hash    string   `xml:"hash"`
	Iter    int      `xml:"iter"`
}

// <?xml version="1.0" encoding="UTF-8"?><request>
// <Index>-1</Index><Phones><Phone>${(phones)}</Phone></Phones><Sca></Sca><Content>${message}</Content>
// <Length>${message.length}</Length><Reserved>1</Reserved>
// <Date>2021-10-27 00:12:24</Date><nonce>${encrpt}</nonce></request>
type SMSSendReq struct {
	XMLName  xml.Name `xml:"request"`
	Index    int      `xml:"Index"`
	Phones   []Phone  `xml:"Phones"`
	Sca      any      `xml:"Sca"`
	Content  string   `xml:"Content"`
	Length   int      `xml:"Length"`
	Reserved int      `xml:"Reserved"`
	Date     string   `xml:"Date"`
	Nonce    string   `xml:"nonce"`
}

type Phone struct {
	Phone string `xml:"Phone"`
}
