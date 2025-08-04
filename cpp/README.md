# Sample Game Backend API (C++ Version)

This is a C++ implementation of the sample game backend that manages session-specific asset information using in-memory storage.

## Prerequisites

- C++17 compatible compiler (GCC 7+, Clang 5+, MSVC 2017+)
- CMake 3.16 or higher
- nlohmann/json library

## Building

### 1. Install dependencies

```bash
# Ubuntu/Debian
sudo apt-get install nlohmann-json3-dev

# macOS
brew install nlohmann-json

# Or install via vcpkg
vcpkg install nlohmann-json
```

### 2. Build the project

```bash
mkdir build
cd build
cmake ..
make
```

### 3. Run the application

```bash
./sample_game_backend_cpp
```

## Project Structure

```
cpp/
├── CMakeLists.txt          # Build configuration
├── include/                # Header files
│   ├── config.hpp         # Configuration management
│   ├── database.hpp       # Database operations
│   ├── handlers.hpp       # HTTP request handlers
│   ├── models.hpp         # Data structures
│   └── services.hpp       # Business logic
├── src/                   # Source files
│   ├── main.cpp          # Application entry point
│   ├── config.cpp        # Configuration implementation
│   ├── database.cpp      # Database implementation
│   ├── handlers.cpp      # Handler implementation
│   └── services.cpp      # Service implementation
└── README.md             # This file
```
