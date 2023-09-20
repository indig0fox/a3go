package a3interface

/*
#include <stdlib.h>
#include <stdio.h>
#include <string.h>
*/
import "C"
import (
	"context"
	"fmt"
	"strings"
	"time"
	"unsafe"
)

// Config defines how calls to this extension will be handled
// it can be configured using method calls against it
var Config configStruct = configStruct{}

var activeContext context.Context
var closeReasonFnc context.CancelCauseFunc

func init() {
	Config.Init()
}

// called by Arma to get the version of the extension
//
//export RVExtensionVersion
func RVExtensionVersion(output *C.char, outputsize C.size_t) {
	result := Config.rvExtensionVersion
	replyToSyncArmaCall(result, output, outputsize)
}

type key int

var armaCtxKey key

type ArmaExtensionContext struct {
	SteamID           string
	FileSource        string
	MissionNameSource string
	ServerName        string
}

// passed just before all calls of exported functions
// in C/C++: void __stdcall RVExtensionContext(const char **args, int argsCnt)
//
//export RVExtensionContext
func RVExtensionContext(args **C.char, argsCnt C.int) {
	// convert args into context object
	// process the C vector into a Go slice
	var offset = unsafe.Sizeof(uintptr(0))
	var data []string
	for index := C.int(0); index < argsCnt; index++ {
		data = append(data, C.GoString(*args))
		args = (**C.char)(unsafe.Pointer(uintptr(unsafe.Pointer(args)) + offset))
	}

	callContextObj := ArmaExtensionContext{
		SteamID:           data[0],
		FileSource:        data[1],
		MissionNameSource: data[2],
		ServerName:        data[3],
	}

	// create a new context that will timeout after 800ms, so that we don't block Arma and we avoid a 301 EXECUTION_WARNING_TAKES_TOO_LONG
	activeContext, _ = context.WithTimeout(context.Background(), 800*time.Millisecond)
	activeContext = context.WithValue(activeContext, armaCtxKey, callContextObj)
	activeContext, closeReasonFnc = context.WithCancelCause(activeContext)
}

// called by Arma when in the format of: "extensionName" callExtension "command"
//
//export RVExtension
func RVExtension(output *C.char, outputsize C.size_t, input *C.char) {

	var command string = C.GoString(input)
	var commandSubstr string = strings.Split(command, "|")[0]
	var desiredCommand string
	var response string = "OK"

	if command == ":TIMESTAMP:" {
		response = getTimestamp()
		replyToSyncArmaCall(response, output, outputsize)
		return
	}

	// send default or timestamp reply immediately
	replyToSyncArmaCall(response, output, outputsize)

	// check if the callback channel is set for this command
	// first with the full command
	if _, ok := Config.rvExtensionFunctions[command]; !ok {
		// then with the substring
		if _, ok := Config.rvExtensionFunctions[commandSubstr]; !ok {
			// log an error if it isn't
			writeErrChan(command, fmt.Errorf("no channel set"))
			return
		}
		desiredCommand = commandSubstr
	} else {
		desiredCommand = command
	}

	// get channel
	fnc := Config.rvExtensionFunctions[desiredCommand]
	if fnc == nil {
		writeErrChan(command, fmt.Errorf("channel not set"))
		return
	}
	// call the function pointer
	go (*fnc)(activeContext.Value(armaCtxKey).(ArmaExtensionContext), closeReasonFnc, command)

	// on context close, send a reply to Arma
	go func() {
		<-activeContext.Done()

		// get reason provided by context closure
		reason := context.Cause(activeContext)
		if reason != nil {
			// if the reason is not nil, send it to Arma
			replyToSyncArmaCall(reason.Error(), output, outputsize)
		} else {
			// otherwise, send the default response
			replyToSyncArmaCall(response, output, outputsize)
		}
	}()

}

// called by Arma when in the format of: "extensionName" callExtension ["command", ["data"]]
//
//export RVExtensionArgs
func RVExtensionArgs(output *C.char, outputsize C.size_t, input *C.char, argv **C.char, argc C.int) {

	// get command as Go string
	command := C.GoString(input)
	// set default response
	response := fmt.Sprintf(`["Function: %s", "nb params: %d"]`, command, argc)

	replyToSyncArmaCall(response, output, outputsize)

	// get channel
	fnc := Config.rvExtensionArgsFunctions[command]
	if fnc == nil {
		writeErrChan(command, fmt.Errorf("function not set"))
		return
	}

	// now, we'll process the data
	// process the C vector into a Go slice
	var offset = unsafe.Sizeof(uintptr(0))
	var data []string
	for index := C.int(0); index < argc; index++ {
		data = append(data, C.GoString(*argv))
		argv = (**C.char)(unsafe.Pointer(uintptr(unsafe.Pointer(argv)) + offset))
	}

	// append timestamp in nanoseconds
	data = append(data, fmt.Sprintf("%d", time.Now().UnixNano()))

	// call the function pointer
	go (*fnc)(activeContext.Value(armaCtxKey).(ArmaExtensionContext), closeReasonFnc, command, data)

	// on context close, send a reply to Arma
	go func() {
		<-activeContext.Done()

		defaultResponse := fmt.Sprintf(`["Function: %s", "nb params: %d"]`, command, argc)
		// get reason provided by context closure
		reason := context.Cause(activeContext)
		if reason != nil {
			// if the reason is not nil, send it to Arma
			replyToSyncArmaCall(response, output, outputsize)
		} else {
			// otherwise, send the default response
			replyToSyncArmaCall(defaultResponse, output, outputsize)
		}
	}()

}

// replyToSyncArmaCall will respond to a synchronous extension call from Arma
// it returns a single string and any wait time will block Arma
func replyToSyncArmaCall(
	response string,
	output *C.char,
	outputsize C.size_t,
) {
	// Reply to a synchronous call from Arma with a string response
	result := C.CString(response)
	defer C.free(unsafe.Pointer(result))
	var size = C.strlen(result) + 1
	if size > outputsize {
		size = outputsize
	}
	C.memmove(unsafe.Pointer(output), unsafe.Pointer(result), size)
}

// writeErrChan will write an error to the error channel for a command
func writeErrChan(command string, err error) {
	if Config.errChan == nil {
		return
	}
	go func() {
		Config.errChan <- []string{command, err.Error()}
	}()
}

func getTimestamp() string {
	// get the current unix timestamp in nanoseconds
	return fmt.Sprintf("%d", time.Now().UTC().UnixNano())
	// return time.Now().Format("2006-01-02 15:04:05")
}
