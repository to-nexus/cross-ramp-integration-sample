#pragma once

#include <string>
#include <vector>
#include <memory>
#include <openssl/ec.h>
#include <openssl/evp.h>

namespace keystore {

class KeystoreService {
public:
    KeystoreService();
    ~KeystoreService();
    
    std::vector<uint8_t> Sign(const std::vector<uint8_t>& digest);
    
private:
    bool DecryptKeystore(const std::string& keystore_json, const std::string& passphrase);
    std::vector<uint8_t> DecryptAES128CTR(const std::vector<uint8_t>& ciphertext, 
                                          const std::vector<uint8_t>& key, 
                                          const std::vector<uint8_t>& iv);
    std::vector<uint8_t> DeriveKeyFromScrypt(const std::string& passphrase, 
                                             const std::vector<uint8_t>& salt,
                                             int n, int r, int p, int dklen);
    std::vector<uint8_t> HexToBytes(const std::string& hex);
    std::string BytesToHex(const std::vector<uint8_t>& bytes);
    
private:
    EVP_PKEY* private_key_;
    std::string keystore_json_;
    std::string passphrase_;
};

} // namespace keystore 