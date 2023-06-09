package debugging

import (
	"encoding/json"
	"fmt"
	"time"
)

type Debug struct {
	TimesSpent map[string]executionTime `json:"timesSpent"`
	TotalTime  float64                  `json:"totalTime"`
	Distance  float64                  `json:"distance"`
}

type executionTime struct {
	Time float64 `json:"time"`
	Call int     `json:"call"`
	TimeByCall float64 `json:"timeByCall"`
}

func NewDebug() *Debug {
	return &Debug{
		TimesSpent: make(map[string]executionTime),
	}
}

func (d *Debug) GetTimeUsing(funcname string, f func()) {
	start := time.Now()
	f()
	elapsed := time.Since(start)
	if _, exists := d.TimesSpent[funcname]; !exists {
		d.TimesSpent[funcname] = executionTime{
			Call: 1,
			Time: elapsed.Seconds(),
		}
	} else {
		d.TimesSpent[funcname] = executionTime{
			Call: d.TimesSpent[funcname].Call + 1,
			Time: d.TimesSpent[funcname].Time + elapsed.Seconds(),
		}
	}
}

func (d *Debug) Print() {
	d.SetTimeByCall()
	maxLen := 0
	for k := range d.TimesSpent {
		if len(k) > maxLen {
			maxLen = len(k)
		}
	}
	fmt.Printf("{\n")
	for k, v := range d.TimesSpent {
		nameFormat := fmt.Sprintf("%%%ds", maxLen)
		fmt.Printf("   Action: "+nameFormat+" => took %9f s   in %9d calls\n, meaning ", k, v.Time, v.Call)
	}
	fmt.Printf("}\n")
	d.Reset()
}

func (d *Debug) JSON() []byte {
	d.SetTimeByCall()
	res, _ := json.Marshal(d)
	return res
}

func (d *Debug) Reset() {
	d.TimesSpent = make(map[string]executionTime)
}

func (d *Debug) SetTimeByCall() {
	for k, v := range d.TimesSpent {
		d.TimesSpent[k] = executionTime{
			Call: v.Call,
			Time: v.Time,
			TimeByCall: v.Time / float64(v.Call) * 1000000000,
		}
	}
}
