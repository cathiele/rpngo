package tinygo

import (
	"runtime"
	"strconv"
)

func DumpMemStats() {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	print("alloc: ")
	print(strconv.Itoa(int(ms.Alloc)))
	print(" sys: ")
	print(strconv.Itoa(int(ms.Sys)))
	print(" free: ")
	println(strconv.Itoa(int(ms.HeapIdle)))
}
