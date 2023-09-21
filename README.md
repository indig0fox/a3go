# a3go

## Description

Go library for Arma 3 extension development. Tested on Go 1.20.7.

All calls from Arma will receive an immediate basic response. To send data back to Arma from your extension, use the [Extension Callback handler](https://community.bistudio.com/wiki/Arma_3:_Mission_Event_Handlers#ExtensionCallback).

## Example

See [template](./template) for a working example of an addon and extension. You can follow the build steps below to compile the extension and addon. Use `./template` as your working directory when executing the commands!

## Building using Docker

You will need Docker Engine installed and running. This can be done on Windows or on Linux. However, you will need to use Linux containers if you're on Windows (specified in Docker Desktop settings). Building this way ensures that the CGo compiler will use the correct toolchain for the target platform -- there are issues when cross-compiling using the CGo compiler on Windows, and it will crash your game when trying to load the extension.

This assumes at least the following project structure beneath where you run the commands from:

```
|- build
|  |- Dockerfile.build
|- cmd
|  |- EXTENSION_NAME
|  |  |- main.go
|- dist
|- go.mod
```

*See [here](https://github.com/golang-standards/project-layout) for more information on Go project structure.*

### COMPILING FOR WINDOWS

```bash
docker pull x1unix/go-mingw:1.20

# Compile x64 Windows DLL
docker run --rm -it -v ${PWD}:/go/work -w /go/work x1unix/go-mingw:1.20 go build -o dist/EXTENSION_NAME_x64.dll -buildmode=c-shared -ldflags '-w -s' ./cmd/EXTENSION_NAME

# Compile x86 Windows DLL
docker run --rm -it -v ${PWD}:/go/work -w /go/work -e GOARCH=386 x1unix/go-mingw:1.20 go build -o dist/EXTENSION_NAME.dll -buildmode=c-shared -ldflags '-w -s' ./cmd/EXTENSION_NAME

# Compile x64 Windows EXE
docker run --rm -it -v ${PWD}:/go/work -w /go/work x1unix/go-mingw:1.20 go build -o dist/EXTENSION_NAME_x64.exe -ldflags '-w -s' ./cmd/EXTENSION_NAME
```

### COMPILING FOR LINUX

```bash
docker build -t indifox926/build-a3go:linux-so -f ./build/Dockerfile.build ./cmd

# Compile x64 Linux .so
docker run --rm -it -v ${PWD}:/app -e GOOS=linux -e GOARCH=amd64 -e CGO_ENABLED=1 -e CC=gcc indifox926/build-a3go:linux-so go build -o dist/EXTENSION_NAME_x64.so -linkshared -ldflags '-w -s' ./cmd/EXTENSION_NAME

# Compile x86 Linux .so
docker run --rm -it -v ${PWD}:/app -e GOOS=linux -e GOARCH=386 -e CGO_ENABLED=1 -e CC=gcc indifox926/build-a3go:linux-so go build -o dist/EXTENSION_NAME.so -linkshared -ldflags '-w -s' ./cmd/EXTENSION_NAME
```
