package shm

import (
	"os"
	"os/signal"
	"unsafe"
	"sync"
	"time"

	"github.com/gdygd/goshm/shmlinux"
	"github.com/go-co-op/gocron"
)

type Publisher struct {
	shmInfo *ShmMemInfo
	shmData []byte
	segment *shmlinux.Linuxshm

	sysSignal  chan os.Signal
	DoneSignal chan struct{}

	IsClosed bool
	msgID    uint32 //當日第幾筆訊息

	Scheduler *gocron.Scheduler
	msgIDLock sync.Mutex
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
		DoneSignal: make(chan struct{}),
		msgID:    0,
		Scheduler: gocron.NewScheduler(time.UTC),
	}

	signal.Notify(publisher.sysSignal, os.Interrupt)

	go func() {
		sig := <-publisher.sysSignal
		if sig == os.Interrupt {
			segmentInfo.DeleteShm()
		}
	}()
	publisher.Scheduler.Every(1).Day().At("00:00").Do(publisher.ResetMsgID)
	publisher.Scheduler.StartAsync()
	return publisher
}
func (p *Publisher) ResetMsgID() {
	p.msgIDLock.Lock()
	defer p.msgIDLock.Unlock()
	p.msgID = 0
}

func (p *Publisher) IncreaseMsgID() {
	p.msgIDLock.Lock()
	defer p.msgIDLock.Unlock()
	p.msgID++
}

func (p *Publisher) Write(data []byte) {
	if p.IsClosed {
		return
	}
	dataLen := uint(len(data))
	if p.shmInfo.WritePtr+p.shmInfo.writeLen+dataLen > uint(p.shmInfo.Size) {
		p.shmInfo.WritePtr = 0
		p.ResetMsgID()
	} else {
		p.shmInfo.WritePtr += p.shmInfo.writeLen
		p.IncreaseMsgID()
	}
	p.shmInfo.writeLen = dataLen
	copy((p.shmData)[p.shmInfo.WritePtr:p.shmInfo.WritePtr+p.shmInfo.writeLen], data)
	Logger.Debugf("MsgID : %d", p.msgID)
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
