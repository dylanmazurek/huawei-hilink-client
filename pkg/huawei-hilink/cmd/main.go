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

	// err = client.NewSession2()
	// if err != nil {
	// 	panic(err)
	// }

	err = client.Login()
	if err != nil {
		panic(err)
	}

	getSignalInfo(client)
	// getSMSInfo(client)
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
	smsInfo, err := client.GetSMSCount()
	if err != nil {
		panic(err)
	}

	tab := tabulate.New(tabulate.CompactUnicode)
	row := tab.Row()
	row.Column("Unread")
	row.Column(fmt.Sprintf("%d", *smsInfo))
	fmt.Println(tab.String())
}
