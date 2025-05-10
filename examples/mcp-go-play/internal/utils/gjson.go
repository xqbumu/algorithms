package utils

import "github.com/tidwall/gjson"

func JSONGetBytesByPath(json []byte, path, final string) string {
	result := gjson.GetBytes(json, path)
	if !result.Exists() {
		return final
	}
	return result.String()
}
