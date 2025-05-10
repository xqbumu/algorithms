package tools

import (
	"algorithms/examples/mcp-go-play/internal/utils"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/tidwall/gjson"
)

func QueryWeather() server.ServerTool {
	return server.ServerTool{
		Tool: mcp.NewTool("queryWeather",
			mcp.WithDescription("输入指定城市的英文名称，返回今日天气查询结果."),
			mcp.WithString("city",
				mcp.Required(),
				mcp.Description("城市名称（需使用英文或拼音）"),
			),
		),
		Handler: handleQueryWeather,
	}
}

// API Docs: https://openweathermap.org/current
func handleQueryWeather(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	city := ""
	if b, ok := request.Params.Arguments["city"].(string); ok {
		city = b
	}

	// OpenWeather API 配置
	OPENWEATHER_API_BASE := "https://api.openweathermap.org/data/2.5/weather"
	API_KEY := os.Getenv("OPENWEATHER_API_KEY")
	USER_AGENT := "weather-app/1.0"

	// Create and send request
	var req *http.Request
	var err error
	req, err = http.NewRequest(http.MethodGet, OPENWEATHER_API_BASE, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	query := req.URL.Query()
	query.Add("q", city)
	query.Add("appid", API_KEY)
	query.Add("units", "metric")
	query.Add("lang", "zh_cn")
	req.URL.RawQuery = query.Encode()

	req.Header.Add("User-Agent", USER_AGENT)

	client := &http.Client{
		Timeout: time.Second * 30,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	// Return response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	return mcp.NewToolResultText(formatQueryWeather(respBody)), nil
}

func formatQueryWeather(body []byte) string {
	if !gjson.ParseBytes(body).IsObject() {
		return `{"error": "无法解析天气数据"}`
	}

	// 如果数据中包含错误信息，直接返回错误信息
	if s := gjson.GetBytes(body, "error").String(); len(s) > 0 {
		return fmt.Sprintf(`{"error": "%s"}`, s)
	}

	//  提取数据时做容错处理
	city := utils.JSONGetBytesByPath(body, "name", "未知")
	country := utils.JSONGetBytesByPath(body, "sys.country", "未知")
	temp := utils.JSONGetBytesByPath(body, "main.temp", "N/A")
	humidity := utils.JSONGetBytesByPath(body, "main.humidity", "N/A")
	wind_speed := utils.JSONGetBytesByPath(body, "wind.speed", "N/A")
	description := utils.JSONGetBytesByPath(body, "weather.0.description", "未知")

	result := []string{
		fmt.Sprintf("🌏 %s, %s", city, country),
		fmt.Sprintf("🌡️ 温度: %s", temp),
		fmt.Sprintf("💧 湿度: %s", humidity),
		fmt.Sprintf("🌬️ 风速: %s", wind_speed),
		fmt.Sprintf("🌤️ 天气: %s", description),
	}

	return strings.Join(result, "\n")
}
