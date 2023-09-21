package a3interface

// ConfigStruct is the central configuration used by this library
type configStruct struct {

	// version is the value that will be returned when the extension is first called by Arma. This is a string value and is logged by the game engine to the RPT file
	version string

	// rvExtensionRegistrations is a collection of registrations that will be used to determine how to handle calls to the extension
	registrations []RVExtensionRegistration

	// errChan is the channel that errors will be sent to. the string slice will contain the command that caused the error and the error itself
	errChan chan []string
}

// Init method initializes the config struct
func (c *configStruct) init() {
	c.version = "No version set"
	c.registrations = make([]RVExtensionRegistration, 0)
}

func (c *configStruct) getRegistration(
	command string,
) *RVExtensionRegistration {
	for _, registration := range c.registrations {
		if registration.Command == command {
			return &registration
		}
	}
	return nil
}

// SetVersion sets the version string that will be returned when the extension is first called by Arma. This is a string value and is logged by the game engine to the RPT file
func SetVersion(version string) {
	config.version = version
}

// RegisterErrorChan triggered when an error occurs in the extension, this will send the error to the designated channel
func RegisterErrorChan(
	channel chan []string,
) {
	config.errChan = channel
}
