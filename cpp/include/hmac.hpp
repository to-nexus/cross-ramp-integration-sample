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
    static const std::string HMAC_SALT;
    
    // Generate HMAC signature for given data (following guide specification)
    static std::string GenerateHmac(const std::string& data);
    
    // Validate HMAC signature
    static bool ValidateHmac(const std::string& requestBody, const std::string& hmacSignature);
    
    // Convert bytes to hex string
    static std::string BytesToHex(const std::vector<unsigned char>& bytes);
    
    // Convert hex string to bytes
    static std::vector<unsigned char> HexToBytes(const std::string& hex);
    
private:
    // Base64 URL decoding function (as per guide specification)
    static std::vector<unsigned char> Base64UrlDecode(const std::string& str);
};

} // namespace hmac_utils 