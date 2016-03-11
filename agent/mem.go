package main

import (
	"encoding/binary"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"unsafe"
)

type MemStatus struct {
	Total       uint64  `json:"all"`
	Used        uint64  `json:"used"`
	Free        uint64  `json:"free"`
	UsedPercent float64 `json:"usedpercent"`
}

func DoSysctrl(mib string) ([]string, error) {
	os.Setenv("LC_ALL", "C")
	out, err := exec.Command("/usr/sbin/sysctl", "-n", mib).Output()
	if err != nil {
		return []string{}, err
	}
	v := strings.Replace(string(out), "{ ", "", 1)
	v = strings.Replace(string(v), " }", "", 1)
	values := strings.Fields(string(v))

	return values, nil
}

func CallSyscall(mib []int32) ([]byte, uint64, error) {
	miblen := uint64(len(mib))

	// get required buffer size
	length := uint64(0)
	_, _, err := syscall.Syscall6(
		syscall.SYS___SYSCTL,
		uintptr(unsafe.Pointer(&mib[0])),
		uintptr(miblen),
		0,
		uintptr(unsafe.Pointer(&length)),
		0,
		0)
	if err != 0 {
		var b []byte
		return b, length, err
	}
	if length == 0 {
		var b []byte
		return b, length, err
	}
	// get proc info itself
	buf := make([]byte, length)
	_, _, err = syscall.Syscall6(
		syscall.SYS___SYSCTL,
		uintptr(unsafe.Pointer(&mib[0])),
		uintptr(miblen),
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(unsafe.Pointer(&length)),
		0,
		0)
	if err != 0 {
		return buf, length, err
	}

	return buf, length, nil
}

func getMemsize() (uint64, error) {
	totalString, err := syscall.Sysctl("hw.memsize")
	if err != nil {
		return 0, err
	}

	totalString += "\x00"
	total := uint64(binary.LittleEndian.Uint64([]byte(totalString)))
	return total, nil
}

func SwapMemory() (*MemStatus, error) {
	var ret *MemStatus
	swapUsage, err := DoSysctrl("vm.swapusage")
	if err != nil {
		return ret, err
	}

	total := strings.Replace(swapUsage[2], "M", "", 1)
	used := strings.Replace(swapUsage[5], "M", "", 1)
	free := strings.Replace(swapUsage[8], "M", "", 1)
	total_v, err := strconv.ParseFloat(total, 64)
	if err != nil {
		return nil, err
	}
	used_v, err := strconv.ParseFloat(used, 64)
	if err != nil {
		return nil, err
	}
	free_v, err := strconv.ParseFloat(free, 64)
	if err != nil {
		return nil, err
	}
	u := float64(0)
	if total_v != 0 {
		u = ((total_v - free_v) / total_v) * 100.0
	}
	ret = &MemStatus{
		Total:       uint64(total_v * 1000),
		Used:        uint64(used_v * 1000),
		Free:        uint64(free_v * 1000),
		UsedPercent: u,
	}
	return ret, nil
}
