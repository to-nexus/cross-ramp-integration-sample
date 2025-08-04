#pragma once

#include "models.hpp"
#include <string>
#include <vector>

namespace services {

class ValidationService {
public:
    static bool ValidateIntent(const models::ValidateRequest::Intent& intent);
    static std::string GenerateValidatorSignature(const std::string& user_sig, const std::string& digest);
    static bool ValidateAndProcessMint(const std::string& session_id, const std::vector<models::ValidateRequest::Intent::FromAsset>& from_assets);
};

class ExchangeService {
public:
    static bool ProcessExchangeResult(const std::string& session_id, const std::vector<models::ExchangeResultRequest::Intent::OutputAsset>& outputs, uint64_t receipt_status);
};

} // namespace services 