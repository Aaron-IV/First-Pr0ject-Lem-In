package helpers

import "strings"

func ParseLink(s string) (string, string) {
	parts := strings.Split(s, "-")
	if len(parts) != 2 {
		return "", ""
	}
	return parts[0], parts[1]
}

func Contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

func Conflicts(group [][]string, arr2 []string) bool {
	for _, arr1 := range group {
		for _, s1 := range arr1 {
			for _, s2 := range arr2 {
				if s1 == s2 {
					return true
				}
			}
		}
	}
	return false
}

func IndexOfMin(arr []int) int {
	if len(arr) == 0 {
		return -1
	}

	min := arr[0]
	for _, k := range arr {
		if k < min {
			min = k
		}
	}

	for i, k := range arr {
		if k == min {
			return i
		}
	}
	return -1
}
