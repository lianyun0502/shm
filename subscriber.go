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
	shm         *ShmInfo
	segmentInfo *shmlinux.Linuxshm
	segmentData *shmlinux.Linuxshm
	Handle      func(data []byte)
	Data        []byte

	stopSignal chan struct{}
	sysSignal  chan os.Signal

	dataCH chan []byte

	startFlag   bool
	preWritePtr uint
}

func NewSubscriber(skey int, shmSize int) *Subscriber {
	segmentInfo, _ := NewSegment(skey, int(InfoSize))
	segmentData, _ := NewSegment(skey|0x6666, shmSize)

	shmInfo := (*ShmInfo)(unsafe.Pointer(segmentInfo.Addr))
	p := (*byte)(unsafe.Pointer(segmentData.Addr))
	shmData := unsafe.Slice(p, shmSize)

	subscriber := &Subscriber{
		shm:         shmInfo,
		segmentInfo: segmentInfo,
		segmentData: segmentData,
		stopSignal:  make(chan struct{}),
		sysSignal:   make(chan os.Signal, 1),
		dataCH:      make(chan []byte),
		Data:        shmData,
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
		Logger.Debugf("Ptr : %d, Len : %d, MsgID : %d", s.shm.WritePtr, s.shm.writeLen, s.shm.MsgID)
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
