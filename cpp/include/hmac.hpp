#pragma once

#include <string>
#include <vector>
#include <openssl/hmac.h>
#include <openssl/sha.h>
#include <iomanip>
#include <sstream>

namespace hmac_utils {

class HmacService {
public:
    // TODO: HMAC salt - In actual implementation, load from environment variables or configuration file
    static const std::string HMAC_SALT;
    
    // Generate HMAC signature for given data
    static std::string GenerateHmac(const std::string& data);
    
    // Validate HMAC signature
    static bool ValidateHmac(const std::string& requestBody, const std::string& hmacSignature);
    
    // Convert bytes to hex string
    static std::string BytesToHex(const std::vector<unsigned char>& bytes);
    
    // Convert hex string to bytes
    static std::vector<unsigned char> HexToBytes(const std::string& hex);
};

} // namespace hmac_utils 