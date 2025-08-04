#pragma once

#include "handlers.hpp"
#include "models.hpp"
#include <httplib.h>
#include <string>
#include <memory>

namespace server {

class GameBackendServer {
public:
    GameBackendServer(int port = 8080);
    ~GameBackendServer() = default;
    
    bool Start();
    void Stop();
    
private:
    void SetupRoutes();
    void SetupCORS();
    
    // Route handlers
    void HandleGetAssets(const httplib::Request& req, httplib::Response& res);
    void HandleValidateUserAction(const httplib::Request& req, httplib::Response& res);
    void HandleExchangeResult(const httplib::Request& req, httplib::Response& res);
    void HandleHealthCheck(const httplib::Request& req, httplib::Response& res);
    
    // Helper functions
    std::string GetSessionIDFromHeaders(const httplib::Request& req);
    void SetCORSHeaders(httplib::Response& res);
    bool ValidateHmacRequest(const httplib::Request& req, httplib::Response& res);
    
private:
    std::unique_ptr<httplib::Server> server_;
    int port_;
};

} // namespace server 