package shm

import (
	"fmt"
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
}

func NewPublisher(skey int, shmSize int) *Publisher {
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
	sharedMemInfo := (*ShmMemInfo)(unsafe.Pointer(segmentInfo.Addr))
	sharedMemInfo.Size = uint(shmSize)
	sharedMemInfo.Flag = false
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
	dataLen := uint(len(data))
	if p.shmInfo.WritePtr+p.shmInfo.writeLen+dataLen > uint(p.shmInfo.Size) {
		p.shmInfo.WritePtr = 0
	} else {
		p.shmInfo.WritePtr += p.shmInfo.writeLen
	}
	p.shmInfo.writeLen = dataLen
	copy((p.shmData)[p.shmInfo.WritePtr:p.shmInfo.WritePtr+p.shmInfo.writeLen], data)

	if p.shmInfo.Flag {
		p.shmInfo.Flag = false
	} else {
		p.shmInfo.Flag = true
	}
}

func (p *Publisher) Close() (err error) {
	p.Write([]byte("EOF"))
	err = p.segment.DeleteShm()
	if err != nil {
		fmt.Println("DeleteShm err : ", err)
		return err
	}
	fmt.Println("Publisher Close")
	p.DoneSignal <- struct{}{}
	return nil
}
