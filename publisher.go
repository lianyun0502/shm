package shm

import (
	"os"
	"os/signal"
	"sync"
	"time"
	"unsafe"

	"github.com/gdygd/goshm/shmlinux"
	"github.com/go-co-op/gocron"
)

type Publisher struct {
	shmInfo     *ShmInfo
	shmData     []byte
	segmentInfo *shmlinux.Linuxshm
	segmentData *shmlinux.Linuxshm

	sysSignal  chan os.Signal
	DoneSignal chan struct{}

	IsClosed bool

	Scheduler *gocron.Scheduler
	msgIDLock sync.Mutex
	copyLock sync.RWMutex
}

func NewPublisher(skey int, shmSize int) *Publisher {
	segmentInfo, _ := NewSegment(skey, int(InfoSize))
	segmentData, _ := NewSegment(skey|0x6666, shmSize)

	sharedMemInfo := (*ShmInfo)(unsafe.Pointer(segmentInfo.Addr))
	sharedMemInfo.Size = uint(shmSize)
	p := (*byte)(unsafe.Pointer(segmentData.Addr))
	sharedMemData := unsafe.Slice(p, shmSize)

	publisher := &Publisher{
		shmInfo:     sharedMemInfo,
		shmData:     sharedMemData,
		segmentInfo: segmentInfo,
		segmentData: segmentData,
		sysSignal:   make(chan os.Signal),
		DoneSignal:  make(chan struct{}, 1),
		Scheduler:   gocron.NewScheduler(time.UTC),
	}

	signal.Notify(publisher.sysSignal, os.Interrupt)

	go func() {
		sig := <-publisher.sysSignal
		if sig == os.Interrupt {
			segmentInfo.DeleteShm()
			segmentData.DeleteShm()
		}
	}()
	publisher.Scheduler.Every(1).Day().At("00:00").Do(publisher.ResetMsgID)
	publisher.Scheduler.StartAsync()
	return publisher
}
func (p *Publisher) ResetMsgID() {
	p.msgIDLock.Lock()
	defer p.msgIDLock.Unlock()
	p.shmInfo.MsgID = 0
}

func (p *Publisher) IncreaseMsgID() {
	p.msgIDLock.Lock()
	defer p.msgIDLock.Unlock()
	p.shmInfo.MsgID++
}

func (p *Publisher) Write(data []byte) {
	if p.IsClosed {
		return
	}
	p.copyLock.Lock()
	var writePtr uint
	dataLen := uint(len(data))
	if p.shmInfo.WritePtr+p.shmInfo.writeLen+dataLen > uint(p.shmInfo.Size) {
		writePtr = 0
		p.ResetMsgID()
	} else {
		writePtr = p.shmInfo.WritePtr + p.shmInfo.writeLen
		p.IncreaseMsgID()
	}
	copy((p.shmData)[writePtr:writePtr+dataLen], data)
	p.copyLock.Unlock()
	p.shmInfo.writeLen = dataLen
	p.shmInfo.WritePtr = writePtr
	Logger.Debugf("MsgID : %d", p.shmInfo.MsgID)
}

func (p *Publisher) WriteStruct(obj any) {
	binData, err := GetBinary(obj)
	if err != nil {
		Logger.Info("GetBinary err : ", err)
		return
	}
	p.Write(binData)
}

func (p *Publisher) Close() (err error) {
	p.Write([]byte("EOF"))
	p.IsClosed = true
	err = p.segmentInfo.DeleteShm()
	if err != nil {
		Logger.Info("DeleteShm err : ", err)
		return err
	}
	err = p.segmentData.DeleteShm()
	if err != nil {
		Logger.Info("DeleteShm err : ", err)
		return err
	}
	Logger.Info("Publisher Close")
	p.DoneSignal <- struct{}{}
	return nil
}
