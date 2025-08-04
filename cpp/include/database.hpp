#pragma once

#include "models.hpp"
#include <memory>
#include <unordered_map>
#include <mutex>
#include <string>
#include <chrono>
#include <iomanip>
#include <sstream>

namespace database {

class Database {
public:
    static Database& GetInstance();
    
    bool InitDB();
    void CloseDB();
    
    std::shared_ptr<models::SessionAssets> GetOrCreateSessionAssets(const std::string& session_id);
    bool CheckAndDeductAssets(const std::string& session_id, const std::vector<models::ValidateRequest::Intent::FromAsset>& from_assets);
    bool AddAssets(const std::string& session_id, const std::vector<models::ExchangeResultRequest::Intent::OutputAsset>& outputs);
    bool StoreUUIDMapping(const std::string& uuid, const std::string& session_id);
    std::string GetSessionIDByUUID(const std::string& uuid);

private:
    Database() = default;
    ~Database() = default;
    Database(const Database&) = delete;
    Database& operator=(const Database&) = delete;
    
    std::unordered_map<std::string, models::SessionAssets> session_assets_;
    std::unordered_map<std::string, std::string> uuid_mapping_;
    std::mutex session_mutex_;
    std::mutex uuid_mutex_;
    
    std::map<std::string, std::string> GenerateRandomAssets();
    std::string GetCurrentTimeString();
};

} // namespace database 