package shm

import (
	"bytes"
	"os"
	"os/signal"
	"unsafe"

	// "sync"
	"github.com/gdygd/goshm/shmlinux"
)

type Subscriber struct {
	shm     *ShmMemInfo
	segmentInfo *shmlinux.Linuxshm
	segmentData *shmlinux.Linuxshm
	Handle  func(data []byte)
	Data    []byte

	stopSignal chan struct{}
	sysSignal  chan os.Signal

	dataCH 	chan []byte

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

	subscriber := &Subscriber{
		shm:         sharedMem,
		segmentInfo:     segmentInfo,
		segmentData:     segmentData,
		stopSignal:  make(chan struct{}),
		sysSignal:   make(chan os.Signal, 1),
		dataCH: 	make(chan []byte),
		Data:        sharedMemData,
		startFlag:   false,
		preWritePtr: 0,
	}

	signal.Notify(subscriber.sysSignal, os.Interrupt)

	go func() {
		sig := <-subscriber.sysSignal
		if sig == os.Interrupt {
			segmentInfo.DeleteShm()
			segmentData.DeleteShm()
		}
	}()

	return subscriber
}

func (s *Subscriber) ReadLoop() {
	for {
		if (s.preWritePtr == s.shm.WritePtr) && s.startFlag {
			continue
		}
		s.startFlag = true
		s.preWritePtr = s.shm.WritePtr
		data := make([]byte, s.shm.writeLen)
		Logger.Debugf("Ptr : %d, Len : %d", s.shm.WritePtr, s.shm.writeLen)
		if bytes.Equal(data, []byte("EOF")) {
			s.Close()
			return
		}
		copy(data, s.Data[s.shm.WritePtr:s.shm.WritePtr+s.shm.writeLen])
		s.Handle(data)
	}
}

func (s *Subscriber) Close() {
	s.segmentInfo.DeleteShm()
	s.segmentData.DeleteShm()
	Logger.Info("Subscriber Close")
}

