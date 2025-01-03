package models

import (
	"encoding/xml"
	"fmt"
	"regexp"
)

var (
	txPowerRegex = regexp.MustCompile(`PPusch:(?P<usch>\d+)dBm PPucch:(?P<ucch>\d+)dBm PSrs:(?P<rs>\d+)dBm PPrach:(?P<rach>\d+)dBm`)
	earfcnRegex  = regexp.MustCompile(`DL:(?P<dl>\d+) UL:(?P<ul>\d+)`)
)

// <?xml version=\"1.0\" encoding=\"UTF-8\"?>\r\n<response>\r\n<pci>266</pci>\r\n<sc></sc>\r\n<cell_id>21890627</cell_id>\r\n<rsrq>-14.0dB</rsrq>\r\n<rsrp>-105dBm</rsrp>\r\n<rssi>-71dBm</rssi>\r\n<sinr>-7dB</sinr>\r\n<rscp></rscp>\r\n<ecio></ecio>\r\n<mode>7</mode>\r\n<ulbandwidth>20MHz</ulbandwidth>\r\n<dlbandwidth>20MHz</dlbandwidth>\r\n<txpower>PPusch:21dBm PPucch:13dBm PSrs:23dBm PPrach:14dBm</txpower>\r\n<tdd></tdd>\r\n<ul_mcs>mcsUpCarrier1:23</ul_mcs>\r\n<dl_mcs>mcsDownCarrier1Code0:0 mcsDownCarrier1Code1:0</dl_mcs>\r\n<earfcn>DL:299 UL:18299</earfcn>\r\n<rrc_status></rrc_status>\r\n<rac></rac>\r\n<lac></lac>\r\n<tac>53281</tac>\r\n<band>1</band>\r\n<nei_cellid></nei_cellid>\r\n<plmn>50502</plmn>\r\n<ims>0</ims>\r\n<wdlfreq></wdlfreq>\r\n<lteulfreq>19499</lteulfreq>\r\n<ltedlfreq>21399</ltedlfreq>\r\n<transmode>TM[4]</transmode>\r\n<enodeb_id>0085510</enodeb_id>\r\n<cqi0>10</cqi0>\r\n<cqi1>127</cqi1>\r\n<ulfrequency>1949900kHz</ulfrequency>\r\n<dlfrequency>2139900kHz</dlfrequency>\r\n<arfcn></arfcn>\r\n<bsic></bsic>\r\n<rxlev></rxlev>\r\n</response>
type SignalInfoResp struct {
	Pci    int    `xml:"pci"`           // Physical Cell ID
	Sc     string `xml:"sc,omitempty"`  // Synchronization Code
	CellId int    `xml:"cell_id"`       // Cell ID
	Rscp   string `xml:"rscp"`          // Received Signal Code Power
	Ecio   string `xml:"ecio"`          //
	Mode   int    `xml:"mode"`          // Network Mode
	Tdd    string `xml:"tdd,omitempty"` // Time Division Duplex

	RrcStatus    string `xml:"rrc_status"`        // Radio Resource Control Status
	Rac          string `xml:"rac,omitempty"`     // Routing Area Code
	Lac          string `xml:"lac,omitempty"`     // Location Area Code
	Tac          int    `xml:"tac"`               // Tracking Area Code
	Band         int    `xml:"band"`              // Band
	NeiCellId    string `xml:"nei_cellid"`        // Neighbor Cell ID
	Plmn         int    `xml:"plmn"`              // Public Land Mobile Network
	Ims          int    `xml:"ims"`               // IP Multimedia Subsystem
	WdlFreqkHz   int    `xml:"wdlfreq,omitempty"` // WCDMA Downlink Frequency
	LteUlFreqkHz int    `xml:"lteulfreq"`         // LTE Uplink Frequency
	LteDlFreqkHz int    `xml:"ltedlfreq"`         // LTE Downlink Frequency
	TransMode    string `xml:"transmode"`         // Transmission Mode
	EnodebId     string `xml:"enodeb_id"`         // E-UTRAN Node B Identifier
	Cqi0         int    `xml:"cqi0"`              // Channel Quality Indicator 0
	Cqi1         int    `xml:"cqi1"`              // Channel Quality Indicator 1
	Arfcn        string `xml:"arfcn"`             // Absolute Radio Frequency Channel Number
	Bsic         string `xml:"bsic"`              // Base Station Identity Code
	RxLev        string `xml:"rxlev"`             // Received Signal Level

	PPusch         string   `xml:"-"` // Transmit Power
	PPucch         string   `xml:"-"` // Transmit Power
	PSrs           string   `xml:"-"` // Transmit Power
	PPrach         string   `xml:"-"` // Transmit Power
	EarfcnDL       string   `xml:"-"` // Downlink E-UTRA Absolute Radio Frequency Channel Number
	EarfcnUL       string   `xml:"-"` // Uplink E-UTRA Absolute Radio Frequency Channel Number
	Rsrq           *float64 `xml:"-"` // Reference Signal Received Quality
	Rsrp           *float64 `xml:"-"` // Reference Signal Received Power
	Rssi           *float64 `xml:"-"` // Received Signal Strength Indicator
	Sinr           *float64 `xml:"-"` // Signal to Interference plus Noise Ratio
	UlBandwidthMHz *int     `xml:"-"` // Uplink Bandwidth
	DlBandwidthMHz *int     `xml:"-"` // Downlink Bandwidth
	UlFrequencykHz *int     `xml:"-"` // Uplink Frequency
	DlFrequencykHz *int     `xml:"-"` // Downlink Frequency
}

