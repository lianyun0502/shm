package shm

import (
	"os"
	"unsafe"

	"github.com/gdygd/goshm/shmlinux"
	"github.com/sirupsen/logrus"
)

type ShmInfo struct {
	WritePtr uint // 當下寫入位置
	writeLen uint // 當下寫入長度

	// Flag bool
	Size uint // shm大小
	MsgID    uint32 //當日第幾筆訊息
}

var info ShmInfo

const InfoSize = unsafe.Sizeof(info)

var Logger = &logrus.Logger{
	Out:          os.Stderr,
	Formatter:    &logrus.TextFormatter{DisableColors: true, TimestampFormat: "2006-01-02 15:04:05.000"},
	Hooks:        make(logrus.LevelHooks),
	Level:        logrus.InfoLevel,
	ExitFunc:     os.Exit,
	ReportCaller: false,
}

func NewSegment(skey int, size int) (*shmlinux.Linuxshm, error) {
	segment := shmlinux.NewLinuxShm()
	segment.InitShm(skey, size)
	err := segment.CreateShm()
	if err != nil {
		Logger.Warning("CreateShm err : ", err)
		return nil, err
	}
	err = segment.CreateShm()
	if err != nil {
		Logger.Warning("CreateShm err : ", err)
		return nil, err
	}
	return segment, nil
}
