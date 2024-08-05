package shm

import (
	"fmt"
	"unsafe"

	// "sync"
	"github.com/gdygd/goshm/shmlinux"
)

type Subscriber struct {
	shm     *ShmMemInfo
	segment *shmlinux.Linuxshm
	Handle  func(data []byte)
	Data    []byte

	stopSignal chan struct{}

	flag bool
}

func NewSubscriber(skey int, shmSize int) *Subscriber {
	segmentInfo := shmlinux.NewLinuxShm()
	segmentData := shmlinux.NewLinuxShm()
	segmentInfo.InitShm(skey, int(InfoSize))
	segmentData.InitShm(skey|0x6666, shmSize)

	err := segmentInfo.CreateShm()
	if err != nil {
		fmt.Println("CreateShm err : ", err)
	}
	err = segmentData.CreateShm()
	if err != nil {
		fmt.Println("CreateShm err : ", err)
	}

	err = segmentInfo.AttachShm()
	if err != nil {
		fmt.Println("AttachShm err : ", err)
	}
	err = segmentData.AttachShm()
	if err != nil {
		fmt.Println("AttachShm err : ", err)
	}
	sharedMem := (*ShmMemInfo)(unsafe.Pointer(segmentInfo.Addr))
	p := (*byte)(unsafe.Pointer(segmentData.Addr))
	sharedMemData := unsafe.Slice(p, shmSize)
	return &Subscriber{shm: sharedMem, segment: segmentInfo, stopSignal: make(chan struct{}), Data: sharedMemData}
}

func (s *Subscriber) ReadLoop() {
	for {
		if s.shm.Flag == s.flag {
			continue
		}
		s.flag = s.shm.Flag
		data := make([]byte, s.shm.writeLen)
		fmt.Printf("Ptr : %d, Len : %d\n", s.shm.WritePtr, s.shm.writeLen)
		copy(data, s.Data[s.shm.WritePtr:s.shm.WritePtr+s.shm.writeLen])
		s.Handle(data)
	}
}

func (s *Subscriber) Close() {
	s.segment.DeleteShm()
}
