package assemblyfinder

import (
	"path/filepath"
	"unsafe"
)

/*
#cgo windows LDFLAGS: -lpsapi
#cgo linux LDFLAGS: -ldl

#ifdef _WIN32
#define WIN32_LEAN_AND_MEAN
#include <windows.h>
#include <libloaderapi.h>
#include <stdlib.h>
#include <stdio.h>

char* GetModulePath() {
    HMODULE hModule = NULL;
    if (!GetModuleHandleExA(GET_MODULE_HANDLE_EX_FLAG_FROM_ADDRESS |
                           GET_MODULE_HANDLE_EX_FLAG_UNCHANGED_REFCOUNT,
                           (LPCTSTR)GetModulePath,
                           &hModule)) {
        return NULL;
    }

    DWORD size = MAX_PATH;
    char* buffer = NULL;
    while (1) {
        buffer = (char*)realloc(buffer, size);
        if (!buffer) {
            return NULL;
        }
        DWORD result = GetModuleFileNameA(hModule, buffer, size);
        if (result == 0) {
            free(buffer);
            return NULL;
        } else if (result < size) {
            return buffer;
        }
        size *= 2;
    }
}

#elif __linux__

#define _GNU_SOURCE
#include <dlfcn.h>
#include <stdlib.h>
#include <string.h>

char* GetModulePath() {
    Dl_info dl_info;
    dladdr((void*)GetModulePath, &dl_info);

    char* path = dl_info.dli_fname;
    char* result = (char*)malloc(strlen(path) + 1);
    strcpy(result, path);

    return result;
}

#endif
*/
import "C"

// GetModulePath returns the absolute path to the DLL or SO file this runtime was loaded from
func GetModulePath() string {
	modPath := C.GetModulePath()
	defer C.free(unsafe.Pointer(modPath))
	// fmt.Printf("Running on %s, module path: %s\n", runtime.GOOS, C.GoString(modPath))
	pathStr := C.GoString(modPath)
	absPath, err := filepath.Abs(pathStr)
	if err != nil {
		panic(err)
	}
	return absPath
}
