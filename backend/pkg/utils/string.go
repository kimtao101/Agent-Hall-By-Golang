// utils 提供通用工具函数
package utils

// TruncateString 截断字符串到指定长度
func TruncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen]
}

// GetStringValue 从 map 中获取字符串值，如果不存在返回默认值
func GetStringValue(config map[string]interface{}, key string, defaultValue string) string {
	if val, ok := config[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return defaultValue
}
