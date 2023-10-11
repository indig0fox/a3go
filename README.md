# a3go

## Description

Go library for Arma 3 extension development. Tested on Go 1.20.7.

## Usage in your Go program

```go
go get github.com/indig0fox/a3go
```

## Features

- Synchronous and Asynchronous command handling
- Helper functions for parsing, including (nested) SQF arrays and hashmaps
- Send callbacks for use with the [Extension Callback handler](https://community.bistudio.com/wiki/Arma_3:_Mission_Event_Handlers#ExtensionCallback)
- Find the absolute path of the loaded DLL or SO file

To see an example of this library in use, see the [template](./template) folder and the [Attendance Tracker](https://github.com/indig0fox/Arma3-AttendanceTracker/) addon.

## a3interface API

### Registering Commands

The optimal method of registering a command to listen for is via a method chaining API.

> Note: See [a3interface.ArmaExtensionContext](###a3interface.ArmaExtensionContext) for more information on the context object passed to your function. See [a3interface Helper Functions](###a3interface-Helper-Functions) for more information on helper functions for removing extra escape quotations for strings, and for processing SQF arrays and hashmaps.

```go

/* REGISTER COMMAND HANDLER
Takes a single string parameter
Provide the exact text (case-sensitive) that will be used to call the command from Arma
This will be read from extension calls in the following orders:
For RVExtension:
 "extension" callExtension "commandText" (1)
 "extension" callExtension "commandText|arg1|arg2" (2)
For RVExtensionArgs:
 "extension" callExtension ["commandText", ["arg1", "arg2"]] (1) */
a3interface.NewRegistration(commandText).

/* SYNC
Takes a single boolean parameter
Run synchronously and respond to Arma with a string response or error
This is the default behavior */
  SetRunInBackground(false).
/* ASYNC
Run asynchronously and respond to Arma with a default response
Then start a goroutine to run the function in the background */
  SetRunInBackground(true).

/* DEFAULT RESPONSE
Takes a single string parameter
Configure the default response if RunInBackground is true
This will default to `["Command ` + command + ` called"]` */
  SetDefaultResponse(
    `["Received command ` + 
    commandText + 
    `, starting background process"]`
  ).

/* SET RVEXTENSION FUNCTION
Takes a single function parameter in format
  func(
    ctx a3interface.ArmaExtensionContext, data string,
  ) (string, error)
You can define the function inline or pass a function variable
This function will be called when the command is received from Arma and
the format `"extension" callExtension "commandText|arg1|arg2"` is used

In this example, we search the data string for a specific value and return an error if not found.

SYNCHRONOUS BEHAVIOR
If RunInBackground is false, then the function will be run synchronously and the return value will be sent to Arma as a string response. In this case (assuming parseSimpleArray is used on it), it would be:
["Found specific value in data"]
If 'specific value' was not in the data, however, then the return value sent to Arma in this case would be:
["commandText", "Error: Invalid data"]
This allows you to parse the array and detect the command text as well as use the SQF find command to search for the error string.
It's generally recommended to design your return data to Arma 3 in a stringified array format, as this allows you to send multiple values back to Arma in a single response and use parseSimpleArray to get your elements.

ASYNCHRONOUS BEHAVIOR
If RunInBackground is true, then the function will be run asynchronously and the default response will be sent to Arma immediately. In this case, it would be ["Received command commandText, starting background process"] because we set it above.
The function itself will then be called, as if its original defined scope, but with the parameters passed from Arma and in a non-blocking goroutine.
*/
  SetFunction(
    func(
      ctx a3interface.ArmaExtensionContext, data string,
    ) (string, error) {
      // If specific value not in data, return error
      if !strings.Contains(data, "specific value") {
        return "", errors.New("Invalid data")
      }
      // Do something with data
      return `["Found specific value in data"]`, nil
    },
  ).

/* SET RVEXTENSIONARGS FUNCTION
Takes a single function parameter in format
  func(
    ctx a3interface.ArmaExtensionContext, command string, data []string,
  ) (string, error)
You can define the function inline or pass a function variable
This function will be called when the command is received from Arma and
the format `"extension" callExtension ["commandText", ["arg1", "arg2"]]` is used

In this example, we search the data array for a specific value and return an error if not found.

SYNCHRONOUS BEHAVIOR
If RunInBackground is false, then the function will be run synchronously and the return value will be sent to Arma as a string response. In this case (assuming parseSimpleArray is used on it), it would be:
["Found specific value in data"]
If 'specific value' was not in the data, however, then the return value sent to Arma in this case would be:
["commandText", "Error: Invalid data"]
This allows you to parse the array and detect the command text as well as use the SQF find command to search for the error string.
It's generally recommended to design your return data to Arma 3 in a stringified array format, as this allows you to send multiple values back to Arma in a single response and use parseSimpleArray to get your elements.

ASYNCHRONOUS BEHAVIOR
If RunInBackground is true, then the function will be run asynchronously and the default response will be sent to Arma immediately. In this case, it would be ["Received command commandText, starting background process"] because we set it above.
The function itself will then be called, as if its original defined scope, but with the parameters passed from Arma and in a non-blocking goroutine.
*/
  SetArgsFunction(
    func(
      ctx a3interface.ArmaExtensionContext, command string, data []string,
    ) (string, error) {
      // preprocess the elements to remove double quotes
      // see below for more information on helper functions
      data = a3interface.RemoveEscapeQuotes(data)
      // If specific value not in data, return error
      for _, v := range data {
        if !strings.Contains(v, "specific value") {
          return "", errors.New("Invalid data")
        }
      }
      // Do something with data
      return `["Found specific value in data"]`, nil
    },
  ).

  /* REGISTER THE COMMAND
  This will register the command with the package so that calls with this command text will be handled.
  If you do not call this, then the command will not be registered and will not be handled.
  */
  Register()
```

### a3interface.ArmaExtensionContext

The context object passed to your function when a command is received from Arma contains four fields that provide context behind the call.

> See [A3 Wiki - callExtension](https://community.bistudio.com/wiki/callExtension) for more info.

```go
type ArmaExtensionContext struct {
  SteamID           string
  FileSource        string
  MissionNameSource string
  ServerName        string
}
```

### a3interface Helper Functions

#### RemoveEscapeQuotes

When strings are passed from Arma to Go, they are escaped with double quotes. This function will remove the double quotes from the string.

This function is important to use so you get your expected values when parsing the data.

```go
// definition
func RemoveEscapeQuotes(input string) string

// backticks indicate a raw string literal as you would process in Go
// `"my string"` -> `my string`
// `"[""my string""]"` -> `["my string"]`
// `"[""my string"", 34]"` -> `["my string", 34]`

// For RVExtensionArgs:
for _, v := range data {
  v = a3interface.RemoveEscapeQuotes(v)
}
```

#### ParseSQF

This function will take a raw string, expecting an SQF array or hashmap, and return an interface that you can check the indexes of and typecast to the appropriate type.

> It's important to note that this function includes a call to RemoveEscapeQuotes for the common use case of sending arrays and hashes to the extension, so do not preprocess the data with that function before passing it to this one!

```go
// definition
func ParseSQF(input string) interface{}

// !! all numerics (without quotes) should be parsed as float64

// backticks indicate a raw string literal as you would process in Go
// Arma: ["1", 2, 3] -> Extension: `"[""1"", 2, 3]" ->
// ParseSQF return: []interface{}{"1", 2, 3}
// r[0].(string) -> "1"
// r[1].(float64) -> 2.0
//
// Arma: [1, 2, [3, 4]] -> Extension: `"[1, 2, [3, 4]]"` ->
// ParseSQF return: `[]interface{}{1, 2, []interface{}{3, 4}}`
// r[1].(float64) -> 2.00
// r[2].([]interface{})[1].(float64) -> 4.00
// `"[""my string"", 34.2]"` -> `[]interface{}{"my string", 34.2}`
```

#### ParseSQFHashMap

This function will take an interface from ParseSQF, expecting an SQF HashMap, and return a map[string]interface{} with the keys and values. It will process nested values.

```go
// definition
func ParseSQFHashMap(input interface{}) (map[string]interface{}, error) 

/* 
backticks indicate a raw string literal as you would process in Go
Arma: [["key1", "value1"], ["keysExtra", ["myKey", "yeah!"], ["twokey", "oh no!"]]] -> 
Extension: `"[[""key1"", ""value1""], [""keysExtra"", [[""myKey"", ""yeah!""], [""twokey"", ""oh no!""]]]"` ->
ParseSQFHashMap return: map[string]interface{}{
  "key1": "value1",
  "keysExtra": map[string]interface{}{
    "myKey": "yeah!",
    "twokey": "oh no!",
  },
} 
*/

// example
func ReturnJSONFromHashMapArgs(
  ctx a3interface.ArmaExtensionContext,
  command string,
  args []string,
) (string, error) {

  JSONInterface, err := a3interface.ParseSQF(args[0])
  if err != nil {
    return "", err
  }
  JSONMapStringInterface, err := a3interface.ParseSQFHashMap(JSONInterface)
  if err != nil {
    return "", err
  }

  /* 
  JSONMapStringInterface["key1"].(string) -> "value1"
  extraKeys := JSONMapStringInterface["keysExtra"].(map[string]interface{})
  extraKeys["myKey"].(string) -> "yeah!"
  extraKeys["twokey"].(string) -> "oh no!"
  */

  JSONString, err := json.Marshal(JSONMapStringInterface)
  if err != nil {
    return "", err
  }

 return fmt.Sprintf(`%s`, JSONString), nil
 /* 
  returns
  "{""key1"":""value1",""key2":""value2""}"
  to Arma
 */
}
```

#### ToArmaHashMap

This function will take one of the following types and return a string that can be used in Arma. It is meant for HashMaps, but can also be used for arrays.

```go
// definition
func ToArmaHashMap(input interface{}) string

// example
func getVulcanTranslation(textToTranslate string) (map[string]interface{}, error) {
  // Make an HTTP request to the fun translation API and get our text in Vulcan
  url := url.URL{
  Scheme: "https",
  Host:   "api.funtranslations.com",
  Path:   "/translate/vulcan.json",
  RawQuery: url.Values{
    "text": []string{textToTranslate},
  }.Encode(),
  }
  resp, err := http.Get(url.String())
  if err != nil {
  return map[string]interface{}{}, err
  }
  defer resp.Body.Close()

  // Parse the JSON response and print the titles of the top headlines
  var data map[string]interface{}
  err = json.NewDecoder(resp.Body).Decode(&data)
  if err != nil {
  return map[string]interface{}{}, err
  }
  translationContent := data["contents"]
  if translationContent == nil {
  return map[string]interface{}{}, fmt.Errorf("No translation content found")
  }

  /* 
    "contents": {
      "translated": "Sos leralmin ik i ma ri thoroughly tizh-tor serving k' komihn? I talal ish-veh riolozhikaik heh duhik zherka wuh ek'kayik ornat.",
      "text": "May I say that I have not thoroughly enjoyed serving with humans? I find their illogical and foolish emotions a constant irritant.",
      "translation": "vulcan"
    } 
  */
  return translationContent.(map[string]interface{}), nil
  
}

// this would be registered as a command args function
a3interface.NewRegistration("translateToVulcan").
  SetArgsFunction(translateToVulcan).
  SetRunInBackground(false).
  Register()

func translateToVulcan(
  ctx a3interface.ArmaExtensionContext,
  command string,
  args []string,
) (string, error) {

  // get the translation
  translation, err := getVulcanTranslation(args[0])
  if err != nil {
    return "", err
  }

  // this map[string]interface{} (like a JSON object) will be parsed into a stringified HashMap format for SQF
  translationHashData := a3interface.ToArmaHashMap(translation)
  // return to Arma
  return translationHashData, nil
  // returns something like
  /* 
    "[[""contents"", [[""translated"", ""Sos leralmin ik i ma ri thoroughly tizh-tor serving k' komihn? I talal ish-veh riolozhikaik heh duhik zherka wuh ek'kayik ornat.""], [""text"", ""May I say that I have not thoroughly enjoyed serving with humans? I find their illogical and foolish emotions a constant irritant.""], [""translation"", ""vulcan""]]]]"
  */
}
```

To use the above in SQF:
  
```sqf
private _result = "extensionName" callExtension [
  "translateToVulcan", 
  ["Hello, how are you?"]
];

// NOTE: Args functions return an array with
// [<ourData>, 0, 0]
private _resultParsed = parseSimpleArray (_result select 0);
if (count _resultParsed isEqualTo 0) exitWith {};

private _resultHash = createHashMapFromArray _resultParsed;
if (isNil _resultHash) exitWith {};

hint formatText[
  "Language: %1\nTranslation: %2",
  _resultHash getOrDefault ["translation", "Unknown"],
  _resultHash getOrDefault ["translated", "Unknown"],
];
// this would show in a hint:
// Language: vulcan
// Translation: Tonk'peh,  uf nam-tor du?


diag_log formatText[
  "[Translator]: Translated English to $1.\nORIGINAL:%2\nTRANSLATED: %3",
  _resultHash getOrDefault ["translation", "Unknown"],
  _resultHash getOrDefault ["text", "Unknown"],
  _resultHash getOrDefault ["translated", "Unknown"],
];
// this would log
// [Translator]: Translated English to vulcan.
// ORIGINAL: Hello, how are you?
// TRANSLATED: Tonk'peh,  uf nam-tor du?
```

#### WriteArmaCallback

This function can be used to send a callback function to Arma from your extension. You can listen for these use the [Extension Callback handler](https://community.bistudio.com/wiki/Arma_3:_Mission_Event_Handlers#ExtensionCallback).

```go
// definition
func WriteArmaCallback(
extensionName string,
functionName string,
data ...string,
) (
err error,
)

// example Go function
// we would register this command using SetArgsFunction, since we're expecting
// "example_extension" callExtension ["example_callback", ["arg1", "arg2"]]
// we'll also assume here that SetRunInBackground(false) was used to
// demonstrate what Arma would receive (i.e. the return of this function
// instead of a default response)
func SendALogEntryAsACallback(
  ctx a3interface.ArmaExtensionContext,
  command string,
  args []string,
) (string, error) {
  err := a3interface.WriteArmaCallback(
    "example_extension",
    "LOG",
    "ERROR",
    "I didn't count high enough!",
  )
  if err != nil {
    return "", errors.New("Oh no!")
  }
  return `["Callback sent"]`, nil
}
```

The SQF would look like this:
  
```sqf
// add a callback listener to the mission
addMissionEventHandler ["ExtensionCallback", {
  params ["_extension", "_function", "_data"];
  if (_extension == "example_extension" && _function == "LOG") then {
    private arr = parseSimpleArray _data;
    // this will catch if the array is empty or couldn't be parsed in SQF
    if (count arr isEqualTo 0) exitWith {};

    // _extension = "example_extension"
    // _function = "LOG"
    // _data = "[""ERROR"", ""I didn't count high enough!""]"
    // _arr = ["ERROR", "I didn't count high enough!"]
    // _arr[0] -> "ERROR"
    // _arr[1] -> "I didn't count high enough!"

    ["[%1] %2", arr[0], arr[1]] call BIS_fnc_logFormat;
    // this will log "21:20:04 [ERROR] I didn't count high enough!" to the RPT
  };
}];

// make our extension call
private _immediateResult = "example_extension" callExtension ["example_callback", ["arg1", "arg2"]];
hint formatText ["%1", _immediateResult];
// _immediateResult -> "[""Callback sent""]"
// parseSimpleArray _immediateResult -> ["Callback sent"]

// if an error was returned, it would look like this:
// _immediateResult -> "[""example_callback"", ""Error: Oh no!""]"
// parseSimpleArray _immediateResult -> ["example_callback", "Error: Oh no!"]
```

## assemblyfinder API

This package is provided to locate the absolute path of the loaded DLL or SO file. This is useful for locating the addon directory (regardless of what it may be named) when you want to load a resource file from the same directory.

```go
// definition
func GetModulePath() string

// example
var dllFolder string = filepath.Dir(assemblyfinder.GetModulePath())
var logFilePathInAddonFolder = filepath.Join(
  dllFolder,
  "log.txt",
)
func GetAConfigFile() string {
  return filepath.Join(
    dllFolder,
    "config.json",
  )
  // returns something like C:\Program Files (x86)\Steam\steamapps\common\Arma 3\@example_addon\config.json
  // so long as the dll is in the @example_addon folder
  // NOTE: Extensions can also be loaded from the Arma 3 root directory, so you may need to check for that
}
```

## Template

See [template](./template) for a working example of an addon and extension. You can follow the build steps below to compile the extension and addon.

## Packaging your addon

Once you've built the extension then the addon using the build steps below, you will have a ./template/.hemttout/build folder containing `addons` and `dist` folders, as well as the license file.

Create your desired addon folder (i.e. `@example_addon`) and move the `build/addons` folder into it. Copy the extensions within `build/dist` into the `@example_addon` folder alongside the `addons` folder. Copy the LICENSE file under `build` into your `@example_addon` folder.

Add this `@example_addon` folder to your Arma launcher as a "Local Mod" and you should be good to launch.

*Note: [./template/addons/main/script_version.hpp](./template/addons/main/script_version.hpp) is used to provide versioning information to HEMTT. There is no CBA dependency in the example addon, but everything is included to implement that SQF Macro library if you wish.*

## Building using Docker

You will need Docker Engine installed and running. This can be done on Windows or on Linux. However, you will need to use Linux containers if you're on Windows (specified in Docker Desktop settings). Building this way ensures that the CGo compiler will use the correct toolchain for the target platform. Otherwise, there are issues when cross-compiling using the CGo or GCC compiler on Windows, and it will crash your game when trying to load the extension.

*See [here](https://github.com/golang-standards/project-layout) for more information on Go project structure.*

Build the extension first, then we can use HEMTT to build the addon and include the dll and so files.

### EXTENSION: COMPILING FOR WINDOWS

Run this from the project root.

```powershell
docker pull x1unix/go-mingw:1.20

# Compile x64 Windows DLL
docker run --rm -it -v ${PWD}:/go/work -w /go/work -e GOARCH=amd64 -e CGO_ENABLED=1 x1unix/go-mingw:1.20  go build -o ./template/dist/EXTENSION_NAME_x64.dll -buildmode=c-shared -ldflags '-w -s' ./template/EXTENSION_NAME

# Compile x86 Windows DLL
docker run --rm -it -v ${PWD}:/go/work -w /go/work -e GOARCH=386 -e CGO_ENABLED=1 x1unix/go-mingw:1.20 go build -o ./template/dist/EXTENSION_NAME.dll -buildmode=c-shared -ldflags '-w -s' ./template/EXTENSION_NAME

# Compile x64 Windows EXE
docker run --rm -it -v ${PWD}:/go/work -w /go/work -e GOARCH=amd64 -e CGO_ENABLED=1 x1unix/go-mingw:1.20 go build -o ./template/dist/EXTENSION_NAME_x64.exe -ldflags '-w -s' ./template/EXTENSION_NAME
```

### EXTENSION: COMPILING FOR LINUX

Run this from the project root.

```powershell
docker build -t indifox926/build-a3go:linux-so -f ./build/Dockerfile.build .

# Compile x64 Linux .so
docker run --rm -it -v ${PWD}:/app -e GOOS=linux -e GOARCH=amd64 -e CGO_ENABLED=1 -e CC=gcc indifox926/build-a3go:linux-so go build -o ./template/dist/EXTENSION_NAME_x64.so -linkshared -ldflags '-w -s' ./template/EXTENSION_NAME

# Compile x86 Linux .so
docker run --rm -it -v ${PWD}:/app -e GOOS=linux -e GOARCH=386 -e CGO_ENABLED=1 -e CC=gcc indifox926/build-a3go:linux-so go build -o ./template/dist/EXTENSION_NAME.so -linkshared -ldflags '-w -s' ./template/EXTENSION_NAME
```

### ADDON: COMPILE USING HEMTT

Download the [HEMTT binary](https://github.com/BrettMayson/HEMTT/releases/latest) and place it in [./template](./template), or wherever your .hemtt folder is located. The configuration inside will be read by the HEMTT exe and defines the build process.

```powershell
cd ./template
./hemtt.exe release
```

## LICENSE

Arma Public License Share Alike (APL-SA) - See [LICENSE](./LICENSE) for more information.
