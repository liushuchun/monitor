package main

import "testing"
import "fmt"

func TestSwapMemory(t *testing.T) {

	memStatus, err := SwapMemory()
	if err != nil {
		t.Error("Failed get the memory status")
	}
	fmt.Println(memStatus)
}
