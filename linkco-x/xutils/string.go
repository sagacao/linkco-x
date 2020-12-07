package xutils

// InArrayString is array string
func InArrayString(s string, arr []string) bool {
	for _, i := range arr {
		if i == s {
			return true
		}
	}
	return false
}

//Substr 字符串的截取
func Substr(str string, start int64, end int64) string {
	length := int64(len(str))
	if start < 0 || start > length {
		return ""
	}
	if end < 0 {
		return ""
	}
	if end > length {
		end = length
	}
	return string(str[start:end])
}
