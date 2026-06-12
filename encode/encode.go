// Package encode provides helpers that format values into RESP
// (Redis Serialization Protocol) wire-format strings.
package encode

import "strconv"

// EncodeSimpleString encodes s as a RESP simple string: "+<s>\r\n".
func EncodeSimpleString(s string) string {
	return "+" + s + "\r\n"
}

// EncodeError encodes s as a RESP error: "-<s>\r\n".
func EncodeError(s string) string {
	return "-" + s + "\r\n"
}

// EncodeNull returns the RESP null bulk-string "$-1\r\n", used when a key is missing.
func EncodeNull() string {
	return "$-1\r\n"
}

// EncodeBulkString encodes s as a RESP bulk string: "$<len>\r\n<s>\r\n".
func EncodeBulkString(s string) string {
	return "$" + strconv.Itoa(len(s)) + "\r\n" + s + "\r\n"
}

// EncodeInteger encodes i as a RESP integer: ":<i>\r\n".
func EncodeInteger(i int) string {
	return ":" + strconv.Itoa(i) + "\r\n"
}

// EncodeArray encodes a slice of strings as a RESP array of bulk strings.
func EncodeArray(l []string) string {
	result := "*" + strconv.Itoa(len(l)) + "\r\n"
	for _, s := range l {
		result += EncodeBulkString(s)
	}
	return result
}

// EncodeArrayWithNulls encodes a mixed slice as a RESP array where nil elements
// become null bulk strings and string elements become bulk strings.
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
