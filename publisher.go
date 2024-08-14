package shm

import (
	"os"
	"os/signal"
	"unsafe"

	"github.com/gdygd/goshm/shmlinux"
)

type Publisher struct {
	shmInfo *ShmMemInfo
	shmData []byte
	segment *shmlinux.Linuxshm

	sysSignal  chan os.Signal
	DoneSignal chan struct{}

	IsClosed bool
}

func NewPublisher(skey int, shmSize int) *Publisher {
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
	sharedMemInfo := (*ShmMemInfo)(unsafe.Pointer(segmentInfo.Addr))
	sharedMemInfo.Size = uint(shmSize)
	p := (*byte)(unsafe.Pointer(segmentData.Addr))
	sharedMemData := unsafe.Slice(p, shmSize)

	publisher := &Publisher{
		shmInfo:   sharedMemInfo,
		shmData:   sharedMemData,
		segment:   segmentInfo,
		sysSignal: make(chan os.Signal),
	}

	signal.Notify(publisher.sysSignal, os.Interrupt)

	go func() {
		sig := <-publisher.sysSignal
		if sig == os.Interrupt {
			segmentInfo.DeleteShm()
		}
	}()
	return publisher
}

func (p *Publisher) Write(data []byte) {
	if p.IsClosed {
		return
	}
	dataLen := uint(len(data))
	if p.shmInfo.WritePtr+p.shmInfo.writeLen+dataLen > uint(p.shmInfo.Size) {
		p.shmInfo.WritePtr = 0
	} else {
		p.shmInfo.WritePtr += p.shmInfo.writeLen
	}
	p.shmInfo.writeLen = dataLen
	copy((p.shmData)[p.shmInfo.WritePtr:p.shmInfo.WritePtr+p.shmInfo.writeLen], data)
}

func (p *Publisher) Close() (err error) {
	p.Write([]byte("EOF"))
	p.IsClosed = true
	err = p.segment.DeleteShm()
	if err != nil {
		Logger.Info("DeleteShm err : ", err)
		return err
	}
	Logger.Info("Publisher Close")
	p.DoneSignal <- struct{}{}
	return nil
}
