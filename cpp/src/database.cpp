#include "database.hpp"
#include <iostream>
#include <random>
#include <sstream>
#include <iomanip>

namespace database {

Database& Database::GetInstance() {
    static Database instance;
    return instance;
}

bool Database::InitDB() {
    std::cout << "Database initialized successfully" << std::endl;
    return true;
}

void Database::CloseDB() {
    // In-memory database, no cleanup needed
}

std::map<std::string, std::string> Database::GenerateRandomAssets() {
    std::map<std::string, std::string> assets;
    std::random_device rd;
    std::mt19937 gen(rd());
    std::uniform_int_distribution<> dis(500, 5000);
    
    assets["asset_money"] = std::to_string(dis(gen));
    assets["asset_gold"] = std::to_string(dis(gen));
    assets["item_gem"] = std::to_string(dis(gen));
    assets["item_banana"] = std::to_string(dis(gen));
    assets["asset_silver"] = std::to_string(dis(gen));
    assets["item_apple"] = std::to_string(dis(gen));
    assets["item_fish"] = std::to_string(dis(gen));
    assets["item_branch"] = std::to_string(dis(gen));
    assets["item_horn"] = std::to_string(dis(gen));
    assets["item_maple"] = std::to_string(dis(gen));
    
    return assets;
}

std::string Database::GetCurrentTimeString() {
    auto now = std::chrono::system_clock::now();
    auto time_t = std::chrono::system_clock::to_time_t(now);
    std::stringstream ss;
    ss << std::put_time(std::gmtime(&time_t), "%Y-%m-%dT%H:%M:%SZ");
    return ss.str();
}

std::shared_ptr<models::SessionAssets> Database::GetOrCreateSessionAssets(const std::string& session_id) {
    std::lock_guard<std::mutex> lock(session_mutex_);
    
    auto it = session_assets_.find(session_id);
    if (it != session_assets_.end()) {
        return std::make_shared<models::SessionAssets>(it->second);
    }
    
    // Create new session assets
    models::SessionAssets new_assets;
    new_assets.session_id = session_id;
    new_assets.assets = GenerateRandomAssets();
    new_assets.created_at = GetCurrentTimeString();
    new_assets.updated_at = GetCurrentTimeString();
    
    session_assets_[session_id] = new_assets;
    return std::make_shared<models::SessionAssets>(new_assets);
}

bool Database::CheckAndDeductAssets(const std::string& session_id, const std::vector<models::ValidateRequest::Intent::FromAsset>& from_assets) {
    std::lock_guard<std::mutex> lock(session_mutex_);
    
    auto it = session_assets_.find(session_id);
    if (it == session_assets_.end()) {
        return false;
    }
    
    for (const auto& asset : from_assets) {
        auto asset_it = it->second.assets.find(asset.id);
        if (asset_it == it->second.assets.end()) {
            return false;
        }
        
        int current_balance = std::stoi(asset_it->second);
        if (current_balance < asset.amount) {
            return false;
        }
        
        it->second.assets[asset.id] = std::to_string(current_balance - asset.amount);
    }
    
    it->second.updated_at = GetCurrentTimeString();
    return true;
}

bool Database::AddAssets(const std::string& session_id, const std::vector<models::ExchangeResp::Intent::OutputAsset>& outputs) {
    std::lock_guard<std::mutex> lock(session_mutex_);
    
    auto it = session_assets_.find(session_id);
    if (it == session_assets_.end()) {
        return false;
    }
    
    for (const auto& output : outputs) {
        auto asset_it = it->second.assets.find(output.asset_id);
        if (asset_it == it->second.assets.end()) {
            it->second.assets[output.asset_id] = std::to_string(output.amount);
        } else {
            int current_balance = std::stoi(asset_it->second);
            it->second.assets[output.asset_id] = std::to_string(current_balance + output.amount);
        }
    }
    
    it->second.updated_at = GetCurrentTimeString();
    return true;
}

bool Database::StoreUUIDMapping(const std::string& uuid, const std::string& session_id) {
    std::lock_guard<std::mutex> lock(uuid_mutex_);
    uuid_mapping_[uuid] = session_id;
    return true;
}

std::string Database::GetSessionIDByUUID(const std::string& uuid) {
    std::lock_guard<std::mutex> lock(uuid_mutex_);
    auto it = uuid_mapping_.find(uuid);
    if (it == uuid_mapping_.end()) {
        return "";
    }
    return it->second;
}

} // namespace database 