package main

import (
	"fmt"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	hunyuan "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/hunyuan/v20230901"
)

func main() {
	credential := common.NewCredential(
		os.Getenv("TENCENT_SECRET_ID"),
		os.Getenv("TENCENT_SECRET_KEY"),
	)

	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "hunyuan.tencentcloudapi.com"

	client, _ := hunyuan.NewClient(credential, "ap-guangzhou", cpf)

	// 实例化一个请求对象,每个接口都会对应一个request对象

	request := hunyuan.NewChatCompletionsRequest()

	request.Model = common.StringPtr("hunyuan-lite")
	request.Messages = []*hunyuan.Message{
		{
			Role:    common.StringPtr("user"),
			Content: common.StringPtr("hi"),
		},
	}
	request.Stream = common.BoolPtr(false)
	request.StreamModeration = common.BoolPtr(false)

	// 返回的resp是一个ChatCompletionsResponse的实例，与请求对象对应

	response, err := client.ChatCompletions(request)

	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		fmt.Printf("An API error has returned: %s", err)
		return
	}

	if err != nil {
		panic(err)
	}

	// 输出json格式的字符串回包
	if response.Response != nil {
		// 非流式响应
		fmt.Println(response.ToJsonString())
	} else {
		// 流式响应
		for event := range response.Events {
			fmt.Println(string(event.Data))
		}
	}
}
