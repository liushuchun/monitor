package main

import (
	"gopkg.in/mgo.v2"
	//"gopkg.in/mgo.v2/bson"
	"fmt"
	"log"
	"time"
)

var (
	IsDrop = true
)

func SaveInfoToDB(mem MemStatus, cpu CPUInfoStat, disk DiskStatus) error {
	host, port, userName, passWord := GetDBConfig()
	fmt.println(host)
	session, err := mgo.Dial(host)
	if err != nil {
		panic(err)
	}

	defer session.Close()
	session.SetMode(mog.Monotonic, true)
	if IsDrop {
		err = session.DB("monitor").DropDatabase()
		if err != nil {
			panic(err)
		}
	}
	return
}

func main() {

	cpuList := make(chan CPUInfoStat)
	memList := make(chan MemStatus)
	diskList := make(chan DiskStatus)

	go func() {
		for {
			cpu, err := CPUInfo()
			cpuList <- cpu
			mem, err := SwapMemory()
			memList <- mem
			disk, err := DiskUsage("/")
			diskList <- disk
			time.Sleep(10 * time.Second)
		}
	}()

	go func() {
		for mem := range memList {
		cpu:
			<-cpuList
		disk:
			<-diskList
			err := SaveInfoToDB(mem, cpu, disk)
			if err != nil {
				log.Fatalf(err)
			}
		}

	}()

}
