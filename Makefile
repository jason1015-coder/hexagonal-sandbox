# ============================================================================
# Multi-Platform Makefile for SFML Game
# Supports: Windows (.exe), Linux (native + AppImage), macOS (.app)
# ============================================================================

# --- Platform Detection ---
UNAME := $(shell uname -s)
ifeq ($(UNAME),Darwin)
    PLATFORM := macos
else ifeq ($(UNAME),Linux)
    PLATFORM := linux
else
    PLATFORM := unknown
endif

# --- Project Configuration ---
PROJECT_NAME = game
SRC_DIR = src
BUILD_DIR_BASE = build
DIST_DIR = dist

# Source files
SRCS = $(wildcard $(SRC_DIR)/*.cpp)

# --- Common Compiler Flags ---
CXXFLAGS_COMMON = -O2 -Wall -std=c++17 -Isrc

# ============================================================================
# MACOS BUILD CONFIGURATION
# ============================================================================
ifeq ($(PLATFORM),macos)
    CXX = clang++
    BUILD_DIR = $(BUILD_DIR_BASE)/macos
    OBJS = $(patsubst $(SRC_DIR)/%.cpp, $(BUILD_DIR)/%.o, $(SRCS))
    
    # macOS specific flags
    CXXFLAGS = $(CXXFLAGS_COMMON)
    
    # Check if SFML is installed via Homebrew
    HOMEBREW_SFML = /opt/homebrew/Cellar/sfml
    HOMEBREW_SFML_INTEL = /usr/local/Cellar/sfml
    
    ifneq ($(wildcard $(HOMEBREW_SFML)/.),)
        SFML_VERSION := $(shell ls $(HOMEBREW_SFML) | sort -V | tail -n 1)
        SFML_PATH = $(HOMEBREW_SFML)/$(SFML_VERSION)
        CXXFLAGS += -I$(SFML_PATH)/include
        LDFLAGS = -L$(SFML_PATH)/lib -lsfml-graphics -lsfml-window -lsfml-network -lsfml-system
    else ifneq ($(wildcard $(HOMEBREW_SFML_INTEL)/.),)
        SFML_VERSION := $(shell ls $(HOMEBREW_SFML_INTEL) | sort -V | tail -n 1)
        SFML_PATH = $(HOMEBREW_SFML_INTEL)/$(SFML_VERSION)
        CXXFLAGS += -I$(SFML_PATH)/include
        LDFLAGS = -L$(SFML_PATH)/lib -lsfml-graphics -lsfml-window -lsfml-network -lsfml-system
    else
        # Fallback to system paths
        LDFLAGS = -lsfml-graphics -lsfml-window -lsfml-network -lsfml-system
    endif
    
    # Add rpath for dylib location
    LDFLAGS += -Wl,-rpath,@executable_path/../Frameworks
    
    TARGET = $(DIST_DIR)/$(PROJECT_NAME)
    APP_BUNDLE = $(DIST_DIR)/$(PROJECT_NAME).app
endif

# ============================================================================
# LINUX BUILD CONFIGURATION
# ============================================================================
ifeq ($(PLATFORM),linux)
    CXX = g++
    BUILD_DIR = $(BUILD_DIR_BASE)/linux
    OBJS = $(patsubst $(SRC_DIR)/%.cpp, $(BUILD_DIR)/%.o, $(SRCS))
    
    # Linux specific flags
    CXXFLAGS = $(CXXFLAGS_COMMON)
    LDFLAGS = -lsfml-graphics -lsfml-window -lsfml-network -lsfml-system -lpthread
    
    TARGET = $(DIST_DIR)/$(PROJECT_NAME)
    APPIMAGE = $(DIST_DIR)/$(PROJECT_NAME)-x86_64.AppImage
endif

# ============================================================================
# WINDOWS CROSS-COMPILATION CONFIGURATION
# ============================================================================
CXX_WIN = x86_64-w64-mingw32-g++-posix
BUILD_DIR_WIN = $(BUILD_DIR_BASE)/win
OBJS_WIN = $(patsubst $(SRC_DIR)/%.cpp, $(BUILD_DIR_WIN)/%.o, $(SRCS))

CXXFLAGS_WIN = -O2 -Wall -std=c++17 -DSFML_STATIC -D_GLIBCXX_HAS_GTHREADS \
               -Isrc -I./sfml_win/include

LDFLAGS_WIN = -L./sfml_win/lib \
              -lsfml-graphics-s -lsfml-window-s -lsfml-network-s -lsfml-system-s \
              -lfreetype -lopengl32 -lwinmm -lgdi32 -lws2_32 \
              -lwinpthread -static-libgcc -static-libstdc++ -static

TARGET_WIN = $(DIST_DIR)/$(PROJECT_NAME).exe

# ============================================================================
# MAIN TARGETS
# ============================================================================

.PHONY: all clean help setup

all: native

help:
	@echo "Multi-Platform Build System"
	@echo "============================"
	@echo ""
	@echo "Current platform: $(PLATFORM)"
	@echo ""
	@echo "Available targets:"
	@echo "  make native      - Build for current platform ($(PLATFORM))"
	@echo "  make macos       - Build macOS executable"
	@echo "  make macos-app   - Build macOS .app bundle"
	@echo "  make linux       - Build Linux executable"
	@echo "  make appimage    - Build Linux AppImage"
	@echo "  make windows     - Cross-compile Windows .exe"
	@echo "  make all-platforms - Build for all platforms"
	@echo "  make clean       - Remove all build artifacts"
	@echo "  make setup-*     - Setup SFML for specific platform"
	@echo ""

native:
ifeq ($(PLATFORM),macos)
	@$(MAKE) macos
else ifeq ($(PLATFORM),linux)
	@$(MAKE) linux
else
	@echo "ERROR: Native builds only supported on macOS and Linux"
	@echo "For Windows .exe, use: make windows"
	@exit 1
endif

all-platforms:
	@echo "Building for all platforms..."
	@if [ "$(PLATFORM)" = "macos" ]; then $(MAKE) macos-app; fi
	@if [ "$(PLATFORM)" = "linux" ]; then $(MAKE) appimage; fi
	@$(MAKE) windows
	@echo ""
	@echo "All platforms built successfully!"
	@ls -lh $(DIST_DIR)/

# ============================================================================
# MACOS TARGETS
# ============================================================================

.PHONY: macos macos-app check-sfml-macos setup-sfml-macos

macos: check-sfml-macos setup $(TARGET)
	@echo ""
	@echo "✓ macOS build complete: $(TARGET)"
	@echo "  Run with: ./$(TARGET)"

macos-app: check-sfml-macos setup $(APP_BUNDLE)
	@echo ""
	@echo "✓ macOS app bundle complete: $(APP_BUNDLE)"
	@echo "  Run with: open $(APP_BUNDLE)"

check-sfml-macos:
ifeq ($(PLATFORM),macos)
	@command -v brew >/dev/null 2>&1 || { echo "ERROR: Homebrew not found. Install from https://brew.sh"; exit 1; }
	@brew list sfml >/dev/null 2>&1 || { echo "ERROR: SFML not installed."; echo "Run: make setup-sfml-macos"; exit 1; }
endif

setup-sfml-macos:
ifeq ($(PLATFORM),macos)
	@echo "Installing SFML via Homebrew..."
	brew install sfml
	@echo "✓ SFML installed successfully!"
else
	@echo "ERROR: This target is only for macOS"
	@exit 1
endif

$(TARGET): $(OBJS)
	@mkdir -p $(DIST_DIR)
	$(CXX) $(OBJS) -o $@ $(LDFLAGS)
	@echo "Copying SFML frameworks/dylibs..."
	@mkdir -p $(DIST_DIR)/lib
ifneq ($(SFML_PATH),)
	@cp -f $(SFML_PATH)/lib/*.dylib $(DIST_DIR)/lib/ 2>/dev/null || true
endif

$(APP_BUNDLE): $(OBJS)
	@echo "Creating macOS app bundle..."
	@mkdir -p $(APP_BUNDLE)/Contents/MacOS
	@mkdir -p $(APP_BUNDLE)/Contents/Resources
	@mkdir -p $(APP_BUNDLE)/Contents/Frameworks
	$(CXX) $(OBJS) -o $(APP_BUNDLE)/Contents/MacOS/$(PROJECT_NAME) $(LDFLAGS)
	@echo "Copying SFML frameworks..."
ifneq ($(SFML_PATH),)
	@cp -R $(SFML_PATH)/lib/*.dylib $(APP_BUNDLE)/Contents/Frameworks/ 2>/dev/null || true
endif
	@echo '<?xml version="1.0" encoding="UTF-8"?>' > $(APP_BUNDLE)/Contents/Info.plist
	@echo '<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">' >> $(APP_BUNDLE)/Contents/Info.plist
	@echo '<plist version="1.0">' >> $(APP_BUNDLE)/Contents/Info.plist
	@echo '<dict>' >> $(APP_BUNDLE)/Contents/Info.plist
	@echo '    <key>CFBundleExecutable</key>' >> $(APP_BUNDLE)/Contents/Info.plist
	@echo '    <string>$(PROJECT_NAME)</string>' >> $(APP_BUNDLE)/Contents/Info.plist
	@echo '    <key>CFBundleIdentifier</key>' >> $(APP_BUNDLE)/Contents/Info.plist
	@echo '    <string>com.hexagonal.$(PROJECT_NAME)</string>' >> $(APP_BUNDLE)/Contents/Info.plist
	@echo '    <key>CFBundleName</key>' >> $(APP_BUNDLE)/Contents/Info.plist
	@echo '    <string>$(PROJECT_NAME)</string>' >> $(APP_BUNDLE)/Contents/Info.plist
	@echo '    <key>CFBundleVersion</key>' >> $(APP_BUNDLE)/Contents/Info.plist
	@echo '    <string>1.0</string>' >> $(APP_BUNDLE)/Contents/Info.plist
	@echo '    <key>LSMinimumSystemVersion</key>' >> $(APP_BUNDLE)/Contents/Info.plist
	@echo '    <string>10.13</string>' >> $(APP_BUNDLE)/Contents/Info.plist
	@echo '</dict>' >> $(APP_BUNDLE)/Contents/Info.plist
	@echo '</plist>' >> $(APP_BUNDLE)/Contents/Info.plist

$(BUILD_DIR)/%.o: $(SRC_DIR)/%.cpp
	@mkdir -p $(BUILD_DIR)
	$(CXX) $(CXXFLAGS) -c $< -o $@

# ============================================================================
# LINUX TARGETS
# ============================================================================

.PHONY: linux appimage check-sfml-linux setup-sfml-linux

linux: check-sfml-linux setup $(TARGET)
	@echo ""
	@echo "✓ Linux build complete: $(TARGET)"
	@echo "  Run with: ./$(TARGET)"

appimage: check-sfml-linux setup $(APPIMAGE)
	@echo ""
	@echo "✓ AppImage created: $(APPIMAGE)"
	@echo "  Run with: ./$(APPIMAGE)"

check-sfml-linux:
ifeq ($(PLATFORM),linux)
	@ldconfig -p | grep libsfml >/dev/null 2>&1 || { echo "ERROR: SFML not installed."; echo "Run: make setup-sfml-linux"; exit 1; }
endif

setup-sfml-linux:
ifeq ($(PLATFORM),linux)
	@echo "Installing SFML..."
	@if command -v apt-get >/dev/null 2>&1; then \
		sudo apt-get update && sudo apt-get install -y libsfml-dev; \
	elif command -v dnf >/dev/null 2>&1; then \
		sudo dnf install -y SFML-devel; \
	elif command -v pacman >/dev/null 2>&1; then \
		sudo pacman -S --noconfirm sfml; \
	else \
		echo "ERROR: Unknown package manager. Please install SFML manually."; \
		exit 1; \
	fi
	@echo "✓ SFML installed successfully!"
else
	@echo "ERROR: This target is only for Linux"
	@exit 1
endif

$(TARGET): $(OBJS)
	@mkdir -p $(DIST_DIR)
	$(CXX) $(OBJS) -o $@ $(LDFLAGS)

$(BUILD_DIR)/%.o: $(SRC_DIR)/%.cpp
	@mkdir -p $(BUILD_DIR)
	$(CXX) $(CXXFLAGS) -c $< -o $@

$(APPIMAGE): $(TARGET)
	@echo "Creating AppImage..."
	@mkdir -p $(BUILD_DIR)/AppDir/usr/bin
	@mkdir -p $(BUILD_DIR)/AppDir/usr/lib
	@mkdir -p $(BUILD_DIR)/AppDir/usr/share/applications
	@mkdir -p $(BUILD_DIR)/AppDir/usr/share/icons/hicolor/256x256/apps
	@cp $(TARGET) $(BUILD_DIR)/AppDir/usr/bin/
	@echo "Copying SFML libraries..."
	@ldd $(TARGET) | grep sfml | awk '{print $$3}' | xargs -I {} cp {} $(BUILD_DIR)/AppDir/usr/lib/ || true
	@echo "[Desktop Entry]" > $(BUILD_DIR)/AppDir/$(PROJECT_NAME).desktop
	@echo "Type=Application" >> $(BUILD_DIR)/AppDir/$(PROJECT_NAME).desktop
	@echo "Name=$(PROJECT_NAME)" >> $(BUILD_DIR)/AppDir/$(PROJECT_NAME).desktop
	@echo "Exec=$(PROJECT_NAME)" >> $(BUILD_DIR)/AppDir/$(PROJECT_NAME).desktop
	@echo "Icon=$(PROJECT_NAME)" >> $(BUILD_DIR)/AppDir/$(PROJECT_NAME).desktop
	@echo "Categories=Game;" >> $(BUILD_DIR)/AppDir/$(PROJECT_NAME).desktop
	@cp $(BUILD_DIR)/AppDir/$(PROJECT_NAME).desktop $(BUILD_DIR)/AppDir/usr/share/applications/
	@if command -v appimagetool >/dev/null 2>&1; then \
		ARCH=x86_64 appimagetool $(BUILD_DIR)/AppDir $(APPIMAGE); \
	else \
		echo "WARNING: appimagetool not found. Installing..."; \
		wget -q https://github.com/AppImage/AppImageKit/releases/download/continuous/appimagetool-x86_64.AppImage -O /tmp/appimagetool; \
		chmod +x /tmp/appimagetool; \
		ARCH=x86_64 /tmp/appimagetool $(BUILD_DIR)/AppDir $(APPIMAGE); \
	fi

# ============================================================================
# WINDOWS TARGETS
# ============================================================================

.PHONY: windows win check-sfml-win setup-sfml-win

windows: win

win: check-sfml-win setup $(TARGET_WIN)
	@echo ""
	@echo "✓ Windows build complete: $(TARGET_WIN)"

check-sfml-win:
	@command -v $(CXX_WIN) >/dev/null 2>&1 || { echo "ERROR: MinGW-w64 not found."; echo "Install it first, then run: make setup-sfml-win"; exit 1; }
	@if [ ! -d "./sfml_win/include" ] || [ ! -d "./sfml_win/lib" ]; then \
		echo "ERROR: SFML Windows libraries not found in ./sfml_win/"; \
		echo ""; \
		echo "To fix this, run: make setup-sfml-win"; \
		exit 1; \
	fi

setup-sfml-win:
	@echo "Setting up SFML for Windows cross-compilation..."
	@mkdir -p sfml_win
	@if [ ! -d "./sfml_win/include" ] || [ ! -d "./sfml_win/lib" ]; then \
		echo "Downloading SFML 2.5.1 for Windows..."; \
		cd sfml_win && \
		wget -q --show-progress https://www.sfml-dev.org/files/SFML-2.5.1-windows-gcc-7.3.0-mingw-64-bit.zip -O sfml.zip && \
		unzip -q sfml.zip && \
		mv SFML-2.5.1/* . && \
		rmdir SFML-2.5.1 && \
		rm sfml.zip && \
		echo "✓ SFML Windows setup complete!"; \
	else \
		echo "SFML already set up in ./sfml_win/"; \
	fi

$(TARGET_WIN): $(OBJS_WIN)
	@mkdir -p $(DIST_DIR)
	$(CXX_WIN) $(OBJS_WIN) -o $@ $(LDFLAGS_WIN)

$(BUILD_DIR_WIN)/%.o: $(SRC_DIR)/%.cpp
	@mkdir -p $(BUILD_DIR_WIN)
	$(CXX_WIN) $(CXXFLAGS_WIN) -c $< -o $@

# ============================================================================
# UTILITY TARGETS
# ============================================================================

setup:
	@mkdir -p $(BUILD_DIR) $(DIST_DIR)

clean:
	rm -rf build/ dist/
	@echo "✓ Clean complete"

.PHONY: test-run

test-run:
ifeq ($(PLATFORM),macos)
	@if [ -d "$(APP_BUNDLE)" ]; then \
		open $(APP_BUNDLE); \
	elif [ -f "$(TARGET)" ]; then \
		./$(TARGET); \
	else \
		echo "No executable found. Run 'make' first."; \
	fi
else ifeq ($(PLATFORM),linux)
	@if [ -f "$(APPIMAGE)" ]; then \
		./$(APPIMAGE); \
	elif [ -f "$(TARGET)" ]; then \
		./$(TARGET); \
	else \
		echo "No executable found. Run 'make' first."; \
	fi
else
	@echo "test-run is only supported on macOS and Linux"
	@echo "For Windows .exe, copy $(TARGET_WIN) to a Windows machine"
endif