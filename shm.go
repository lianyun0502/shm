package shm

import (
	"github.com/sirupsen/logrus"
	"os"
	"unsafe"
)

type ShmMemInfo struct {
	WritePtr uint // 當下寫入位置
	writeLen uint // 當下寫入長度

	// Flag bool
	Size uint // shm大小
}

var info ShmMemInfo

const InfoSize = unsafe.Sizeof(info)

var Logger = &logrus.Logger{
	Out:          os.Stderr,
	Formatter:    &logrus.TextFormatter{DisableColors: true, TimestampFormat: "2006-01-02 15:04:05.000"},
	Hooks:        make(logrus.LevelHooks),
	Level:        logrus.InfoLevel,
	ExitFunc:     os.Exit,
	ReportCaller: false,
}