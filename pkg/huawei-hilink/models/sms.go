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
