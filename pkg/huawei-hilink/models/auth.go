package models

import "encoding/xml"

type SessionTokenInfo struct {
	XMLName     xml.Name `xml:"response"`
	SessionInfo string   `xml:"SesInfo"`
	TokenInfo   string   `xml:"TokInfo"`
}

type TokenResp struct {
	XMLName xml.Name `xml:"response"`
	Token   string   `xml:"token"`
}

type TokenChallengeReq struct {
	XMLName    xml.Name `xml:"request"`
	Username   string   `xml:"username"`
	FirstNonce string   `xml:"firstnonce"`
	Mode       int      `xml:"mode"`
}

type ChallengeResp struct {
	XMLName      xml.Name `xml:"response"`
	Salt         string   `xml:"salt"`
	ModeSelected int      `xml:"mode_selected"`
	ServerNonce  string   `xml:"servernonce"`
	NewType      int      `xml:"newtype"`
	Iterations   int      `xml:"iterations"`
}

type ClientProofReq struct {
	XMLName     xml.Name `xml:"request"`
	ClientProof string   `xml:"clientproof"`
	FinalNonce  string   `xml:"finalnonce"`
}

type AuthenticationLoginResp struct {
	XMLName            xml.Name `xml:"response"`
	ServerSignature    string   `xml:"serversignature"`
	RsaPubKeySignature string   `xml:"rsapubkeysignature"`
	RsaE               string   `xml:"rsae"`
	RsaN               string   `xml:"rsan"`
}
