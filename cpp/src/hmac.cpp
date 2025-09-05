#include "hmac.hpp"
#include <cstring>
#include <algorithm>
#include <openssl/bio.h>
#include <openssl/evp.h>

namespace hmac_utils {

const std::string HmacService::HMAC_SALT = "my_secret_salt_value_!@#$%^&*";

std::vector<unsigned char> HmacService::Base64UrlDecode(const std::string& str) {
    std::string modified = str;
    
    // Convert URL safe base64 to standard base64
    std::replace(modified.begin(), modified.end(), '-', '+');
    std::replace(modified.begin(), modified.end(), '_', '/');
    
    // Add padding (if needed)
    while (modified.length() % 4 != 0) {
        modified += '=';
    }
    
    // Decode base64
    BIO* bio = BIO_new_mem_buf(modified.c_str(), modified.length());
    BIO* b64 = BIO_new(BIO_f_base64());
    BIO_set_flags(b64, BIO_FLAGS_BASE64_NO_NL);
    bio = BIO_push(b64, bio);
    
    std::vector<unsigned char> decoded(modified.length());
    int decodedLength = BIO_read(bio, decoded.data(), modified.length());
    decoded.resize(decodedLength);
    
    BIO_free_all(bio);
    return decoded;
}

std::string HmacService::GenerateHmac(const std::string& data) {
    std::vector<unsigned char> hmac(SHA256_DIGEST_LENGTH);
    
    // Use Base64 URL decoding as per guide
    std::vector<unsigned char> saltBytes = Base64UrlDecode(HMAC_SALT);
    
    HMAC(EVP_sha256(), 
         saltBytes.data(), 
         saltBytes.size(),
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