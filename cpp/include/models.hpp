#pragma once

#include <string>
#include <vector>
#include <map>
#include <memory>
#include <nlohmann/json.hpp>

using json = nlohmann::json;

namespace models {

// Asset asset information structure
struct Asset {
    std::string id;
    std::string balance;
    
    NLOHMANN_DEFINE_TYPE_INTRUSIVE(Asset, id, balance)
};

// SessionAssets session-specific asset information structure
struct SessionAssets {
    std::string session_id;
    std::map<std::string, std::string> assets;
    std::string created_at;
    std::string updated_at;
    
    NLOHMANN_DEFINE_TYPE_INTRUSIVE(SessionAssets, session_id, assets, created_at, updated_at)
};

// V1Data v1 guide data structure
struct V1Data {
    std::string player_id;
    std::string name;
    std::string wallet_address;
    std::string server;
    std::vector<Asset> assets;
    
    NLOHMANN_DEFINE_TYPE_INTRUSIVE(V1Data, player_id, name, wallet_address, server, assets)
};

// Response API response structure
struct Response {
    bool success;
    std::string error_code;
    struct {
        V1Data v1;
        json guide;
    } data;
    
    NLOHMANN_DEFINE_TYPE_INTRUSIVE(Response, success, error_code, data)
};

// Order validation request structure
struct ValidateRequest {
    std::string uuid;
    std::string user_sig;
    std::string user_address;
    std::string project_id;
    std::string digest;
    
    struct Intent {
        std::string type;
        std::string method;
        struct FromAsset {
            std::string type;
            std::string id;
            int amount;
        };
        struct ToAsset {
            std::string type;
            std::string id;
            int amount;
        };
        std::vector<FromAsset> from;
        std::vector<ToAsset> to;
    } intent;
    
    NLOHMANN_DEFINE_TYPE_INTRUSIVE(ValidateRequest, uuid, user_sig, user_address, project_id, digest, intent)
};

// Order validation response structure
struct ValidateResponse {
    bool success;
    std::string error_code;
    struct {
        std::string user_sig;
        std::string validator_sig;
    } data;
    
    NLOHMANN_DEFINE_TYPE_INTRUSIVE(ValidateResponse, success, error_code, data)
};

// Exchange result request structure
struct ExchangeResultRequest {
    std::string uuid;
    std::string tx_hash;
    struct {
        int status;
    } receipt;
    struct Intent {
        std::string type;
        std::string method;
        struct PairAsset {
            std::string type;
            std::string id;
            int amount;
        };
        std::vector<PairAsset> from;
        std::vector<PairAsset> to;
    };
    Intent intent;
    
    NLOHMANN_DEFINE_TYPE_INTRUSIVE(ExchangeResultRequest, uuid, tx_hash, receipt, intent)
};

} // namespace models 