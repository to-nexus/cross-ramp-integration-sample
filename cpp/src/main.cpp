#include "config.hpp"
#include "database.hpp"
#include "handlers.hpp"
#include "server.hpp"
#include <iostream>
#include <string>
#include <signal.h>

static server::GameBackendServer* g_server = nullptr;

void SignalHandler(int signal) {
    if (g_server) {
        std::cout << "\nShutting down server..." << std::endl;
        g_server->Stop();
    }
    exit(0);
}

int main() {
    try {
        // Initialize configuration
        auto cfg = config::ConfigManager::InitConfig();
        
        // Initialize database
        if (!database::Database::GetInstance().InitDB()) {
            std::cerr << "Failed to initialize database" << std::endl;
            return 1;
        }
        
        // Create and start HTTP server
        server::GameBackendServer server(8080);
        g_server = &server;
        
        // Set up signal handlers for graceful shutdown
        signal(SIGINT, SignalHandler);
        signal(SIGTERM, SignalHandler);
        
        std::cout << "Server started on port 8080" << std::endl;
        std::cout << "API endpoint: http://localhost:8080/api/assets?language=ko" << std::endl;
        std::cout << "Order validation API: http://localhost:8080/api/validate" << std::endl;
        std::cout << "Health check: http://localhost:8080/health" << std::endl;
        std::cout << "Session-specific asset information is stored in memory" << std::endl;
        std::cout << "Press Ctrl+C to stop the server" << std::endl;
        
        // Start the server (this will block until server stops)
        if (!server.Start()) {
            std::cerr << "Failed to start server" << std::endl;
            return 1;
        }
        
        // Cleanup
        database::Database::GetInstance().CloseDB();
        
    } catch (const std::exception& e) {
        std::cerr << "Error: " << e.what() << std::endl;
        return 1;
    }
    
    return 0;
} 