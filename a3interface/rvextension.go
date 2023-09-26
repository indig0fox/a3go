package a3interface

/*
#include <stdlib.h>
#include <stdio.h>
#include <string.h>
*/
import "C"
import (
	"fmt"
	"strings"
	"unsafe"
)

var activeContext *ArmaExtensionContext

// Config is the central configuration used by this library
var config *configStruct = new(configStruct)

func init() {
	// initialize the config struct
	config.init()
}

// called by Arma to get the version of the extension
//
//export RVExtensionVersion
func RVExtensionVersion(output *C.char, outputsize C.size_t) {
	replyToSyncArmaCall(config.version, output, outputsize)
}

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

	activeContext = &ArmaExtensionContext{
		SteamID:           data[0],
		FileSource:        data[1],
		MissionNameSource: data[2],
		ServerName:        data[3],
	}
}

// called by Arma when in the format of: "extensionName" callExtension "command"
//
//export RVExtension
func RVExtension(output *C.char, outputsize C.size_t, input *C.char) {

	var command string = C.GoString(input)
	var commandSubstr string = strings.Split(command, "|")[0]

	// look for registration
	registration := config.getRegistration(command)
	if registration == nil {
		registration = config.getRegistration(commandSubstr)
		if registration == nil {
			writeErrChan(command, fmt.Errorf("command not registered"))
			replyToSyncArmaCall(
				fmt.Sprintf(`["Command %s not registered!"]`, command),
				output, outputsize)
			return
		}
	}

	// get function pointer
	fnc := registration.Function
	if fnc == nil {
		writeErrChan(command, fmt.Errorf("function not set"))
		replyToSyncArmaCall(
			fmt.Sprintf(`["RVExtension function not set for command %s"]`, command),
			output, outputsize)
		return
	}

	// if RunInBackground is true for this registration, send default response
	// to Arma and run the function in the background
	// data can be sent back to arma using WriteArmaCallback
	if registration.RunInBackground {
		replyToSyncArmaCall(registration.DefaultResponse, output, outputsize)
		go func() {
			_, err := (fnc)(*activeContext, command)
			if err != nil {
				writeErrChan(command, err)
			}
		}()
		return
	}

	// otherwise, Arma is awaiting a reply
	response, err := (fnc)(*activeContext, command)
	if err != nil {
		writeErrChan(command, err)
		replyToSyncArmaCall(
			fmt.Sprintf(
				`[%q, %q]`,
				command,
				fmt.Sprintf("Error: %q", err.Error()),
			),
			output, outputsize)
		return
	} else {
		replyToSyncArmaCall(response, output, outputsize)
	}
}

// called by Arma when in the format of: "extensionName" callExtension ["command", ["data"]]
//
//export RVExtensionArgs
func RVExtensionArgs(output *C.char, outputsize C.size_t, input *C.char, argv **C.char, argc C.int) {

	// get command as Go string
	command := C.GoString(input)

	// look for registration
	registration := config.getRegistration(command)
	if registration == nil {
		writeErrChan(command, fmt.Errorf("command not registered"))
		replyToSyncArmaCall(
			fmt.Sprintf(`["Command %s not registered!"]`, command),
			output, outputsize)
		return
	}

	// if RunInBackground is true for this registration, send default response
	// to Arma and run the function in the background
	// data can be sent back to arma using WriteArmaCallback
	// do this so parsing arguments doesn't block Arma
	if registration.RunInBackground {
		replyToSyncArmaCall(registration.DefaultResponse, output, outputsize)
	}

	// now, we'll process the data
	// process the C vector into a Go slice
	var offset = unsafe.Sizeof(uintptr(0))
	var data []string
	for index := C.int(0); index < argc; index++ {
		data = append(data, C.GoString(*argv))
		argv = (**C.char)(unsafe.Pointer(uintptr(unsafe.Pointer(argv)) + offset))
	}

	// get function pointer
	fnc := registration.ArgsFunction
	if fnc == nil {
		writeErrChan(command, fmt.Errorf("function not set"))
		replyToSyncArmaCall(
			fmt.Sprintf(`["RVExtensionArgs function not set for command %s"]`, command),
			output, outputsize)
		return
	}

	// if running in background, launch the function in an asynchronous goroutine and return
	if registration.RunInBackground {
		go func() {
			_, err := (fnc)(*activeContext, command, data)
			if err != nil {
				writeErrChan(command, err)
			}
		}()
		return
	}

	// otherwise, Arma is awaiting a reply
	response, err := (fnc)(*activeContext, command, data)
	if err != nil {
		writeErrChan(command, err)
		replyToSyncArmaCall(
			fmt.Sprintf(
				`[%q, %q]`,
				command,
				fmt.Sprintf(
					"Error: %s",
					err.Error()),
			),
			output, outputsize)
		return
	} else {
		replyToSyncArmaCall(response, output, outputsize)
	}
}
