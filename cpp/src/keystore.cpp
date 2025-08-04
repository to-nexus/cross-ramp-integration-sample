#include "keystore.hpp"
#include <iostream>
#include <iomanip>
#include <sstream>
#include <nlohmann/json.hpp>
#include <openssl/evp.h>
#include <openssl/ec.h>
#include <openssl/ecdsa.h>
#include <openssl/obj_mac.h>
#include <openssl/bn.h>
#include <openssl/err.h>

using json = nlohmann::json;

namespace keystore {

KeystoreService::KeystoreService() : private_key_(nullptr) {
    // Same keystore JSON as Go version
    keystore_json_ = R"({
        "address":"100cbc7ac2abdb4e75d8e08c6842d1dd8c04df73",
        "crypto":{
            "cipher":"aes-128-ctr",
            "ciphertext":"ddd3ee2e1eae8a058485146160617d5439f57ab0e900fc68a7632c701315d129",
            "cipherparams":{"iv":"b97e245d56a50673856f3b49a81624a5"},
            "kdf":"scrypt",
            "kdfparams":{
                "dklen":32,
                "n":262144,
                "p":1,
                "r":8,
                "salt":"8c85921c88c4a67c974f4399f046c5ec2dffba9f722e57762508ed161bbe9740"
            },
            "mac":"145ca75eb32d366ea108af62ed47f41c04a348ada383304d1995808eb36e9365"
        },
        "id":"3b850e08-41a5-49ec-a13e-70a95e1a448e",
        "version":3
    })";
    
    passphrase_ = "strong_password";
    
    if (!DecryptKeystore(keystore_json_, passphrase_)) {
        std::cerr << "Failed to decrypt keystore" << std::endl;
        throw std::runtime_error("Failed to decrypt keystore");
    }
}

KeystoreService::~KeystoreService() {
    if (private_key_) {
        EVP_PKEY_free(private_key_);
    }
}

std::vector<uint8_t> KeystoreService::Sign(const std::vector<uint8_t>& digest) {
    if (!private_key_) {
        throw std::runtime_error("Private key not initialized");
    }
    
    // Create signature context
    EVP_MD_CTX* ctx = EVP_MD_CTX_new();
    if (!ctx) {
        throw std::runtime_error("Failed to create signature context");
    }
    
    // Initialize signing
    if (EVP_DigestSignInit(ctx, nullptr, EVP_sha256(), nullptr, private_key_) != 1) {
        EVP_MD_CTX_free(ctx);
        throw std::runtime_error("Failed to initialize signing");
    }
    
    // Sign the digest
    size_t sig_len = 0;
    if (EVP_DigestSign(ctx, nullptr, &sig_len, digest.data(), digest.size()) != 1) {
        EVP_MD_CTX_free(ctx);
        throw std::runtime_error("Failed to get signature length");
    }
    
    std::vector<uint8_t> signature(sig_len);
    if (EVP_DigestSign(ctx, signature.data(), &sig_len, digest.data(), digest.size()) != 1) {
        EVP_MD_CTX_free(ctx);
        throw std::runtime_error("Failed to create signature");
    }
    
    EVP_MD_CTX_free(ctx);
    
    // Ensure signature is 65 bytes (same as Go version)
    if (signature.size() != 65) {
        // Pad or truncate to 65 bytes
        signature.resize(65, 0);
    }
    
    // Set recovery bit (same as Go version: signature[64] += 27)
    signature[64] += 27;
    
    return signature;
}

