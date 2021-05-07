package logger

import (
	"log"
	"os"
)

var (
	outfile, _ = os.OpenFile("./logger/info.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
	LogFile    = log.New(outfile, "", 0)
)

func ForError(err error) {
	if err != nil {
		LogFile.Println(err)
		// LogFile.Fatalln(er)
	}
}

func LogCommandResult(str string){
	if str != "" {
		LogFile.Println(str)
		// LogFile.Fatalln(er)
	}
}
