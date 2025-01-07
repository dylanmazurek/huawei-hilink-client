package main

import (
	"context"
	"fmt"
	"reflect"

	huaweihilink "github.com/dylanmazurek/huawei-hilink-client/pkg/huawei-hilink"
	"github.com/markkurossi/tabulate"
)

func main() {
	ctx := context.Background()

	client, err := huaweihilink.New(ctx)
	if err != nil {
		panic(err)
	}

	err = client.Login()
	if err != nil {
		panic(err)
	}

	sendSMS(client)
}

func getSignalInfo(client *huaweihilink.Client) {
	signalInfo, err := client.SignalInfo()
	if err != nil {
		panic(err)
	}

	tab := tabulate.New(tabulate.CompactUnicode)

	v := reflect.ValueOf(*signalInfo)

	for i := 0; i < v.NumField(); i++ {
		row := tab.Row()
		fieldName := v.Type().Field(i).Name
		row.Column(fieldName)

		if v.Field(i).Kind() != reflect.Ptr {
			valueStr := fmt.Sprintf("%v", v.Field(i).Interface())
			row.Column(valueStr)
		} else {
			row.Column("-")
		}
	}

	fmt.Println(tab.String())
}

func getSMSInfo(client *huaweihilink.Client) {
	smsPages, err := client.GetSMSPages()
	if err != nil {
		panic(err)
	}

	fmt.Println("SMS Pages:", *smsPages)
}

func getSMSPage(client *huaweihilink.Client) {
	smsPages, err := client.GetSMSPages()
	if err != nil {
		panic(err)
	}

	fmt.Println("hash:", smsPages.Hash)
	fmt.Println("pwd:", smsPages.Pwd)
}

func sendSMS(client *huaweihilink.Client) {
	_, err := client.PostSMS("+61413073143", "hello how are you?")
	if err != nil {
		panic(err)
	}
}
