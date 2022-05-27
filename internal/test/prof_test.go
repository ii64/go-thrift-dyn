package test

import (
	"fmt"
	"os"
	"runtime/pprof"
	"time"
)

func profCpu() func() {
	f, err := os.Create(fmt.Sprintf("%d.cpu", time.Now().Unix()))
	if err != nil {
		panic(err)
	}
	if err = pprof.StartCPUProfile(f); err != nil {
		panic(err)
	}
	return pprof.StopCPUProfile
}
