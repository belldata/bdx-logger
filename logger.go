package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

var _ ILogger = &Logger{}

func itoa(i int, wid int) string {
	// Assemble decimal in reverse order.
	var b [20]byte
	bp := len(b) - 1
	for i >= 10 || wid > 1 {
		wid--
		q := i / 10
		b[bp] = byte('0' + i - q*10)
		bp--
		i = q
	}
	// i < 10
	b[bp] = byte('0' + i)
	return string(b[bp:])
}

const (
	levelDebug   = "DEBUG"
	levelInfo    = "INFO"
	levelWarning = "WARNING"
	levelError   = "ERROR"
	levelFatal   = "FATAL"
	levelALL     = "ALL"
	levelOFF     = "OFF"
)

// ILogger ログ出力インターフェイス
type ILogger interface {
	SetWriter(w io.Writer)
	SetLevel(level LogLevel)
	SetPrefix(prefix string)
	// SetLogColor(OnOff bool)
	SetLogFormat(format string)

	Debug(v ...interface{})
	Info(v ...interface{})
	Warning(v ...interface{})
	Error(v ...interface{})
	Fatal(v ...interface{})

	Debugf(format string, v ...interface{})
	Infof(format string, v ...interface{})
	Warningf(format string, v ...interface{})
	Errorf(format string, v ...interface{})
	Fatalf(format string, v ...interface{})
}

// LogLevel ログレベル列挙型
type LogLevel int

const (
	// Debug DEBUG
	Debug LogLevel = 1 << iota
	// Info INFO
	Info
	// Warning WARNING
	Warning
	// Error ERROR
	Error
	// Fatal FATAL
	Fatal
	// All ALL
	All
	// Off OFF
	Off
	unKnown
)

func (b LogLevel) String() string {
	switch b {
	case Debug:
		return levelDebug
	case Info:
		return levelInfo
	case Warning:
		return levelWarning
	case Error:
		return levelError
	case Fatal:
		return levelFatal
	case All:
		return levelALL
	case Off:
		return levelOFF
	default:
		return "UNKNOWN"
	}
}

// StringToLebel 文字列からBxLevelへ変換
func (b LogLevel) StringToLebel(s string) LogLevel {
	switch s {
	case levelDebug:
		return Debug
	case levelInfo:
		return Info
	case levelWarning:
		return Warning
	case levelError:
		return Error
	case levelFatal:
		return Fatal
	case levelALL:
		return All
	case levelOFF:
		return Off
	default:
		return unKnown
	}
}

// Logger ロガー構造体
type Logger struct {
	level         LogLevel
	writer        io.Writer
	printLevel    LogLevel
	log           *log.Logger
	mu            sync.Mutex
	prefix        string
	logFormat     string
	logColor      bool
	logColorStart string
	logColorEnd   string
}

// SetLevel ログレベルの変更
func (b *Logger) SetLevel(level LogLevel) {
	b.level = level
}

// SetWriter Writerをデフォルトの`os.Stdout`から変更する`
func (b *Logger) SetWriter(w io.Writer) {
	b.writer = w
}

// SetPrefix prefix文字の設定
func (b *Logger) SetPrefix(prefix string) {
	if prefix != "" {
		prefix = "[" + prefix + "]"
	}
	b.prefix = prefix
}

// SetLogFormat ログフォーマットを指定
// default -> $date$ $time$ $level$ $prefix$ $file$:$func$:$linenumber$: $message$
func (b *Logger) SetLogFormat(format string) {
	b.logFormat = format
}

func (b *Logger) getFuncName(pc uintptr) string {
	fn := runtime.FuncForPC(pc)
	return fn.Name()
}

func (b *Logger) getFileName(skip int) (string, string, uintptr) {
	var file string
	var line int
	var ok bool
	var counter uintptr
	counter, file, line, ok = runtime.Caller(skip)
	if !ok {
		file = "???"
		line = 0
	}
	short := file
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			short = file[i+1:]
			break
		}
	}
	return short, itoa(line, -1), counter
}

