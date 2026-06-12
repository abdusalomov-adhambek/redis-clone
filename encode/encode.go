package encode

import "strconv"

func EncodeSimpleString(s string) string {
	return "+" + s + "\r\n"
}

func EncodeError(s string) string {
	return "-" + s + "\r\n"
}

func EncodeNull() string {
	return "$-1\r\n"
}

func EncodeBulkString(s string) string {
	return "$" + strconv.Itoa(len(s)) + "\r\n" + s + "\r\n"
}

func EncodeInteger(i int) string {
	return ":" + strconv.Itoa(i) + "\r\n"
}

func EncodeArray(l []string) string {
	result := "*" + strconv.Itoa(len(l)) + "\r\n"
	for _, s := range l {
		result += EncodeBulkString(s)
	}
	return result
}

func EncodeArrayWithNulls(l []any) string {
	result := "*" + strconv.Itoa(len(l)) + "\r\n"
	for _, s := range l {
		if s == nil {
			result += EncodeNull()
		} else {
			result += EncodeBulkString(s.(string))
		}
	}
	return result
}
