#!/bin/bash

source ../buildScript/init/env_ndk.sh

# Set up Go cross-compilation for Android
export GOOS=android
export CGO_ENABLED=1

# Build directory
BUILD_DIR=".build"
OUTPUT_DIR="../app/src/main/jniLibs"

rm -rf $BUILD_DIR/lumine
mkdir -p $BUILD_DIR/lumine

# Android architectures
declare -A ARCHS=(
    ["arm64"]="arm64"
    ["arm"]="arm"
    ["x86"]="386"
    ["x86_64"]="amd64"
)

declare -A NDK_ARCHS=(
    ["arm64"]="aarch64-linux-android"
    ["arm"]="armv7a-linux-androideabi"
    ["x86"]="i686-linux-android"
    ["x86_64"]="x86_64-linux-android"
)

API_LEVEL=21

# Detect OS for NDK path
if [[ "$OSTYPE" == "darwin"* ]]; then
    NDK_HOST="darwin-x86_64"
else
    NDK_HOST="linux-x86_64"
fi

for arch in "${!ARCHS[@]}"; do
    echo "Building lumine for $arch..."
    
    export GOARCH=${ARCHS[$arch]}
    
    # Set up NDK toolchain
    NDK_ARCH=${NDK_ARCHS[$arch]}
    export CC="$NDK/toolchains/llvm/prebuilt/$NDK_HOST/bin/${NDK_ARCH}${API_LEVEL}-clang"
    export CXX="$NDK/toolchains/llvm/prebuilt/$NDK_HOST/bin/${NDK_ARCH}${API_LEVEL}-clang++"
    
    # Special case for arm (armv7a)
    if [ "$arch" = "arm" ]; then
        export GOARM=7
    fi
    
    # Build the binary
    mkdir -p "$BUILD_DIR/lumine/$arch"
    go build -trimpath -ldflags="-s -w" -o "$BUILD_DIR/lumine/$arch/liblumine.so" ./lumine_cmd || exit 1
    
    # Copy to jniLibs
    case $arch in
        arm64)
            mkdir -p "$OUTPUT_DIR/arm64-v8a"
            cp "$BUILD_DIR/lumine/$arch/liblumine.so" "$OUTPUT_DIR/arm64-v8a/liblumine.so"
            chmod +x "$OUTPUT_DIR/arm64-v8a/liblumine.so"
            ;;
        arm)
            mkdir -p "$OUTPUT_DIR/armeabi-v7a"
            cp "$BUILD_DIR/lumine/$arch/liblumine.so" "$OUTPUT_DIR/armeabi-v7a/liblumine.so"
            chmod +x "$OUTPUT_DIR/armeabi-v7a/liblumine.so"
            ;;
        x86)
            mkdir -p "$OUTPUT_DIR/x86"
            cp "$BUILD_DIR/lumine/$arch/liblumine.so" "$OUTPUT_DIR/x86/liblumine.so"
            chmod +x "$OUTPUT_DIR/x86/liblumine.so"
            ;;
        x86_64)
            mkdir -p "$OUTPUT_DIR/x86_64"
            cp "$BUILD_DIR/lumine/$arch/liblumine.so" "$OUTPUT_DIR/x86_64/liblumine.so"
            chmod +x "$OUTPUT_DIR/x86_64/liblumine.so"
            ;;
    esac
    
    echo "Built liblumine.so for $arch"
done

echo "All architectures built successfully"
echo "Lumine plugin installed to $OUTPUT_DIR"

