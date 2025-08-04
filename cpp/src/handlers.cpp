#include "handlers.hpp"
#include "database.hpp"
#include "services.hpp"
#include <iostream>
#include <sstream>

namespace handlers {

models::Response AssetHandler::GetAssetsHandler(const std::string& session_id, const std::string& language) {
    models::Response response;
    
    if (session_id.empty()) {
        response.success = false;
        response.error_code = "INVALID_SESSION_ID";
        return response;
    }
    
    auto session_assets = database::Database::GetInstance().GetOrCreateSessionAssets(session_id);
    if (!session_assets) {
        response.success = false;
        response.error_code = "DB_ERROR";
        return response;
    }
    
    // Convert assets map to vector
    std::vector<models::Asset> assets;
    for (const auto& [id, balance] : session_assets->assets) {
        models::Asset asset;
        asset.id = id;
        asset.balance = balance;
        assets.push_back(asset);
    }
    
    // Create V1Data
    models::V1Data v1_data;
    v1_data.player_id = session_id;
    v1_data.name = "playerName_" + session_id;
    v1_data.wallet_address = "0xaaaa";
    v1_data.server = "test";
    v1_data.assets = assets;
    
    // Create guide data
    json guide;
    guide["Authorization"] = "Bearer <token>";
    guide["X-Dapp-Authorization"] = "Bearer <token>";
    guide["X-Dapp-SessionID"] = session_id;
    guide["message"] = "The guide field displays header information at request time...";
    
    json session_info;
    session_info["created_at"] = session_assets->created_at;
    session_info["updated_at"] = session_assets->updated_at;
    guide["session_info"] = session_info;
    
    response.success = true;
    response.data.v1 = v1_data;
    response.data.guide = guide;
    
    return response;
}

models::ValidateResponse ValidationHandler::ValidateUserActionHandler(const models::ValidateRequest& request, const std::string& session_id) {
    models::ValidateResponse response;
    
    if (session_id.empty()) {
        response.success = false;
        response.error_code = "INVALID_SESSION_ID";
        return response;
    }
    
    // Validate intent
    if (!services::ValidationService::ValidateIntent(request.intent)) {
        response.success = false;
        response.error_code = "INVALID_INTENT";
        return response;
    }
    
    // Store UUID mapping
    if (!database::Database::GetInstance().StoreUUIDMapping(request.uuid, session_id)) {
        response.success = false;
        response.error_code = "UUID_MAPPING_FAILED";
        return response;
    }
    
    // For mint method, validate and deduct assets
    if (request.intent.method == "mint" || request.intent.method == "transfer") {
        if (!services::ValidationService::ValidateAndProcessMint(session_id, request.intent.from)) {
            response.success = false;
            response.error_code = "INSUFFICIENT_BALANCE";
            return response;
        }
    }
    
    // Generate validator signature
    std::string validator_sig = services::ValidationService::GenerateValidatorSignature(request.user_sig, request.digest);
    
    response.success = true;
    response.data.user_sig = request.user_sig;
    response.data.validator_sig = validator_sig;
    
    return response;
}

bool ExchangeResultHandler::ProcessExchangeResult(const models::ExchangeResp& request) {
    if (request.intent.outputs.empty()) {
        return true;
    }
    
    std::string session_id = database::Database::GetInstance().GetSessionIDByUUID(request.uuid);
    if (session_id.empty()) {
        std::cerr << "Invalid UUID or session not found: " << request.uuid << std::endl;
        return false;
    }
    
    uint64_t receipt_status = static_cast<uint64_t>(request.receipt.status);
    return services::ExchangeService::ProcessExchangeResult(session_id, request.intent.outputs, receipt_status);
}

} // namespace handlers 