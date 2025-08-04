#pragma once

#include <string>
#include <random>
#include <chrono>

namespace config {

struct DBConfig {
    std::string path;
};

struct Config {
    std::string port;
    DBConfig db;
};

class ConfigManager {
public:
    static Config InitConfig();
    
private:
    static void InitRandomSeed();
};

} // namespace config 