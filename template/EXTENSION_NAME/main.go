package main

import (
	"path/filepath"

	"github.com/indig0fox/a3go/a3interface"
	"github.com/indig0fox/a3go/assemblyfinder"
)

var a3ErrorChannel chan []string = make(chan []string)

// we can use the assemblyfinder library to get the absolute path to our DLL. this is useful for finding files relative to the DLL, like sidecar config files in the addon's root directory.
var dllAbsPath string = assemblyfinder.GetModulePath()
var addonDirectory string = filepath.Dir(dllAbsPath)

func init() {
	a3interface.SetVersion("1.0.0")
	a3interface.RegisterErrorChan(a3ErrorChannel)

	// SYNCHRONOUS EXAMPLE
	// calling "test" as a command will expect a string response to be fed back to Arma.
	// we don't want to do anything long-running here as it will block Arma. the default "RunInBackground" setting is false, so if we don't configure it, Arma will be waiting for our function returns.
	testCommand := a3interface.NewRegistration("test")
	// give it something to do when called using "EXTENSION_NAME" callExtension "test|test1|test2"
	testCommand.Function = receiveTestCommand
	// give it something to do when called using "EXTENSION_NAME" callExtension ["test", ["test1", "test2"]]
	testCommand.ArgsFunction = receiveTestCommandArgs
	// NOTE: providing no default response will cause the library to return ["Command test called"] to Arma
	testCommand.Register()

	// ASYNCHRONOUS EXAMPLE
	// calling testAsync as a command will instead return a default response to Arma and run the function in the background. we can use the a3interface.WriteArmaCallback function to send data back to Arma and the SQF `addMissionEventHandler ["ExtensionCallback", {}]` function to receive it.
	testAsyncCommand := a3interface.NewRegistration("testAsync")
	testAsyncCommand.RunInBackground = true
	testAsyncCommand.DefaultResponse = `["testAsync called"]`
	testAsyncCommand.Function = receiveTestCommand
	testAsyncCommand.ArgsFunction = receiveTestCommandArgs
	testAsyncCommand.Register()

	// CHAIN SYNTAX EXAMPLE
	// this command will log the caller context to a sqlite database
	// here we use the API chain syntax to configure the registration
	a3interface.NewRegistration("saveMyCall").
		SetDefaultResponse(`["saveMyCall called"]`).
		SetRunInBackground(true).
		SetFunction(saveCaller).
		SetArgsFunction(saveCallerArgs).
		Register()
}

// NOTE: This main function must exist for building the DLL, but isn't exposed and won't be called by Arma. You could build an exe or binary using this library for testing or other purposes and, upon running it, this main function would be called.
func main() {}
