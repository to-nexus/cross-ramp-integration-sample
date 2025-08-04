#include "server.hpp"
#include "hmac.hpp"
#include <iostream>
#include <nlohmann/json.hpp>

using json = nlohmann::json;

namespace server {

GameBackendServer::GameBackendServer(int port) : port_(port) {
    server_ = std::make_unique<httplib::Server>();
    SetupCORS();
    SetupRoutes();
}

void GameBackendServer::SetupCORS() {
    server_->set_default_headers({
        {"Access-Control-Allow-Origin", "*"},
        {"Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS"},
        {"Access-Control-Allow-Headers", "Authorization, X-Dapp-Authorization, X-Dapp-SessionID, Content-Type"}
    });
}

void GameBackendServer::SetupRoutes() {
    // Health check endpoint
    server_->Get("/health", [this](const httplib::Request& req, httplib::Response& res) {
        HandleHealthCheck(req, res);
    });
    
    // API routes
    server_->Get("/api/assets", [this](const httplib::Request& req, httplib::Response& res) {
        HandleGetAssets(req, res);
    });
    
    server_->Post("/api/validate", [this](const httplib::Request& req, httplib::Response& res) {
        if (!ValidateHmacRequest(req, res)) {
            return;
        }
        HandleValidateUserAction(req, res);
    });
    
    server_->Post("/api/result", [this](const httplib::Request& req, httplib::Response& res) {
        if (!ValidateHmacRequest(req, res)) {
            return;
        }
        HandleExchangeResult(req, res);
    });
    
    // Handle OPTIONS for CORS
    server_->Options(".*", [this](const httplib::Request& req, httplib::Response& res) {
        SetCORSHeaders(res);
        res.status = 204;
    });
}

void GameBackendServer::HandleHealthCheck(const httplib::Request& req, httplib::Response& res) {
    SetCORSHeaders(res);
    
    json response = {
        {"status", "healthy"},
        {"message", "Server is running normally"}
    };
    
    res.set_content(response.dump(), "application/json");
}

void GameBackendServer::HandleGetAssets(const httplib::Request& req, httplib::Response& res) {
    SetCORSHeaders(res);
    
    std::string session_id = GetSessionIDFromHeaders(req);
    std::string language = req.get_param_value("language");
    if (language.empty()) {
        language = "ko";
    }
    
    auto response = handlers::AssetHandler::GetAssetsHandler(session_id, language);
    
    json json_response;
    json_response["success"] = response.success;
    if (!response.error_code.empty()) {
        json_response["errorCode"] = response.error_code;
    }
    if (response.success) {
        // Manually construct data object
        json data;
        json v1;
        v1["player_id"] = response.data.v1.player_id;
        v1["name"] = response.data.v1.name;
        v1["wallet_address"] = response.data.v1.wallet_address;
        v1["server"] = response.data.v1.server;
        
        json assets = json::array();
        for (const auto& asset : response.data.v1.assets) {
            json asset_json;
            asset_json["id"] = asset.id;
            asset_json["balance"] = asset.balance;
            assets.push_back(asset_json);
        }
        v1["assets"] = assets;
        
        data["v1"] = v1;
        data["guide"] = response.data.guide;
        
        json_response["data"] = data;
    }
    
    res.set_content(json_response.dump(), "application/json");
}

void GameBackendServer::HandleValidateUserAction(const httplib::Request& req, httplib::Response& res) {
    SetCORSHeaders(res);
    
    std::string session_id = GetSessionIDFromHeaders(req);
    
    try {
        auto request_json = json::parse(req.body);
        models::ValidateRequest request;
        
        // Parse request (simplified - in real implementation, use proper JSON parsing)
        request.uuid = request_json["uuid"];
        request.user_sig = request_json["user_sig"];
        request.user_address = request_json["user_address"];
        request.project_id = request_json["project_id"];
        request.digest = request_json["digest"];
        
        // Parse intent
        auto intent_json = request_json["intent"];
        request.intent.method = intent_json["method"];
        
        // Parse from assets
        auto from_json = intent_json["from"];
        for (const auto& from_item : from_json) {
            models::ValidateRequest::Intent::FromAsset asset;
            asset.type = from_item["type"];
            asset.id = from_item["id"];
            asset.amount = from_item["amount"];
            request.intent.from.push_back(asset);
        }
        
        // Parse to assets
        auto to_json = intent_json["to"];
        for (const auto& to_item : to_json) {
            models::ValidateRequest::Intent::ToAsset asset;
            asset.type = to_item["type"];
            asset.id = to_item["id"];
            asset.amount = to_item["amount"];
            request.intent.to.push_back(asset);
        }
        
        auto response = handlers::ValidationHandler::ValidateUserActionHandler(request, session_id);
        
        json json_response;
        json_response["success"] = response.success;
        if (!response.error_code.empty()) {
            json_response["errorCode"] = response.error_code;
        }
        if (response.success) {
            json_response["data"]["userSig"] = response.data.user_sig;
            json_response["data"]["validatorSig"] = response.data.validator_sig;
        }
        
        res.set_content(json_response.dump(), "application/json");
        
    } catch (const std::exception& e) {
        json error_response = {
            {"success", false},
            {"errorCode", "INVALID_REQUEST"}
        };
        res.status = 400;
        res.set_content(error_response.dump(), "application/json");
    }
}

void GameBackendServer::HandleExchangeResult(const httplib::Request& req, httplib::Response& res) {
    SetCORSHeaders(res);
    
    try {
        auto request_json = json::parse(req.body);
        models::ExchangeResultRequest request;
        
        // Parse request (simplified)
        request.uuid = request_json["uuid"];
        request.tx_hash = request_json["tx_hash"];
        request.receipt.status = request_json["receipt"]["status"];
        
        // Parse outputs (simplified)
        auto outputs_json = request_json["intent"]["outputs"];
        for (const auto& output : outputs_json) {
            models::ExchangeResultRequest::Intent::OutputAsset asset;
            asset.asset_id = output["asset_id"];
            asset.amount = output["amount"];
            request.intent.outputs.push_back(asset);
        }
        
        bool success = handlers::ExchangeResultHandler::ProcessExchangeResult(request);
        
        json response = {
            {"success", success}
        };
        
        res.set_content(response.dump(), "application/json");
        
    } catch (const std::exception& e) {
        json error_response = {
            {"success", false},
            {"error", "Failed to parse request"}
        };
        res.status = 400;
        res.set_content(error_response.dump(), "application/json");
    }
}

std::string GameBackendServer::GetSessionIDFromHeaders(const httplib::Request& req) {
    auto it = req.headers.find("X-Dapp-SessionID");
    if (it != req.headers.end()) {
        return it->second;
    }
    return "";
}

void GameBackendServer::SetCORSHeaders(httplib::Response& res) {
    res.set_header("Access-Control-Allow-Origin", "*");
    res.set_header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS");
    res.set_header("Access-Control-Allow-Headers", "Authorization, X-Dapp-Authorization, X-Dapp-SessionID, Content-Type, X-HMAC-Signature");
}

bool GameBackendServer::ValidateHmacRequest(const httplib::Request& req, httplib::Response& res) {
    // Get HMAC signature from header
    auto hmac_it = req.headers.find("X-HMAC-Signature");
    if (hmac_it == req.headers.end()) {
        json error_response = {
            {"success", false},
            {"errorCode", "INVALID_HMAC_SIGNATURE"},
            {"message", "Missing HMAC signature"}
        };
        res.status = 401;
        res.set_content(error_response.dump(), "application/json");
        return false;
    }
    
    // Validate HMAC
    if (!hmac_utils::HmacService::ValidateHmac(req.body, hmac_it->second)) {
        json error_response = {
            {"success", false},
            {"errorCode", "INVALID_HMAC_SIGNATURE"},
            {"message", "Invalid HMAC signature"}
        };
        res.status = 401;
        res.set_content(error_response.dump(), "application/json");
        return false;
    }
    
    return true;
}

bool GameBackendServer::Start() {
    std::cout << "Starting server on port " << port_ << std::endl;
    return server_->listen("0.0.0.0", port_);
}

void GameBackendServer::Stop() {
    server_->stop();
}

} // namespace server 