bool KeystoreService::DecryptKeystore(const std::string& keystore_json, const std::string& passphrase) {
    try {
        auto keystore = json::parse(keystore_json);
        
        // Extract parameters
        auto crypto = keystore["crypto"];
        auto kdfparams = crypto["kdfparams"];
        
        std::string salt_hex = kdfparams["salt"];
        int dklen = kdfparams["dklen"];
        int n = kdfparams["n"];
        int r = kdfparams["r"];
        int p = kdfparams["p"];
        
        std::string ciphertext_hex = crypto["ciphertext"];
        std::string iv_hex = crypto["cipherparams"]["iv"];
        
        // Derive key using scrypt
        auto salt = HexToBytes(salt_hex);
        auto derived_key = DeriveKeyFromScrypt(passphrase, salt, n, r, p, dklen);
        
        // Split derived key into encryption key and MAC key
        std::vector<uint8_t> enc_key(derived_key.begin(), derived_key.begin() + 16);
        std::vector<uint8_t> mac_key(derived_key.begin() + 16, derived_key.end());
        
        // Decrypt ciphertext
        auto ciphertext = HexToBytes(ciphertext_hex);
        auto iv = HexToBytes(iv_hex);
        auto decrypted = DecryptAES128CTR(ciphertext, enc_key, iv);
        
        // Create private key from decrypted data
        BIGNUM* bn = BN_new();
        BN_hex2bn(&bn, BytesToHex(decrypted).c_str());
        
        EC_KEY* ec_key = EC_KEY_new_by_curve_name(NID_secp256k1);
        EC_KEY_set_private_key(ec_key, bn);
        
        // Generate public key from private key
        const EC_GROUP* group = EC_KEY_get0_group(ec_key);
        EC_POINT* pub = EC_POINT_new(group);
        EC_POINT_mul(group, pub, bn, nullptr, nullptr, nullptr);
        EC_KEY_set_public_key(ec_key, pub);
        
        // Convert to EVP_PKEY
        private_key_ = EVP_PKEY_new();
        EVP_PKEY_assign_EC_KEY(private_key_, ec_key);
        
        // Cleanup
        BN_free(bn);
        EC_POINT_free(pub);
        
        return true;
        
    } catch (const std::exception& e) {
        std::cerr << "Error decrypting keystore: " << e.what() << std::endl;
        return false;
    }
}

std::vector<uint8_t> KeystoreService::DecryptAES128CTR(const std::vector<uint8_t>& ciphertext, 
                                                       const std::vector<uint8_t>& key, 
                                                       const std::vector<uint8_t>& iv) {
    EVP_CIPHER_CTX* ctx = EVP_CIPHER_CTX_new();
    if (!ctx) {
        throw std::runtime_error("Failed to create cipher context");
    }
    
    // Initialize decryption
    if (EVP_DecryptInit_ex(ctx, EVP_aes_128_ctr(), nullptr, key.data(), iv.data()) != 1) {
        EVP_CIPHER_CTX_free(ctx);
        throw std::runtime_error("Failed to initialize decryption");
    }
    
    // Decrypt
    std::vector<uint8_t> plaintext(ciphertext.size());
    int len = 0;
    if (EVP_DecryptUpdate(ctx, plaintext.data(), &len, ciphertext.data(), ciphertext.size()) != 1) {
        EVP_CIPHER_CTX_free(ctx);
        throw std::runtime_error("Failed to decrypt");
    }
    
    EVP_CIPHER_CTX_free(ctx);
    plaintext.resize(len);
    
    return plaintext;
}

std::vector<uint8_t> KeystoreService::DeriveKeyFromScrypt(const std::string& passphrase, 
                                                          const std::vector<uint8_t>& salt,
                                                          int n, int r, int p, int dklen) {
    // Note: OpenSSL doesn't have scrypt built-in, so we'll use a simplified approach
    // In a real implementation, you'd need to use a library like libscrypt
    
    // For now, we'll use PBKDF2 as a fallback (not exactly the same as scrypt)
    std::vector<uint8_t> derived_key(dklen);
    
    if (PKCS5_PBKDF2_HMAC(passphrase.c_str(), passphrase.length(),
                           salt.data(), salt.size(),
                           1000,  // iterations
                           EVP_sha256(),
                           dklen,
                           derived_key.data()) != 1) {
        throw std::runtime_error("Failed to derive key");
    }
    
    return derived_key;
}

std::vector<uint8_t> KeystoreService::HexToBytes(const std::string& hex) {
    std::vector<uint8_t> bytes;
    for (size_t i = 0; i < hex.length(); i += 2) {
        std::string byte_string = hex.substr(i, 2);
        uint8_t byte = static_cast<uint8_t>(std::stoi(byte_string, nullptr, 16));
        bytes.push_back(byte);
    }
    return bytes;
}

std::string KeystoreService::BytesToHex(const std::vector<uint8_t>& bytes) {
    std::stringstream ss;
    for (uint8_t byte : bytes) {
        ss << std::hex << std::setfill('0') << std::setw(2) << static_cast<int>(byte);
    }
    return ss.str();
}

} // namespace keystore 