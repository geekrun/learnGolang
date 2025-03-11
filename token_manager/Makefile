# Android NDK 交叉编译配置
NDK_PATH := D:/data/ndk/android-ndk-r24
ANDROID_API := 30
ARCH := arm64
OUTPUT_DIR := ./bin
BIN_NAME := myapp_arm64
ANDROID_PATH := /data/local/tmp/$(BIN_NAME)

# 交叉编译工具链
CC := $(NDK_PATH)/toolchains/llvm/prebuilt/windows-x86_64/bin/aarch64-linux-android$(ANDROID_API)-clang
CXX := $(NDK_PATH)/toolchains/llvm/prebuilt/windows-x86_64/bin/aarch64-linux-android$(ANDROID_API)-clang++

# 目标任务：编译 Android ARM64 可执行文件
build:
	mkdir -p $(OUTPUT_DIR)
	GOOS=android GOARCH=arm64 CGO_ENABLED=1 CC=$(CC) CXX=$(CXX) go build -o $(OUTPUT_DIR)/$(BIN_NAME) main.go
	@echo "Build complete: $(OUTPUT_DIR)/$(BIN_NAME)"

# 编译为 Android 共享库（.so），可用于 JNI 或其他原生调用
build-shared:
	mkdir -p $(OUTPUT_DIR)
	GOOS=android GOARCH=arm64 CGO_ENABLED=1 CC="$(CC)" CXX="$(CXX)" go build -o $(OUTPUT_DIR)/lib$(BIN_NAME).so -buildmode=c-shared main.go
	@echo "Shared Library Build complete: $(OUTPUT_DIR)/lib$(BIN_NAME).so"

# 推送到 Android 设备
push: build
	adb push $(OUTPUT_DIR)/$(BIN_NAME) $(ANDROID_PATH)
	adb shell chmod +x $(ANDROID_PATH)
	@echo "Push complete: $(ANDROID_PATH)"

# 运行 Android 可执行文件
run: push
	adb shell $(ANDROID_PATH)

# 一键编译、推送并运行
all: build push run

# 清理
clean:
	rm -rf $(OUTPUT_DIR)
	@echo "Clean complete"
