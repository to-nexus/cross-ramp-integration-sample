#include "services.hpp"
#include "database.hpp"
#include "keystore.hpp"
#include <iostream>
#include <unordered_set>
#include <iomanip>
#include <sstream>

namespace services {

bool ValidationService::ValidateIntent(const models::ValidateRequest::Intent& intent) {
    // Validate allowed methods
    std::unordered_set<std::string> allowed_methods = {
        "mint", "transfer", "burn", "burn-permit"
    };
    
    if (allowed_methods.find(intent.method) == allowed_methods.end()) {
        return false;
    }
    
    // Special validation for mint method
    if (intent.method == "mint") {
        // Must have at least one from item
        if (intent.from.empty()) {
            return false;
        }
        
        // All from items must be asset type
        for (const auto& from : intent.from) {
            if (from.type != "asset") {
                return false;
            }
        }
    }
    
    return true;
}

std::string ValidationService::GenerateValidatorSignature(const std::string& user_sig, const std::string& digest) {
	// Convert digest string to bytes
	std::vector<uint8_t> digest_bytes;
	for (size_t i = 0; i < digest.length(); i += 2) {
		if (i + 1 < digest.length()) {
			std::string byte_string = digest.substr(i, 2);
			uint8_t byte = static_cast<uint8_t>(std::stoi(byte_string, nullptr, 16));
			digest_bytes.push_back(byte);
		}
	}
	
	// Use keystore service to sign
	static keystore::KeystoreService keystore_service;
	auto signature_bytes = keystore_service.Sign(digest_bytes);
	
	// Convert signature to hex string
	std::stringstream ss;
	ss << "0x";
	for (uint8_t byte : signature_bytes) {
		ss << std::hex << std::setfill('0') << std::setw(2) << static_cast<int>(byte);
	}
	
	return ss.str();
}

bool ValidationService::ValidateAndProcessMint(const std::string& session_id, const std::vector<models::ValidateRequest::Intent::FromAsset>& from_assets) {
    return database::Database::GetInstance().CheckAndDeductAssets(session_id, from_assets);
}

bool ExchangeService::ProcessExchangeResult(const std::string& session_id, const std::vector<models::ExchangeResp::Intent::OutputAsset>& outputs, uint64_t receipt_status) {
    // Skip processing if receipt status is not 0x1
    if (receipt_status != 1) {
        std::cout << "ProcessExchangeResult: skipped - receipt_status: " << receipt_status << std::endl;
        return true;
    }
    
    // Skip processing if output is empty
    if (outputs.empty()) {
        std::cout << "ProcessExchangeResult: skipped - empty outputs" << std::endl;
        return true;
    }
    
    // Process asset increase
    bool success = database::Database::GetInstance().AddAssets(session_id, outputs);
    if (!success) {
        std::cerr << "ProcessExchangeResult: failed to add assets" << std::endl;
        return false;
    }
    
    std::cout << "ProcessExchangeResult: assets added successfully" << std::endl;
    return true;
}

} // namespace services 