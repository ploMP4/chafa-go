CHAFAVERSION = 1.16
PREFIX       = $(CURDIR)/build
LIBDIR       = libs

# Define all target architectures
LINUX_TARGETS = $(LIBDIR)/linux_amd64/libchafa.so $(LIBDIR)/linux_arm64/libchafa.so $(LIBDIR)/linux_386/libchafa.so
DARWIN_TARGETS = $(LIBDIR)/darwin_amd64/libchafa.dylib $(LIBDIR)/darwin_arm64/libchafa.dylib
WINDOWS_TARGETS = $(LIBDIR)/windows_amd64/libchafa.dll

.PHONY: all clean linux darwin windows

all: linux darwin windows

linux: $(LINUX_TARGETS)

darwin: $(DARWIN_TARGETS)

windows: $(WINDOWS_TARGETS)

# Linux AMD64
$(LIBDIR)/linux_amd64/libchafa.so:
	mkdir -p $(CURDIR)/$(LIBDIR)/linux_amd64 && \
	mkdir -p build/chafa-linux-amd64 && cd build/chafa-linux-amd64 && \
	git clone --branch $(CHAFAVERSION) --depth 1 https://github.com/hpjansson/chafa.git . && \
	CC=gcc \
	CFLAGS="-m64" \
	LDFLAGS="-m64" \
	./autogen.sh --without-tools --host=x86_64-linux-gnu && make
	cp build/chafa-linux-amd64/chafa/.libs/libchafa.so $(CURDIR)/$(LIBDIR)/linux_amd64/libchafa.so

# Linux ARM64
$(LIBDIR)/linux_arm64/libchafa.so:
	mkdir -p $(CURDIR)/$(LIBDIR)/linux_arm64 && \
	mkdir -p build/chafa-linux-arm64 && cd build/chafa-linux-arm64 && \
	git clone --branch $(CHAFAVERSION) --depth 1 https://github.com/hpjansson/chafa.git . && \
	CC=aarch64-linux-gnu-gcc \
	./autogen.sh --without-tools --host=aarch64-linux-gnu && make
	cp build/chafa-linux-arm64/chafa/.libs/libchafa.so $(CURDIR)/$(LIBDIR)/linux_arm64/libchafa.so

# Linux 386
$(LIBDIR)/linux_386/libchafa.so:
	mkdir -p $(CURDIR)/$(LIBDIR)/linux_386 && \
	mkdir -p build/chafa-linux-386 && cd build/chafa-linux-386 && \
	git clone --branch $(CHAFAVERSION) --depth 1 https://github.com/hpjansson/chafa.git . && \
	CC=gcc \
	CFLAGS="-m32" \
	LDFLAGS="-m32" \
	./autogen.sh --without-tools --host=i686-linux-gnu && make
	cp build/chafa-linux-386/chafa/.libs/libchafa.so $(CURDIR)/$(LIBDIR)/linux_386/libchafa.so

# Darwin AMD64
$(LIBDIR)/darwin_amd64/libchafa.dylib:
	mkdir -p $(CURDIR)/$(LIBDIR)/darwin_amd64 && \
	mkdir -p build/chafa-darwin-amd64 && cd build/chafa-darwin-amd64 && \
	git clone --branch $(CHAFAVERSION) --depth 1 https://github.com/hpjansson/chafa.git . && \
	CC=clang \
	CFLAGS="-arch x86_64" \
	LDFLAGS="-arch x86_64" \
	LIBTOOL=glibtool \
	LIBTOOLIZE=glibtoolize \
	./autogen.sh --without-tools --host=x86_64-apple-darwin && make
	cp build/chafa-darwin-amd64/chafa/.libs/libchafa.dylib $(CURDIR)/$(LIBDIR)/darwin_amd64/libchafa.dylib

# Darwin ARM64 (Apple Silicon)
$(LIBDIR)/darwin_arm64/libchafa.dylib:
	mkdir -p $(CURDIR)/$(LIBDIR)/darwin_arm64 && \
	mkdir -p build/chafa-darwin-arm64 && cd build/chafa-darwin-arm64 && \
	git clone --branch $(CHAFAVERSION) --depth 1 https://github.com/hpjansson/chafa.git . && \
	CC=clang \
	CFLAGS="-arch arm64" \
	LDFLAGS="-arch arm64" \
	LIBTOOL=glibtool \
	LIBTOOLIZE=glibtoolize \
	./autogen.sh --without-tools --host=aarch64-apple-darwin && make
	cp build/chafa-darwin-arm64/chafa/.libs/libchafa.dylib $(CURDIR)/$(LIBDIR)/darwin_arm64/libchafa.dylib

# Windows x64 (MinGW)
# LDFLAGS statically links GLib and its dependencies into the DLL so that
# only libchafa.dll needs to be distributed (no MinGW runtime DLLs required).
$(LIBDIR)/windows_amd64/libchafa.dll:
	mkdir -p $(CURDIR)/$(LIBDIR)/windows_amd64 && \
	mkdir -p build/chafa-win-amd64 && cd build/chafa-win-amd64 && \
	git clone --branch $(CHAFAVERSION) --depth 1 https://github.com/hpjansson/chafa.git . && \
	CC=x86_64-w64-mingw32-gcc \
	CFLAGS="-O2" \
	LDFLAGS="-static-libgcc" \
	./autogen.sh --without-tools --host=x86_64-w64-mingw32 && make
	cp build/chafa-win-amd64/chafa/.libs/libchafa-0.dll $(CURDIR)/$(LIBDIR)/windows_amd64/libchafa.dll

clean:
	rm -rf build $(LIBDIR)/
