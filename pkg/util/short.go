package util

import (
	"math"
	"strconv"
	"strings"
)

var binaryConversionMap = map[int]string{
	0: "0", 1: "1", 2: "2", 3: "3", 4: "4", 5: "5",
	6: "6", 7: "7", 8: "8", 9: "9", 10: "a", 11: "b",
	12: "c", 13: "d", 14: "e", 15: "f", 16: "g", 17: "h",
	18: "i", 19: "j", 20: "k", 21: "l", 22: "m", 23: "n",
	24: "o", 25: "p", 26: "q", 27: "r", 28: "s", 29: "t",
	30: "u", 31: "v", 32: "w", 33: "x", 34: "y", 35: "z",
}

var encodeURL32 = "0123456789abcdefghijklmnopqrstuvwxyz"

type binaryConvert struct {
	// 转换规则
	ConvertRegx string
	// 进制
	len int
}

func NewBinaryConvert(str string) *binaryConvert {
	b := &binaryConvert{}
	if len(str) < 1 {
		b.ConvertRegx = encodeURL32
	} else {
		b.ConvertRegx = str
	}
	b.len = len(b.ConvertRegx)
	return b
}

var BinaryConvert = NewBinaryConvert(encodeURL32)

func (b *binaryConvert) Encode(num int64) string {
	n := int64(b.len)
	newNumStr := ""
	var remainder int64
	var remainderString string
	for num != 0 {
		remainder = num % n
		if int64(b.len+1) > remainder && remainder > 9 {
			// remainderString = binaryConversionMap[remainder]
			remainderString = string(encodeURL32[remainder])
		} else {
			remainderString = strconv.FormatInt(remainder, 10)
		}
		newNumStr = remainderString + newNumStr
		num = num / n
	}
	return newNumStr
}

func binaryConversionKey(in string) int {
	result := -1
	for k, v := range binaryConversionMap {
		if in == v {
			result = k
		}
	}
	return result
}

func (b *binaryConvert) Decode(num string) string {
	var newNum int64
	nNum := len(strings.Split(num, "")) - 1
	for _, value := range strings.Split(num, "") {
		tmp := binaryConversionKey(value)
		if tmp != -1 {
			newNum = newNum + int64(tmp)*int64(math.Pow(float64(b.len), float64(nNum)))
			nNum--
		} else {
			break
		}
	}
	return strconv.FormatInt(newNum, 10)
}
