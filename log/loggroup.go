package log

import "log"

type LogGroup struct {
	loggers []*Logger
}

var gGroup *LogGroup

//var gLogger *Logger

func init() {
	gLogger, _ := New("debug", "", "", log.LstdFlags|log.Lshortfile)
	gGroup = new(LogGroup)
	gGroup.loggers = append(gGroup.loggers, gLogger)
}
func NewLogGroup(strLevel string, pathname string, isStdout bool, flag int) error {

	if pathname != "" {
		debug, err := New("debug", "/debug", pathname, flag)
		if err != nil {
			return err
		}
		info, err := New("info", "/info", pathname, flag)
		if err != nil {
			return err
		}
		erro, err := New("error", "/error", pathname, flag)
		if err != nil {
			return err
		}
		fatal, err := New("fatal", "/fatal", pathname, flag)
		if err != nil {
			return err
		}
		gGroup = &LogGroup{[]*Logger{debug, info, erro, fatal}}
	} else {
		gGroup = &LogGroup{}
	}

	if isStdout {
		gLogger, _ := New(strLevel, "", "", log.LstdFlags|log.Lshortfile)
		gGroup.loggers = append(gGroup.loggers, gLogger)
	}
	return nil
}

func Debug(format string, a ...interface{}) {
	for _, v := range gGroup.loggers {
		v.Debug(format, a...)
	}
}

func Info(format string, a ...interface{}) {
	for _, v := range gGroup.loggers {
		v.Info(format, a...)
	}
}

func _println(format string, a ...interface{}) {
	//str := fmt.Sprint(a)
	for _, v := range gGroup.loggers {
		v.Println(format, a...)
	}

}

func Error(format string, a ...interface{}) {
	//fmt.Println("Error  ", len(gGroup.loggers), gGroup.loggers)
	for _, v := range gGroup.loggers {
		v.Error(format, a...)
	}
}

func Fatal(format string, a ...interface{}) {
	for _, v := range gGroup.loggers {
		v.Fatal(format, a...)
	}
}

func Close() {
	for _, v := range gGroup.loggers {
		v.Close()
	}
}
