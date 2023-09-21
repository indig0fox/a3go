package a3interface

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

func (r *RVExtensionRegistration) Register() {
	config.registrations = append(config.registrations, *r)
}

func NewRegistration(command string) *RVExtensionRegistration {
	return &RVExtensionRegistration{
		Command:         command,
		DefaultResponse: `["Command ` + command + ` called"]`,
	}
}
