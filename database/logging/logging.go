package logging

import (
	"log"
	"os"
	"io"
	"io/ioutil"
)

var (
	Trace   *log.Logger
	Debug   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

func CheckErr(err error) {
	if err != nil {
		Error.Println(err)
		panic(err)
	}
}

func LoggerInit(
	traceHandle io.Writer,
	debugHandle io.Writer,
	infoHandle io.Writer,
	warnHandle io.Writer,
	errorHandle io.Writer) {

	Trace = log.New(traceHandle,
		"[TRACE]: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Debug = log.New(debugHandle,
		"[DEBUG]: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Info = log.New(infoHandle,
		"[INFO]: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Warning = log.New(warnHandle,
		"[WARNING]: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Error = log.New(errorHandle,
		"[ERROR]: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}

func InitializeLogging(LogFile string) *os.File {
	logFile, err := os.OpenFile(LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open log file", os.Stdout, ":", err)
	}
	_, err = logFile.WriteString("\n\n=======START APPLICATION=======\n")
	CheckErr(err)

	logMultiOut := io.MultiWriter(logFile, os.Stdout)
	logFileOut := io.Writer(logFile)
	logStdErr := io.MultiWriter(logFile, os.Stderr)
	logDiscard := ioutil.Discard
	//logStdOut := os.Stdout

	LoggerInit(logDiscard, logFileOut, logMultiOut, logMultiOut, logStdErr)
	Trace.Println("Some discarded message")
	Debug.Println("Logger initiated")

	return logFile
}
