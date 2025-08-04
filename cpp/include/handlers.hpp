#pragma once

#include "models.hpp"
#include <string>
#include <memory>

namespace handlers {

class AssetHandler {
public:
    static models::Response GetAssetsHandler(const std::string& session_id, const std::string& language = "ko");
};

class ValidationHandler {
public:
    static models::ValidateResponse ValidateUserActionHandler(const models::ValidateRequest& request, const std::string& session_id);
};

class ExchangeResultHandler {
public:
    static bool ProcessExchangeResult(const models::ExchangeResultRequest& request);
};

} // namespace handlers 