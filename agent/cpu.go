package agent

import (
	"encoding/json"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

const (
	CPUser = iota
	CPNice
	CPSys
	CPIntr
	CPIdle
	CPUStates
)

type CPUTimesStat struct {
	CPU       string  `json:"cpu"`
	User      float64 `json:"user"`
	System    float64 `json:"system"`
	Idle      float64 `json:"idle"`
	Nice      float64 `json:"nice"`
	Iowait    float64 `json:"iowait"`
	Irq       float64 `json:"irq"`
	Softirq   float64 `json:"softirq"`
	Steal     float64 `json:"steal"`
	Guest     float64 `json:"guest"`
	GuestNice float64 `json:"guest_nice"`
	Stolen    float64 `json:"stolen"`
}

type CPUInfoStat struct {
	CPU        int32    `json:"cpu"`
	VendorID   string   `json:"vendor_id"`
	Family     string   `json:"family"`
	Model      string   `json:"model"`
	Stepping   int32    `json:"stepping"`
	PhysicalID string   `json:"physical_id"`
	CoreID     string   `json:"core_id"`
	Cores      int32    `json:"cores"`
	ModelName  string   `json:"model_name"`
	Mhz        float64  `json:"mhz"`
	CacheSize  int32    `json:"cache_size"`
	Flags      []string `json:"flags"`
}

var lastCPUTimes []CPUTimesStat
var lastPerCPUTimes []CPUTimesStat

func perCPUTimes() ([]CPUTimesStat, error) {
	return []CPUTimesStat{}, nil
}

func allCPUTimes() ([]CPUTimesStat, error) {
	return []CPUTimesStat{}, nil
}

func CPUCounts(logical bool) (int, error) {
	return runtime.NumCPU(), nil
}

func (c CPUTimesStat) String() string {
	v := []string{
		`"cpu":"` + c.CPU + `"`,
		`"user":` + strconv.FormatFloat(c.User, 'f', 1, 64),
		`"system":` + strconv.FormatFloat(c.System, 'f', 1, 64),
		`"idle":` + strconv.FormatFloat(c.Idle, 'f', 1, 64),
		`"nice":` + strconv.FormatFloat(c.Nice, 'f', 1, 64),
		`"iowait":` + strconv.FormatFloat(c.Iowait, 'f', 1, 64),
		`"irq":` + strconv.FormatFloat(c.Irq, 'f', 1, 64),
		`"softirq":` + strconv.FormatFloat(c.Softirq, 'f', 1, 64),
		`"steal":` + strconv.FormatFloat(c.Steal, 'f', 1, 64),
		`"guest":` + strconv.FormatFloat(c.Guest, 'f', 1, 64),
		`"guest_nice":` + strconv.FormatFloat(c.GuestNice, 'f', 1, 64),
		`"stolen":` + strconv.FormatFloat(c.Stolen, 'f', 1, 64),
	}

	return `{` + strings.Join(v, ",") + `}`
}

// Total returns the total number of seconds in a CPUTimesStat
func (c CPUTimesStat) Total() float64 {
	total := c.User + c.System + c.Nice + c.Iowait + c.Irq + c.Softirq + c.Steal +
		c.Guest + c.GuestNice + c.Idle + c.Stolen
	return total
}

func (c CPUInfoStat) String() string {
	s, _ := json.Marshal(c)
	return string(s)
}

var ClocksPerSec = float64(128)

func CPUTimes(percpu bool) ([]CPUTimesStat, error) {
	if percpu {
		return perCPUTimes()
	}
	return allCPUTimes()
}

func CPUInfo() ([]CPUInfoStat, error) {
	var ret []CPUInfoStat
	out, err := exec.Command("/usr/sbin/sysctl", "machdep.cpu").Output()
	if err != nil {
		return ret, err
	}
	c := CPUInfoStat{}
	for _, line := range strings.Split(string(out), "\n") {
		values := strings.Fields(line)
		if len(values) < 1 {
			continue
		}
		t, err := strconv.ParseInt(values[1], 10, 64)
		switch {
		case strings.HasPrefix(line, "machdep.cpu.brand_string"):
			c.ModelName = strings.Join(values[1:], " ")
		case strings.HasPrefix(line, "machdep.cpu.family"):
			c.Family = values[1]
		case strings.HasPrefix(line, "machdep.cpu.model"):
			c.Model = values[1]
		case strings.HasPrefix(line, "machdep.cpu.stepping"):

			if err != nil {
				return ret, err
			}
			c.Stepping = int32(t)
		case strings.HasPrefix(line, "machdep.cpu.features"),
			strings.HasPrefix(line, "machdep.cpu.leaf7_features"),
			strings.HasPrefix(line, "machdep.cpu.extfeatures"):
			for _, v := range values[1:] {
				c.Flags = append(c.Flags, strings.ToLower(v))
			}
		case strings.HasPrefix(line, "machdep.cpu.core_count"):
			if err != nil {
				return ret, err

			}
			c.Cores = int32(t)
		case strings.HasPrefix(line, "machdep.cpu.cache.size"):
			if err != nil {
				return ret, err
			}
			c.CacheSize = int32(t)
		case strings.HasPrefix(line, "machdep.cpu.vendor"):
			c.VendorID = values[1]
		}
	}
	out, err = exec.Command("/usr/sbin/sysctl", "hw.cpufrequency").Output()
	if err != nil {
		return ret, err
	}
	values := strings.Fields(string(out))
	mhz, err := strconv.ParseFloat(values[1], 64)
	if err != nil {
		return ret, err
	}
	c.Mhz = mhz / 1000000.0
	return append(ret, c), nil
}
