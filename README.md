# a3go

## Description

Go library for Arma 3 extension development. Tested on Go 1.20.7.

All calls from Arma will receive an immediate basic response. To send data back to Arma from your extension, use the [Extension Callback handler](https://community.bistudio.com/wiki/Arma_3:_Mission_Event_Handlers#ExtensionCallback).

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

```bash
docker pull x1unix/go-mingw:1.20

# Compile x64 Windows DLL
docker run --rm -it -v ${PWD}:/go/work -w /go/work x1unix/go-mingw:1.20 go build -o ./template/dist/EXTENSION_NAME_x64.dll -buildmode=c-shared -ldflags '-w -s' ./template/EXTENSION_NAME

# Compile x86 Windows DLL
docker run --rm -it -v ${PWD}:/go/work -w /go/work -e GOARCH=386 x1unix/go-mingw:1.20 go build -o ./template/dist/EXTENSION_NAME.dll -buildmode=c-shared -ldflags '-w -s' ./template/EXTENSION_NAME

# Compile x64 Windows EXE
docker run --rm -it -v ${PWD}:/go/work -w /go/work x1unix/go-mingw:1.20 go build -o ./template/dist/EXTENSION_NAME_x64.exe -ldflags '-w -s' ./template/EXTENSION_NAME
```

### EXTENSION: COMPILING FOR LINUX

Run this from the project root.

```bash
docker build -t indifox926/build-a3go:linux-so -f ./build/Dockerfile.build .

# Compile x64 Linux .so
docker run --rm -it -v ${PWD}:/app -e GOOS=linux -e GOARCH=amd64 -e CGO_ENABLED=1 -e CC=gcc indifox926/build-a3go:linux-so go build -o ./template/dist/EXTENSION_NAME_x64.so -linkshared -ldflags '-w -s' ./template/EXTENSION_NAME

# Compile x86 Linux .so
docker run --rm -it -v ${PWD}:/app -e GOOS=linux -e GOARCH=386 -e CGO_ENABLED=1 -e CC=gcc indifox926/build-a3go:linux-so go build -o ./template/dist/EXTENSION_NAME.so -linkshared -ldflags '-w -s' ./template/EXTENSION_NAME
```

### ADDON: COMPILE USING HEMTT

Download the [HEMTT binary](https://github.com/BrettMayson/HEMTT/releases/latest) and place it in [./template](./template), or wherever your .hemtt folder is located. The configuration inside will be read by the HEMTT exe and defines the build process.

```bash
cd ./template
./hemtt.exe build
```

## LICENSE

Arma Public License Share Alike (APL-SA) - See [LICENSE](./LICENSE) for more information.
