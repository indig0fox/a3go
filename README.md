# a3go

## Description

Go library for Arma 3 extension development. Tested on Go 1.20.7.

All calls from Arma will receive an immediate basic response. To send data back to Arma from your extension, use the [Extension Callback handler](https://community.bistudio.com/wiki/Arma_3:_Mission_Event_Handlers#ExtensionCallback).

## Example

```sqf
addMissionEventHandler ["ExtensionCallback", {
  params ["_extension", "_function", "_args"];

  if (_extension isNotEqualTo "EXTENSION_NAME") exitWith {};

  _argsArr = parseSimpleArray _args;
  if (count _argsArr isEqualTo 0) exitWith {
    diag_log format["a3go: No arguments passed to extension. %1", _args];
  };

  switch (_function) do {
    case "log": {
      diag_log format["a3go: %1", _argsArr];
    };
    case "timeNow": {
      diag_log format["a3go: %1", _argsArr];
      // passed data will always be a string, and might be double double quoted!
      timeNow = _argsArr select 0;
    };
  };
}];
```

```go

import (
  "github.com/indig0fox/a3go/a3interface"
  "github.com/indig0fox/a3go/assemblyfinder"
)

// modulePath is the absolute path to the compiled DLL, which should be the addon folder
var modulePath string = assemblyfinder.GetModulePath()
// modulePathDir is the containing folder
var modulePathDir string = path.Dir(modulePath)

var EXTENSION_NAME = "EXTENSION_NAME"

var RVExtensionChannels = map[string]chan string {
  ":timeNow:" : make(chan string),
}
var RVExtensionArgsChannels = map[string]chan []string{
 ":LOG:JOIN:":    make(chan []string),
 ":LOG:LEAVE:":   make(chan []string),
 ":LOG:MISSION:": make(chan []string),
 ":LOG:WORLD:":   make(chan []string),
}
var a3ErrorChan = make(chan error)


func init() {

  a3interface.SetVersion(EXTENSION_VERSION)
  a3interface.RegisterRvExtensionArgsChannels(RVExtensionArgsChannels)


  go func() {
    for {
      select {
      case v := <-RVExtensionChannels[":timeNow:"]:
        // to call from A3: "EXTENSION_NAME" callExtension ":timeNow:";
        go writeTimeNow(v)
      case v := <-RVExtensionArgsChannels[":LOG:JOIN:"]:
        // to call from A3:"EXTENSION_NAME" callExtension [":LOG:JOIN:", ["Test1", "Test2", "Test3"]];
        go writeAttendance(v)
      case v := <-RVExtensionArgsChannels[":LOG:LEAVE:"]:
        // to call from A3: "EXTENSION_NAME" callExtension [":LOG:LEAVE:", ["Test1", "Test2", "Test3"]];
        go writeDisconnectEvent(v)
      case v := <-RVExtensionArgsChannels[":LOG:MISSION:"]:
        // to call from A3: "EXTENSION_NAME" callExtension [":LOG:MISSION:", ["Test1", "Test2", "Test3"]];
        go writeMissionEvent(v)
      case v := <-a3ErrorChan:
        log.Println(v.Error())
      }
    }
  }()
}

// writeTimeNow triggers the ExtensionCallback handler with:
// ["EXTENSION_NAME", "timeNow", ["1", "2021-01-01 00:00:00"]]
// the parameter is string because it's registered to the RVExtension handler
func writeTimeNow(id string) {
  // get current time
  t := time.Now()
  // format time
  timeNow := t.Format("2006-01-02 15:04:05")
  // send data back to Arma
  a3interface.WriteArmaCallback(EXTENSION_NAME, "timeNow", id, string(timeNow))
}

// writeAttendance triggers the ExtensionCallback handler with:
// ["EXTENSION_NAME", "log", ["Test1", "Test2", "Test3"]]
// the parameter is []string because it's registered to the RVExtensionArgs handler
func writeAttendance(args []string) {
  // ... do something with args
  a3interface.WriteArmaCallback(EXTENSION_NAME, "log", args...)
}
```

## Build

## Building using Docker

You will need Docker Engine installed and running. This can be done on Windows or on Linux. However, you will need to use Linux containers if you're on Windows (specified in Docker Desktop settings).

The below assumes you're running the commands from the `EXTENSION_NAME` directory, and you have (**at minimum**) the following files:
| File | Description |
| --- | --- |
| `EXTENSION_NAME.go` | The main extension file |

### COMPILING FOR WINDOWS

```bash
docker pull x1unix/go-mingw:1.20

# Compile x64 Windows DLL
docker run --rm -it -v ${PWD}:/go/work -w /go/work x1unix/go-mingw:1.20 go build -o dist/EXTENSION_NAME_x64.dll -buildmode=c-shared ./cmd/EXTENSION_NAME

# Compile x86 Windows DLL
docker run --rm -it -v ${PWD}:/go/work -w /go/work -e GOARCH=386 x1unix/go-mingw:1.20 go build -o dist/EXTENSION_NAME.dll -buildmode=c-shared ./cmd/EXTENSION_NAME

# Compile x64 Windows EXE
docker run --rm -it -v ${PWD}:/go/work -w /go/work x1unix/go-mingw:1.20 go build -o dist/EXTENSION_NAME_x64.exe ./cmd/EXTENSION_NAME
```

### COMPILING FOR LINUX

```bash
docker build -t indifox926/build-a3go:linux-so -f ./build/Dockerfile.build ./cmd

# Compile x64 Linux .so
docker run --rm -it -v ${PWD}:/app -e GOOS=linux -e GOARCH=amd64 -e CGO_ENABLED=1 -e CC=gcc indifox926/build-a3go:linux-so go build -o dist/EXTENSION_NAME_x64.so -linkshared ./cmd/EXTENSION_NAME

# Compile x86 Linux .so
docker run --rm -it -v ${PWD}:/app -e GOOS=linux -e GOARCH=386 -e CGO_ENABLED=1 -e CC=gcc indifox926/build-a3go:linux-so go build -o dist/EXTENSION_NAME.so -linkshared ./cmd/EXTENSION_NAME
```
