#include "hmac.hpp"
#include <cstring>
#include <algorithm>

namespace hmac_utils {

// TODO: HMAC salt - In actual implementation, load from environment variables or configuration file
const std::string HmacService::HMAC_SALT = "my_secret_salt_value_!@#$%^&*";

std::string HmacService::GenerateHmac(const std::string& data) {
    std::vector<unsigned char> hmac(SHA256_DIGEST_LENGTH);
    
    HMAC(EVP_sha256(), 
         HMAC_SALT.c_str(), 
         HMAC_SALT.length(),
         reinterpret_cast<const unsigned char*>(data.c_str()),
         data.length(),
         hmac.data(),
         nullptr);
    
    return BytesToHex(hmac);
}

bool HmacService::ValidateHmac(const std::string& requestBody, const std::string& hmacSignature) {
    if (hmacSignature.empty()) {
        return false;
    }
    
    std::string calculatedHmac = GenerateHmac(requestBody);
    
    // Case-insensitive comparison
    std::string lowerCalculated = calculatedHmac;
    std::string lowerSignature = hmacSignature;
    
    std::transform(lowerCalculated.begin(), lowerCalculated.end(), lowerCalculated.begin(), ::tolower);
    std::transform(lowerSignature.begin(), lowerSignature.end(), lowerSignature.begin(), ::tolower);
    
    return lowerCalculated == lowerSignature;
}

std::string HmacService::BytesToHex(const std::vector<unsigned char>& bytes) {
    std::stringstream ss;
    ss << std::hex << std::setfill('0');
    
    for (unsigned char byte : bytes) {
        ss << std::setw(2) << static_cast<int>(byte);
    }
    
    return ss.str();
}

std::vector<unsigned char> HmacService::HexToBytes(const std::string& hex) {
    std::vector<unsigned char> bytes;
    
    for (size_t i = 0; i < hex.length(); i += 2) {
        std::string byteString = hex.substr(i, 2);
        unsigned char byte = static_cast<unsigned char>(std::stoi(byteString, nullptr, 16));
        bytes.push_back(byte);
    }
    
    return bytes;
}

} // namespace hmac_utils 