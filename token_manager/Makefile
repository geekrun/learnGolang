NDK_PATH := D:/data/ndk/android-ndk-r24
ANDROID_API := 30
ARCH := arm64
OUTPUT_DIR := ./bin
BIN_NAME := tokencli
ANDROID_PATH := /data/local/tmp/$(BIN_NAME)

CC := $(NDK_PATH)/toolchains/llvm/prebuilt/windows-x86_64/bin/aarch64-linux-android$(ANDROID_API)-clang

build-windows:
	mkdir -p $(OUTPUT_DIR)
	GOOS=windows GOARCH=amd64 go build -o $(OUTPUT_DIR)/$(BIN_NAME).exe main.go
	@echo "Build complete: $(OUTPUT_DIR)/$(BIN_NAME).exe"

build-android:
	mkdir -p $(OUTPUT_DIR)
	GOOS=android GOARCH=arm64 CGO_ENABLED=1 CC=$(CC) go build -o $(OUTPUT_DIR)/$(BIN_NAME) main.go
	@echo "Build complete: $(OUTPUT_DIR)/$(BIN_NAME)"

push:build-android
	adb push $(OUTPUT_DIR)/$(BIN_NAME) $(ANDROID_PATH)
	adb shell chmod +x $(ANDROID_PATH)
	@echo "Push complete: $(ANDROID_PATH)"

run: push
	adb shell $(ANDROID_PATH)

all: build-android push run
