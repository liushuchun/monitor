package main

import (
	"fmt"
	"testing"
)

func TestCPUInfo(t *testing.T) {
	cpuInfo, err := CPUInfo()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(cpuInfo)

}
