package utils

import "strconv"

func GetNumDigits(i int) int {
	return len(strconv.Itoa(i))
}
