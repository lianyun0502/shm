package shm

import (
	"unsafe"
)

type ShmMemInfo struct {
	WritePtr uint
	writeLen uint

	Flag bool
	Size uint
}
var info ShmMemInfo
const InfoSize = unsafe.Sizeof(info)
