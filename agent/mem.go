package mem

import (
	"log"
	"os"
	"strconv"
	"time"
)

type Mem struct {
	UsedSize  int64
	TotalSize int64
	Time      time.Time
}

func (m *Mem) GetInfo() (resTime *Mem) {

}
