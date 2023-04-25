package utils

import (
	"strconv"
	"strings"
)

func CompareVersions(v1 string, v2 string) int {
	v1List := strings.Split(v1, ".")
	v2List := strings.Split(v2, ".")
	for i := 0; i < len(v1List) && i < len(v2List); i++ {
		v1Int, _ := strconv.Atoi(v1List[i])
		v2Int, _ := strconv.Atoi(v2List[i])
		if v1Int < v2Int {
			return -1
		} else if v1Int > v2Int {
			return 1
		}
	}
	if len(v1List) < len(v2List) {
		return -1
	} else if len(v1List) > len(v2List) {
		return 1
	}
	return 0
}
