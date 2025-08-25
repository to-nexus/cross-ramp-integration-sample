using Nethereum.Signer;
using Nethereum.Web3.Accounts;
using Nethereum.KeyStore;
using Nethereum.Util;
using Nethereum.Hex.HexConvertors.Extensions;
using Nethereum.Model;
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
            // Actual keystore JSON example generated from Go code
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
                
                // Output private key, public key, and address logs
                var account = new Nethereum.Web3.Accounts.Account(_privateKey);                
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
                // Decrypt keystore using Nethereum's KeyStoreService
                var keyStoreService = new KeyStoreService();
                var privateKeyBytes = keyStoreService.DecryptKeyStoreFromJson(passphrase, keystoreJson);
                
                // Convert byte[] to hex string
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
                // Generate signature using EthECKey directly (same approach as CryptoSignTest)
                var key = new EthECKey(_privateKey);
                
                // Convert digest to byte array (handle hex string cases)
                byte[] digestBytes;
                if (digest.StartsWith("0x"))
                {
                    digestBytes = digest.HexToByteArray();
                }
                else if (digest.Length == 64) // hex string without 0x prefix
                {
                    digestBytes = Convert.FromHexString(digest);
                }
                else
                {
                    // Convert to UTF-8 bytes for regular strings
                    digestBytes = Encoding.UTF8.GetBytes(digest);
                }
                
                var signature = key.Sign(digestBytes);
                
                // R and S values (64 bytes)
                var rsBytes = signature.To64ByteArray();
                
                // Process V value (same logic as CryptoSignTest)
                byte vByte;
                if (signature.V != null && signature.V.Length > 0)
                {
                    // Check if signature.V[0] is already ethereum standard v value (27 or 28)
                    if (signature.V[0] >= 27)
                    {
                        vByte = signature.V[0]; // Already correct ethereum v value
                    }
                    else
                    {
                        vByte = (byte)(signature.V[0] + 27); // Add 27 since it's recovery id
                    }
                }
                else
                {
                    // Use default value if V is missing (recovery id 0 + 27 = 27)
                    vByte = 27;
                }
                
                // Full signature = R + S + V (65 bytes)
                var fullSignature = new byte[65];
                Array.Copy(rsBytes, 0, fullSignature, 0, 64);
                fullSignature[64] = vByte;
                
                // Return as hex string (without 0x prefix)
                return fullSignature.ToHex();
            }
            catch (Exception ex)
            {
                throw new Exception($"Failed to generate signature: {ex.Message}");
            }
        }

        public string GetAddress()
        {
            var account = new Nethereum.Web3.Accounts.Account(_privateKey);
            return account.Address;
        }
    }
} 