func (b *Logger) print(format string, v ...interface{}) {
	// $date$ $time$ $level$ $prefix$ $file$:$func$:$linenumber$: $message$
	now := time.Now()
	year, month, day := now.Date()
	date := itoa(year, 4) + "/" + itoa(int(month), 2) + "/" + itoa(day, 2)
	hour, min, sec := now.Clock()
	time := itoa(hour, 2) + ":" + itoa(min, 2) + ":" + itoa(sec, 2)
	fileName, line, counter := b.getFileName(4)
	funcName := b.getFuncName(counter)
	funcs := strings.Split(funcName, ".")
	funcName = funcs[len(funcs)-1]
	message := b.logFormat
	value := fmt.Sprintf(format, v...)

	message = strings.ReplaceAll(message, "$date$", date)
	message = strings.ReplaceAll(message, "$time$", time)
	message = strings.ReplaceAll(message, "$level$", b.logColorStart+"["+b.printLevel.String()+"]"+b.logColorEnd)
	message = strings.ReplaceAll(message, "$prefix$", b.prefix)
	message = strings.ReplaceAll(message, "$linenumber$", line)
	message = strings.ReplaceAll(message, "$file$", fileName)
	message = strings.ReplaceAll(message, "$func$", funcName)
	message = strings.ReplaceAll(message, "$message$", value)
	b.log.Println(message)
}

func (b *Logger) doPrint(v ...interface{}) {
	b.print("%s", v...)
}

func (b *Logger) doPrintf(format string, v ...interface{}) {
	b.print(format, v...)
}

func (b *Logger) levelCheck(level LogLevel) bool {
	switch level {
	case Debug:
		b.printLevel = level
		b.logColorStart = "\x1b[36m"
		b.logColorEnd = "\x1b[0m"
	case Info:
		b.printLevel = level
		b.logColorStart = "\x1b[32m"
		b.logColorEnd = "\x1b[0m"
	case Warning:
		b.printLevel = level
		b.logColorStart = "\x1b[33m"
		b.logColorEnd = "\x1b[0m"
	case Error, Fatal:
		b.printLevel = level
		b.logColorStart = "\x1b[31m"
		b.logColorEnd = "\x1b[0m"
	default:
		b.logColorStart = ""
		b.logColorEnd = ""
	}

	switch b.level {
	case Debug:
		return level&(Debug|Info|Warning|Error|Fatal) != 0
	case Info:
		return level&(Info|Warning|Error|Fatal) != 0
	case Warning:
		return level&(Warning|Error|Fatal) != 0
	case Error:
		return level&(Error|Fatal) != 0
	case Fatal:
		return level&(Fatal) != 0
	case All:
		return true
	case Off:
		return false
	default:
		return false
	}
}

func (b *Logger) debug() bool {
	return b.levelCheck(Debug)
}
func (b *Logger) info() bool {
	return b.levelCheck(Info)
}
func (b *Logger) warning() bool {
	return b.levelCheck(Warning)
}
func (b *Logger) error() bool {
	return b.levelCheck(Error)
}
func (b *Logger) fatal() bool {
	return b.levelCheck(Fatal)
}

// Debug is Same as log.Print
func (b *Logger) Debug(v ...interface{}) {
	if b.debug() {
		b.doPrint(v...)
	}
}

// Info is Same as log.Print
func (b *Logger) Info(v ...interface{}) {
	if b.info() {
		b.doPrint(v...)
	}
}

// Warning is Same as log.Print
func (b *Logger) Warning(v ...interface{}) {
	if b.warning() {
		b.doPrint(v...)
	}
}

// Error is Same as log.Print
func (b *Logger) Error(v ...interface{}) {
	if b.error() {
		b.doPrint(v...)
	}
}

// Fatal is Same as log.Print
func (b *Logger) Fatal(v ...interface{}) {
	if b.fatal() {
		b.doPrint(v...)
	}
}

// Debugf is Same as log.Printf
func (b *Logger) Debugf(format string, v ...interface{}) {
	if b.debug() {
		b.doPrintf(format, v...)
	}
}

// Infof is Same as log.Printf
func (b *Logger) Infof(format string, v ...interface{}) {
	if b.info() {
		b.doPrintf(format, v...)
	}
}

// Warningf is Same as log.Printf
func (b *Logger) Warningf(format string, v ...interface{}) {
	if b.warning() {
		b.doPrintf(format, v...)
	}
}

// Errorf is Same as log.Printf
func (b *Logger) Errorf(format string, v ...interface{}) {
	if b.error() {
		b.doPrintf(format, v...)
	}
}

// Fatalf is Same as log.Printf
func (b *Logger) Fatalf(format string, v ...interface{}) {
	if b.fatal() {
		b.doPrintf(format, v...)
	}
}

// New Loggerインスタンス生成.
// level ... array first only
func New(prefix string, level ...LogLevel) ILogger {
	l := Info
	if len(level) > 0 {
		l = level[0]
	}
	b := &Logger{
		level:     l,
		writer:    os.Stdout,
		logColor:  false,
		logFormat: "$date$ $time$ $level$ $prefix$ $file$:$func$:$linenumber$: $message$",
	}
	b.SetPrefix(prefix)
	b.log = log.New(b.writer, "", 0)
	return b
}
