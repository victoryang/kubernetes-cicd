package logger

import (
	"log"
	"os"

	"github.com/toolkits/file"
)

var (
	LogFile   *os.File
	Logger    *log.Logger
	LogOutput bool
	DetailLog bool
)

func InitLoggerModule(logfile string, logOutput bool) error {
	var err error
	//得到日志文件所在的目录，如果目录不存在，那么创建出来
	dir := file.Dir(logfile)
	if !file.IsExist(dir) {
		//尝试创建，如果这也失败就不管了，后面直接fatal
		err := file.EnsureDir(dir)
		if err != nil {
			log.Fatal(err)
		}
	}
	LogFile, err := os.OpenFile(logfile, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	Logger = log.New(LogFile, "", log.Ldate|log.Ltime|log.Lshortfile)
	LogOutput = logOutput
	return nil
}

func Output(v ...interface{}) {
	if LogOutput {
		log.Println(v)
	}
}

func Debug(v ...interface{}) {
	Output("[debug]", v)
	Logger.Println("[debug]", v)
}

func Info(v ...interface{}) {
	Output("[info]", v)
	Logger.Println("[info]", v)
}

func Warning(v ...interface{}) {
	Output("[warning]", v)
	Logger.Println("[warning]", v)
}

func Error(v ...interface{}) {
	Output("[error]", v)
	Logger.Println("[error]", v)
}

func Fatal(v ...interface{}) {
	Output("[fatal]", v)
	Logger.Println("[Fatal]", v)
	os.Exit(1)
}

func Detail(v ...interface{}) {
	if DetailLog {
		Output("[info]", v)
		Logger.Println("[info]", v)
	}
}