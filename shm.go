package shm

import (
	"github.com/sirupsen/logrus"
	"os"
	"unsafe"
)

type ShmMemInfo struct {
	WritePtr uint
	writeLen uint

	// Flag bool
	Size uint
}

var info ShmMemInfo

const InfoSize = unsafe.Sizeof(info)

var Logger = &logrus.Logger{
	Out:          os.Stderr,
	Formatter:    &logrus.TextFormatter{
		DisableColors: true, 
		TimestampFormat: "2006-01-02 15:04:05.000", 
		FullTimestamp: true,
	},
	Hooks:        make(logrus.LevelHooks),
	Level:        logrus.InfoLevel,
	ExitFunc:     os.Exit,
	ReportCaller: false,
}
