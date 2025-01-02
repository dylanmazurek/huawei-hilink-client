package main

import (
	"context"

	huaweihilink "github.com/dylanmazurek/huawei-wifi/pkg/huawei-hilink"
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

	_, err = client.SignalInfo()
	if err != nil {
		panic(err)
	}
}
