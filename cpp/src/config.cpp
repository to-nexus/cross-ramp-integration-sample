#include "config.hpp"
#include <iostream>

namespace config {

void ConfigManager::InitRandomSeed() {
    auto now = std::chrono::high_resolution_clock::now();
    auto seed = now.time_since_epoch().count();
    std::srand(static_cast<unsigned int>(seed));
}

Config ConfigManager::InitConfig() {
    InitRandomSeed();
    
    Config cfg;
    cfg.port = ":8080";
    cfg.db.path = "./session_db";
    
    return cfg;
}

} // namespace config 