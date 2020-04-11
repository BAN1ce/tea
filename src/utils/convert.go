package utils

import "strconv"

func ConvertToBin(num int) string {
	s := ""

	if num == 0 {
		return "0"
	}

	// num /= 2 每次循环的时候 都将num除以2  再把结果赋值给 num
	for ; num > 0; num /= 2 {
		lsb := num % 2
		// strconv.Itoa() 将数字强制性转化为字符串
		s = strconv.Itoa(lsb) + s
	}
	return s

}

/**
utf 计算字节长度
*/
func UtfLength(l []byte) int {

	return int(l[0])*512 + int(l[1])
}
