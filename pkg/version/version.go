// Package version Copyright 2021 damonchen, netubu@gmail.com
//
package version

import (
	"strconv"
	"strings"
)

var version string = "0.1.0"

func Full() string {
	return version
}

func getSubVersion(v string, position int) int64 {
	arr := strings.Split(v, ".")
	if len(arr) < 3 {
		return 0
	}
	res, _ := strconv.ParseInt(arr[position], 10, 64)
	return res
}

func Proto(v string) int64 {
	return getSubVersion(v, 0)
}

func Major(v string) int64 {
	return getSubVersion(v, 1)
}

func Minor(v string) int64 {
	return getSubVersion(v, 2)
}
