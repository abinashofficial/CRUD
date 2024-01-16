package log

import (
	"crud/tapcontext"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

const (
	//INFO level 1
	INFO = iota
	//HTTP level 2
	HTTP
	//ERROR level 3
	ERROR
	//TRACE level 4
	TRACE
	//WARNING level 5
	WARNING
)

var (
	setLevel = WARNING
	trace    *log.Logger
	info     *log.Logger
	warning  *log.Logger
	httplog  *log.Logger
	errorlog *log.Logger
)

const (
	clusterType      = "CLUSTER_TYPE"
	clusterTypeLocal = "local"
	clusterTypeDev   = "dev1"
)

// FieldsMap map of key value pair to log
type FieldsMap map[string]interface{}

func init() {
	logInit(os.Stdout,
		os.Stdout,
		os.Stdout,
		os.Stdout,
		os.Stderr)
}

func logInit(traceHandle, infoHandle, warningHandle, httpHandle, errorHandle io.Writer) {

	flagWithClusterType := log.LUTC | log.LstdFlags | log.Lshortfile
	flagWithoutClusterType := log.LUTC | log.LstdFlags

	var flag int

	if os.Getenv(clusterType) == clusterTypeLocal || os.Getenv(clusterType) == clusterTypeDev {
		flag = flagWithClusterType
	} else {
		flag = flagWithoutClusterType
	}

	trace = log.New(traceHandle, "TRACE|", flag)

	info = log.New(infoHandle, "INFO|", flag)

	warning = log.New(warningHandle, "WARNING|", flag)

	httplog = log.New(httpHandle, "HTTP|", flag)

	errorlog = log.New(errorHandle, "ERROR|", flagWithClusterType)
}

func doLog(cLog *log.Logger, level, callDepth int, v ...interface{}) {
	if level <= setLevel {
		if level == ERROR {
			cLog.SetOutput(os.Stderr)
			cLog.SetFlags(log.Llongfile)
		}
		//cLog.SetOutput(os.Stdout)
		cLog.Output(callDepth, fmt.Sprintln(v...))
	}
}

func generatePrefix(ctx tapcontext.TContext) string {
	// TODO: Add departmentName once that is implemented fully
	return strings.Join([]string{ctx.UserEmail}, ":")
}

func generateTrackingIDs(ctx tapcontext.TContext) (retString string) {
	requestID := ctx.RequestID

	if requestID != "" {
		retString = "requestId=" + requestID
	}

	return
}

// Trace system gives facility to helps you isolate your system problems by monitoring selected events Ex. entry and exit
func traceLog(v ...interface{}) {
	doLog(trace, TRACE, 4, v...)
}

// Info dedicated for logging valuable information
func infoLog(v ...interface{}) {
	doLog(info, INFO, 4, v...)
}

// Warning for critical error
func warningLog(v ...interface{}) {
	doLog(warning, WARNING, 4, v...)
}

// Error logging error
func errorLog(v ...interface{}) {
	doLog(errorlog, ERROR, 4, v...)
}

func fatalLog(v ...interface{}) {
	doLog(errorlog, ERROR, 4, v...)
	os.Exit(1)
}

// HTTPLog prints the log in the following format:
//
// If any of the value is irrelevant then two consecutive PIPEs are printed:
// HTTP|TIMESTAMP|TenenatName:DelaerID:userName|ServerIP:PORT|RequestMethod|RequestURL|ResponseStatusCode|ResponseWeight|UserAgent|Duration
func HTTPLog(logMessage string) {
	doLog(httplog, HTTP, 6, logMessage)
}

// GenericTrace generates trace log (following standard Tekion log spec)
func GenericTrace(ctx tapcontext.TContext, traceMessage string, data ...FieldsMap) {
	var fields FieldsMap
	if len(data) > 0 {
		fields = data[0]
	}
	if os.Getenv("TEK_SERVICE_TRACE") == "true" {
		prefix := generatePrefix(ctx)
		trackingIDs := generateTrackingIDs(ctx)
		msg := fmt.Sprintf("|%s|%s|",
			prefix,
			trackingIDs)
		if fields != nil && len(fields) > 0 {
			fieldsBytes, _ := json.Marshal(fields)
			fieldsString := string(fieldsBytes)
			traceLog(msg, traceMessage, "|", fieldsString)
		} else {
			traceLog(msg, traceMessage)
		}
	}
}

// GenericInfo generates info log (following standard Tekion log spec)
func GenericInfo(ctx tapcontext.TContext, infoMessage string, data ...FieldsMap) {
	var fields FieldsMap
	if len(data) > 0 {
		fields = data[0]
	}
	prefix := generatePrefix(ctx)
	trackingIDs := generateTrackingIDs(ctx)
	fieldsBytes, _ := json.Marshal(fields)
	fieldsString := string(fieldsBytes)
	msg := fmt.Sprintf("|%s|%s|",
		prefix,
		trackingIDs)
	if fields != nil && len(fields) > 0 {
		infoLog(msg, infoMessage, "|", fieldsString)
	} else {
		infoLog(msg, infoMessage)
	}

}

// GenericWarning generates warning log (following standard Tekion log spec)
func GenericWarning(ctx tapcontext.TContext, warnMessage string, data ...FieldsMap) {
	var fields FieldsMap
	if len(data) > 0 {
		fields = data[0]
	}
	if os.Getenv("TEK_SERVICE_WARN") == "true" {
		prefix := generatePrefix(ctx)
		trackingIDs := generateTrackingIDs(ctx)
		msg := fmt.Sprintf("|%s|%s|",
			prefix,
			trackingIDs)
		if fields != nil && len(fields) > 0 {
			fieldsBytes, _ := json.Marshal(fields)
			fieldsString := string(fieldsBytes)
			warningLog(msg, warnMessage, "|", fieldsString)
		} else {
			warningLog(msg, warnMessage)
		}
	}
}

// GenericError generates error log (following standard Tekion log spec)
func GenericError(ctx tapcontext.TContext, e error, data ...FieldsMap) {
	var fields FieldsMap
	if len(data) > 0 {
		fields = data[0]
	}
	prefix := generatePrefix(ctx)
	trackingIDs := generateTrackingIDs(ctx)
	msg := ""
	if e != nil {
		msg = fmt.Sprintf("|%s|%s|%s", prefix, trackingIDs, e.Error())
	} else {
		msg = fmt.Sprintf("|%s|%s", prefix, trackingIDs)
	}

	if fields != nil && len(fields) > 0 {
		fieldsBytes, _ := json.Marshal(fields)
		fieldsString := string(fieldsBytes)
		errorLog(msg, "|", fieldsString)
	} else {
		errorLog(msg)
	}
}

func FatalLog(ctx tapcontext.TContext, e error, data ...FieldsMap) {
	var fields FieldsMap
	if len(data) > 0 {
		fields = data[0]
	}
	prefix := generatePrefix(ctx)
	trackingIDs := generateTrackingIDs(ctx)

	msg := ""
	if e != nil {
		msg = fmt.Sprintf("|%s|%s|%s", prefix, trackingIDs, e.Error())
	} else {
		msg = fmt.Sprintf("|%s|%s", prefix, trackingIDs)
	}

	if fields != nil && len(fields) > 0 {
		fieldsBytes, _ := json.Marshal(fields)
		fieldsString := string(fieldsBytes)
		fatalLog(msg, "|", fieldsString)
	} else {
		fatalLog(msg)
	}
}
