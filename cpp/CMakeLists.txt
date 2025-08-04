cmake_minimum_required(VERSION 3.16)
project(sample_game_backend_cpp)

set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

# Find required packages
find_package(PkgConfig REQUIRED)
find_package(Threads REQUIRED)
find_package(nlohmann_json 3.2.0 REQUIRED)

# Find OpenSSL
find_package(OpenSSL REQUIRED)

# Add cpp-httplib (header-only library)
include_directories(/opt/homebrew/include)

# Include directories
include_directories(${CMAKE_SOURCE_DIR}/include)
include_directories(${CMAKE_SOURCE_DIR}/src)

# Add source files
set(SOURCES
    src/main.cpp
    src/config.cpp
    src/database.cpp
    src/handlers.cpp
    src/services.cpp
    src/server.cpp
    src/keystore.cpp
    src/hmac.cpp
)

# Create executable
add_executable(${PROJECT_NAME} ${SOURCES})

# Link libraries
target_link_libraries(${PROJECT_NAME} 
    Threads::Threads
    nlohmann_json::nlohmann_json
    OpenSSL::SSL
    OpenSSL::Crypto
    ${CMAKE_DL_LIBS}
)

# Compiler flags
target_compile_options(${PROJECT_NAME} PRIVATE -Wall -Wextra)

# Install target
install(TARGETS ${PROJECT_NAME} DESTINATION bin) 