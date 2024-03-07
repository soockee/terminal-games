ROOT_DIR := $(dir $(realpath $(lastword $(MAKEFILE_LIST))))
ifdef OS
   RM = del /Q
   FixPath = $(subst /,\,$1)
else
   ifeq ($(shell uname), Linux)
      RM = rm -f
      FixPath = $1
   endif
endif

.PHONY: all build build-windows-amd64 build-linux-amd64 test clean run-windows-amd64

all: build

build: build-windows-amd64 build-linux-amd64 build-wasm

build-windows-amd64:
	podman compose run build-windows-amd64

build-linux-amd64:
	podman compose run build-linux-amd64

build-wasm:
	podman compose run build-wasm

run-windows-amd64:
	cmd /K %cd%/build/windows-amd64-tictacgoe.exe

test:
	podman compose run test

help:
	podman compose run help

clean:
	$(RM) $(call FixPath,$(ROOT_DIR)build/*)

