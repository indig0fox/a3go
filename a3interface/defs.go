package a3interface

/*
#include <stdlib.h>
#include <stdio.h>
#include <string.h>
*/
import "C"

// ConfigStruct is the central configuration used by this library
type configStruct struct {

	// rvExtensionVersion is the value that will be returned when the extension is first called by Arma. This is a string value and is logged by the game engine to the RPT file
	rvExtensionVersion string

	// rvExtensionFuncChannels stores a map of channels named by the commands provided with the SQF 'x' callExtension "command" format. The actual command name is checked against both a '|' delimited substring and the full command, then sent to the channel as a string if found.
	rvExtensionChannels map[string]chan string

	// rvExtensionArgsFuncs stores a map of channels named by the commands provided with the SQF 'x' callExtension ["command", ["data"]] format. The data array is sent to the channel as a slice of strings.
	rvExtensionArgsChannels map[string]chan []string

	// errChan is the channel that errors will be sent to. the string slice will contain the command that caused the error and the error itself
	errChan chan []string
}

// Init method initializes the config struct
func (c *configStruct) Init() {
	c.rvExtensionVersion = "No version set"
	c.rvExtensionChannels = make(map[string]chan string)
	c.rvExtensionArgsChannels = make(map[string]chan []string)
}

// SetVersion sets the version string that will be returned when the extension is first called by Arma. This is a string value and is logged by the game engine to the RPT file
func SetVersion(version string) {
	Config.rvExtensionVersion = version
}

// RegisterRvExtensionChannel triggered when SQF calls "x" callExtension "command", this will send the full command to the designated channel. The command is sent to the channel as a string.
func RegisterRvExtensionChannel(
	command string,
	channel chan string,
) {
	Config.rvExtensionChannels[command] = channel
}

// RegisterRvExtensionChannels triggered when SQF calls "x" callExtension "command", this will send the full command to the designated channel. The command is sent to the channel as a string.
func RegisterRvExtensionChannels(
	channels map[string]chan string,
) {
	// merge the new channels into the existing ones
	for k, v := range channels {
		Config.rvExtensionChannels[k] = v
	}
}

// RegisterRvExtensionArgsChannel triggered when SQF calls "x" callExtension ["command", ["data"]], this will send the data array to the designated channel based on the command. The data array is sent to the channel as a slice of strings and a timestamp at receipt is appended to the end of the slice (data[len(data)-1])
func RegisterRvExtensionArgsChannel(
	command string,
	channel chan []string,
) {
	Config.rvExtensionArgsChannels[command] = channel
}

// RegisterRvExtensionArgsChannels a map[string]chan []string, triggered when SQF calls "x" callExtension ["command", ["data"]], this will send the data array to the designated channel based on the command. The data array is sent to the channel as a slice of strings and a timestamp at receipt is appended to the end of the slice (data[len(data)-1])
func RegisterRvExtensionArgsChannels(
	channels map[string]chan []string,
) {
	// merge the new channels into the existing ones
	for k, v := range channels {
		Config.rvExtensionArgsChannels[k] = v
	}
}

// RegisterErrorChan triggered when an error occurs in the extension, this will send the error to the designated channel
func RegisterErrorChan(
	channel chan []string,
) {
	Config.errChan = channel
}

// TODO: add a way to unregister channels
// TODO: add a way to register a sync response for limited data, as subfunctions across channels cannot trigger
