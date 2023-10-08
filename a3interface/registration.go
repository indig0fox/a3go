package a3interface

import "fmt"

type RVExtensionRegistration struct {
	// Command When this command is sent as the first element of a pipe-delimited string in RVExtension or as the command element in RVExtensionArgs, this registration will be referenced. i.e. "command|data" or ["command", ["data"]]. This is case sensitive & will call Function or ArgsFunction based on the call type used.
	Command string
	// DefaultResponse will be returned to Arma if Async is true. If Async is false, this value is ignored
	DefaultResponse string
	// RunInBackground determines whether or not the library will respond instantly to Arma with DefaultResponse or wait for a return from the function
	RunInBackground bool
	// Function is a function pointer that will be called in the "extension" callExtension "command|data" format
	Function func(
		ctx ArmaExtensionContext,
		data string) (string, error)

	// ArgsFunction is a function pointer that will be called in the "extension" callExtension ["command", ["data"]] format
	ArgsFunction func(
		ctx ArmaExtensionContext,
		command string,
		args []string) (string, error)
}

func NewRegistration(command string) *RVExtensionRegistration {
	return &RVExtensionRegistration{
		Command:         command,
		DefaultResponse: `["Command ` + command + ` called"]`,
	}
}

func (r *RVExtensionRegistration) SetDefaultResponse(response string) *RVExtensionRegistration {
	r.DefaultResponse = response
	return r
}

// SetRunInBackground determines whether or not the library will respond instantly to Arma with DefaultResponse and run the function in a goroutine (true), or wait for a return from the function (false)
func (r *RVExtensionRegistration) SetRunInBackground(runInBackground bool) *RVExtensionRegistration {
	r.RunInBackground = runInBackground
	return r
}

// SetFunction sets the function pointer that will be called in the "extension" callExtension "command|data" format
func (r *RVExtensionRegistration) SetFunction(
	fnc func(ctx ArmaExtensionContext, data string) (string, error),
) *RVExtensionRegistration {
	r.Function = fnc
	return r
}

// SetArgsFunction sets the function pointer that will be called in the "extension" callExtension ["command", ["data"]] format
func (r *RVExtensionRegistration) SetArgsFunction(
	fnc func(ctx ArmaExtensionContext, command string, args []string) (string, error),
) *RVExtensionRegistration {
	r.ArgsFunction = fnc
	return r
}

// Register adds this registration to the list of registrations that will be used to determine how to handle calls to the extension
func (r *RVExtensionRegistration) Register() {
	config.registrations = append(config.registrations, *r)
	fmt.Printf("Registered command: %s\n", r.Command)
	fmt.Printf("%+v\n", config.registrations)
}
