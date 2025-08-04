using Nethereum.Signer;
using Nethereum.Web3.Accounts;
using Nethereum.KeyStore;
using System.Security.Cryptography;
using System.Text;
using System.Text.Json;

namespace SampleGameBackend.Services
{
    public class KeystoreService
    {
        private readonly string _privateKey;

        public KeystoreService()
        {
            // Go 코드에서 생성한 실제 keystore JSON
            // private key: 3c9817e3bdaca815773de4bc170e464c036149091783b44469b20abef7a31071
            var keyStore = @"{
                ""address"": ""0x100cbc7ac2abdb4e75d8e08c6842d1dd8c04df73"",
                ""crypto"": {
                    ""cipher"": ""aes-128-ctr"",
                    ""ciphertext"": ""7afecf4973b2c0827c11694dfd4190e81148603a17c4ff35963ad9a0eac7217d"",
                    ""cipherparams"": {
                        ""iv"": ""8a22d794e033ac63e3ce3f28848c79d6""
                    },
                    ""kdf"": ""scrypt"",
                    ""kdfparams"": {
                        ""dklen"": 32,
                        ""n"": 262144,
                        ""p"": 1,
                        ""r"": 8,
                        ""salt"": ""6c2c6fc03712e03c59a14eb39846190998a8e630064e9d6d049e39e7e5b0c3bb""
                    },
                    ""mac"": ""d0767683a99721e6b5184ddcd1fa42d705ac78987ca4516a74723ae67dbf43a8""
                },
                ""id"": ""cde2f57e-60a4-4bad-a160-f930e412e1a9"",
                ""version"": 3
            }";
            var passphrase = "strong_password";

            try
            {
                _privateKey = DecryptKeystore(keyStore, passphrase);
                Console.WriteLine("✅ Keystore decryption successful");
                
                // 개인키, 공개키, 주소 로그 출력
                var account = new Account(_privateKey);
                Console.WriteLine($"🔑 Private Key: {_privateKey}");
                Console.WriteLine($"🔑 Public Key: {account.PublicKey}");
                Console.WriteLine($"🔑 Address: {account.Address}");
            }
            catch (Exception ex)
            {
                throw new Exception($"Failed to load keystore: {ex.Message}");
            }

            Console.WriteLine($"✅ Keystore loaded successfully: {_privateKey}");
        }

        private string DecryptKeystore(string keystoreJson, string passphrase)
        {
            try
            {
                // Nethereum의 KeyStoreService를 사용하여 실제 keystore 디크립션
                var keyStoreService = new KeyStoreService();
                var privateKeyBytes = keyStoreService.DecryptKeyStoreFromJson(passphrase, keystoreJson);
                
                // byte[]를 hex string으로 변환
                return "0x" + Convert.ToHexString(privateKeyBytes).ToLower();
            }
            catch (Exception ex)
            {
                throw new Exception($"Keystore decryption failed: {ex.Message}");
            }
        }

        public string Sign(string digest)
        {
            try
            {
                var account = new Account(_privateKey);
                var signer = new EthereumMessageSigner();
                
                // Go 버전과 동일한 서명 생성
                var signature = signer.Sign(Encoding.UTF8.GetBytes(digest), _privateKey);
                
                // v 값을 27로 조정 (Go 버전과 동일)
                var signatureBytes = Convert.FromHexString(signature);
                signatureBytes[64] += 27;
                
                return Convert.ToHexString(signatureBytes);
            }
            catch (Exception ex)
            {
                throw new Exception($"Failed to generate signature: {ex.Message}");
            }
        }

        public string GetAddress()
        {
            var account = new Account(_privateKey);
            return account.Address;
        }
    }
} 