var (
	Suffix = map[string]int{
		"dB":  1,
		"dBm": 1,
		"kHz": 1,
		"MHz": 1000,
	}
)

func (s *SignalInfoResp) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	type Alias SignalInfoResp
	aux := &struct {
		*Alias
		TxPower string `xml:"txpower"` // Transmit Power `PPusch:xxdBm PPucch:xxdBm PSrs:xxdBm PPrach:xxdBm`
		Rsrq    string `xml:"rsrq"`    // Reference Signal Received Quality
		Rsrp    string `xml:"rsrp"`    // Reference Signal Received Power
		Rssi    string `xml:"rssi"`    // Received Signal Strength Indicator
		Sinr    string `xml:"sinr"`    // Signal to Interference plus Noise Ratio

		UlBandwidth string `xml:"ulbandwidth"` // Uplink Bandwidth
		DlBandwidth string `xml:"dlbandwidth"` // Downlink Bandwidth

		ULFrequency string `xml:"ulfrequency"` // Uplink Frequency
		DLFrequency string `xml:"dlfrequency"` // Downlink Frequency

		Earfcn string `xml:"earfcn"` // E-UTRA Absolute Radio Frequency Channel Number `DL:xx UL:xx``

		UlMcs string `xml:"ul_mcs"` // Uplink Modulation and Coding Scheme
		DlMcs string `xml:"dl_mcs"` // Downlink Modulation and Coding Scheme
	}{
		Alias: (*Alias)(s),
	}

	err := d.DecodeElement(aux, &start)
	if err != nil {
		return err
	}

	txPowerMatches := txPowerRegex.FindStringSubmatch(aux.TxPower)
	for i, name := range txPowerRegex.SubexpNames() {
		switch name {
		case "usch":
			s.PPusch = txPowerMatches[i]
		case "ucch":
			s.PPucch = txPowerMatches[i]
		case "rs":
			s.PSrs = txPowerMatches[i]
		case "rach":
			s.PPrach = txPowerMatches[i]
		}
	}

	earfcnMatches := earfcnRegex.FindStringSubmatch(aux.Earfcn)
	for i, name := range earfcnRegex.SubexpNames() {
		switch name {
		case "dl":
			s.EarfcnDL = earfcnMatches[i]
		case "ul":
			s.EarfcnUL = earfcnMatches[i]
		}
	}

	s.Rsrq = parseFloatValue(aux.Rsrq)
	s.Rsrp = parseFloatValue(aux.Rsrp)
	s.Rssi = parseFloatValue(aux.Rssi)
	s.Sinr = parseFloatValue(aux.Sinr)

	s.UlBandwidthMHz = parseIntValue(aux.UlBandwidth)
	s.DlBandwidthMHz = parseIntValue(aux.DlBandwidth)

	s.UlFrequencykHz = parseIntValue(aux.ULFrequency)
	s.DlFrequencykHz = parseIntValue(aux.DLFrequency)

	return nil
}

func parseIntValue(value string) *int {
	var intValue int
	_, err := fmt.Sscanf(value, "%d", &intValue)
	if err != nil {
		return nil
	}

	return &intValue
}

func parseFloatValue(value string) *float64 {
	var floatValue float64
	_, err := fmt.Sscanf(value, "%f", &floatValue)
	if err != nil {
		return nil
	}

	return &floatValue
}
