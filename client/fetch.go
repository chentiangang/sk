package main

import (
	"fmt"
	"runtime"

	"github.com/shirou/gopsutil/mem"

	"github.com/shirou/gopsutil/cpu"
)

func NewItems() {
	// mem
	mInfo, _ := mem.VirtualMemory()
	fmt.Printf("mem Info: #%s\n", mInfo.String())

	// cpu
	cpuInfos, err := cpu.Info()
	if err != nil {
		panic(err)
	}
	// Printf("info: %s, cores: %d\n", cpuInfos[0].ModelName, cpuInfos[0].CacheSize, runtime.NumCPU())
	fmt.Printf("info: %s, cachaSize: %dk, cores: %d\n", cpuInfos[0].ModelName, cpuInfos[0].CacheSize, runtime.NumCPU())
}

func main() {
	NewItems()
}
