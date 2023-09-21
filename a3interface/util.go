package a3interface

/*
#include <stdlib.h>
#include <stdio.h>
#include <string.h>
*/
import "C"
import (
	"fmt"
	"time"
	"unsafe"
)

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
	if config.errChan == nil {
		return
	}
	go func() {
		config.errChan <- []string{command, err.Error()}
	}()
}

func getTimestamp() string {
	// get the current unix timestamp in nanoseconds
	return fmt.Sprintf("%d", time.Now().UTC().UnixNano())
	// return time.Now().Format("2006-01-02 15:04:05")
}
