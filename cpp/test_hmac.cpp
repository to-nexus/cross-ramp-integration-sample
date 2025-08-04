#include "include/hmac.hpp"
#include <iostream>
#include <nlohmann/json.hpp>

using json = nlohmann::json;

// TODO: HMAC salt - In actual implementation, load from environment variables or configuration file
const std::string salt = "my_secret_salt_value_!@#$%^&*"; // hmac key

struct Body {
    int userId;
    std::string username;
    std::string email;
    std::string role;
    int createdAt;
};

void testSha256() {
    Body body = {
        1234,
        "홍길동",
        "user@example.com",
        "admin",
        1234567890
    };

    json bodyJson = {
        {"userId", body.userId},
        {"username", body.username},
        {"email", body.email},
        {"role", body.role},
        {"createdAt", body.createdAt}
    };

    std::string jsonString = bodyJson.dump();
    std::cout << jsonString << std::endl;

    std::string hashString = hmac_utils::HmacService::GenerateHmac(jsonString);
    std::cout << "hashString: " << hashString << std::endl; // expected X-HMAC-Signature: f96cf60394f6b8ad3c6de2d5b2b1d1a540f9529082a8eb9cee405bfbdd9f37a1
}

int main() {
    testSha256();
    return 0;
} 