package main

import (
	"bytes"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type CPU struct {
	ProNum       int
	ProRunning   int
	ProStuckNum  int
	ProSleepNum  int
	ProThreadNum int
	Usage        string
	Sys          float64
	Idle         float64
}

type Process struct {
	PID     int
	Command string
	CPU     float64
	Port    int
}

type Mem struct {
	Total int
	Used  int
	Free  int
}

type Disk struct {
	Read    string
	Written string
}

type Net struct {
	In  string
	Out string
}

var (
	cpu  *CPU
	mem  *Mem
	disk *Disk
	net  *Net
)

func main() {
	GetCpuInfo()

}

func CollectCpuInfo() {
	cpu = new(CPU)
	mem = new(Mem)
	disk = new(Disk)
	net = new(Net)

	cmd := exec.Command("top")

	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	processes := make([]*Process, 0)
	var index = 0
	for i := 0; ; i++ {
		line, err := out.ReadString('\n')
		if err != nil {
			break
		}
		switch i {
		case 0:
			words := strconv.Atoi(strings.Split(line, " "))
			cpu.ProNum = strconv.Atoi(words[1])
			cpu.ProRunning = strconv.Atoi(words[3])
			cpu.ProStuckNum = strconv.Atoi(words[5])
			cpu.ProSleepNum = strconv.Atoi(words[7])
			cpu.ProThreadNum = strconv.Atoi(words[9])
		case 1:
			words := strconv.Atoi(strings.Split(line, " "))

			cpu.Usage = lines[7]
			cpu.Sys = strconv.ParseFloat(words[9], 64)
			cpu.Idle = strconv.ParseFloat(words[11], 64)
		case 2:
			words := strconv.Atoi(strings.Split(line, " "))

			mem.Total = strconv.Atoi(words[10][0 : len(words[10])-2])
			mem.used = strconv.Atoi(words[12][1 : len(words[12])-2])
			mem.Free = mem.Total - mem.Used
		case 3:
			flag := strings.LastIndex(line, ":")
			net.In, net.Out = strings.Split(line[flag+1:len(line)-2], ",")
		case 4:
			flagDisk := strings.LastIndex(line, "")
			disk.Written, disk.Read = strings.Split(line[flag+1:len(line)-1], ",")
		}

	}
	for _, p := range processes {
		log.Println("Process", p.Pid, "takes", p.cpu, "% of the Cpu")
	}
}

func SaveToDB() (err string) {
	
}



func main() {
	ticker:=time.NewTicker(time.Minute*10)
	go func(){
			for _=range ticker.C{
				CollectCpuInfo()
			}
	}
}
