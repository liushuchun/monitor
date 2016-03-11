package main

import (
	"fmt"
	"testing"
)

func TestDiskUsage(t *testing.T) {
	diskstatus, err := DiskUsage("/")
	if err != nil {
		t.Error("error")
	}
	fmt.Println(diskstatus)
}
