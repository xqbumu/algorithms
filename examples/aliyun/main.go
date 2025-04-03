package main

import (
	"fmt"
	"os"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
)

func main() {
	config := sdk.NewConfig()

	// Please ensure that the environment variables ALIBABA_CLOUD_ACCESS_KEY_ID and ALIBABA_CLOUD_ACCESS_KEY_SECRET are set.
	credential := credentials.NewAccessKeyCredential(os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_ID"), os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_SECRET"))
	/* use STS Token
	credential := credentials.NewStsTokenCredential(os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_ID"), os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_SECRET"), os.Getenv("ALIBABA_CLOUD_SECURITY_TOKEN"))
	*/
	client, err := sdk.NewClientWithOptions("cn-hangzhou", config, credential)
	if err != nil {
		panic(err)
	}

	request := requests.NewCommonRequest()

	request.Method = "POST"
	request.Scheme = "https" // https | http
	request.Domain = "nlp-automl.cn-hangzhou.aliyuncs.com"
	request.Version = "2019-11-11"
	request.ApiName = "RunPreTrainService"
	request.QueryParams["ServiceName"] = "DeepEmotionBert"
	request.QueryParams["ServiceVersion"] = "V1"
	request.QueryParams["PredictContent"] = "{\"input\":{\"content\": \"图形验证码输入老是错误\", \"domain\": \"clothing\"}}"

	response, err := client.ProcessCommonRequest(request)
	if err != nil {
		panic(err)
	}
	fmt.Print(response.GetHttpContentString())
}
