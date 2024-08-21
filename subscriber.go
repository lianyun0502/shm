package shm

import (
	"encoding/binary"
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

	startFlag   bool
	preWritePtr uint
}

func NewSubscriber(skey int, shmSize int) *Subscriber {
	segmentInfo := shmlinux.NewLinuxShm()
	segmentData := shmlinux.NewLinuxShm()
	segmentInfo.InitShm(skey, int(InfoSize))
	segmentData.InitShm(skey|0x6666, shmSize)

	err := segmentInfo.CreateShm()
	if err != nil {
		Logger.Warning("CreateShm err : ", err)
	}
	err = segmentData.CreateShm()
	if err != nil {
		Logger.Warning("CreateShm err : ", err)
	}

	err = segmentInfo.AttachShm()
	if err != nil {
		Logger.Warning("AttachShm err : ", err)
	}
	err = segmentData.AttachShm()
	if err != nil {
		Logger.Warning("AttachShm err : ", err)
	}
	sharedMem := (*ShmMemInfo)(unsafe.Pointer(segmentInfo.Addr))
	p := (*byte)(unsafe.Pointer(segmentData.Addr))
	sharedMemData := unsafe.Slice(p, shmSize)
	return &Subscriber{
		shm:         sharedMem,
		segment:     segmentInfo,
		stopSignal:  make(chan struct{}),
		Data:        sharedMemData,
		startFlag:   false,
		preWritePtr: 0,
	}
}

func (s *Subscriber) ReadLoop() {
	for {
		if (s.preWritePtr == s.shm.WritePtr) && s.startFlag {
			continue
		}
		s.startFlag = true
		s.preWritePtr = s.shm.WritePtr
		Logger.Debugf("Ptr : %d, write len : %d", s.shm.WritePtr, s.shm.writeLen)

		dataLen := binary.BigEndian.Uint32(s.Data[s.shm.WritePtr : s.shm.WritePtr+4])
		Logger.Debugf("Data length: %d", dataLen)

		data := make([]byte, s.shm.writeLen)
		copy(data, s.Data[s.shm.WritePtr+4:s.shm.WritePtr+s.shm.writeLen])
		s.Handle(data)
	}
}

func (s *Subscriber) Close() {
	s.segment.DeleteShm()
}
