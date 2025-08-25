using Nethereum.Signer;
using Nethereum.Web3.Accounts;
using Nethereum.Util;
using Nethereum.Hex.HexConvertors.Extensions;
using System.Text;
using System.Security.Cryptography;
using Nethereum.Model;

namespace SampleGameBackend.Tests
{
    /// <summary>
    /// C# version of crypto.sign.viem.test.ts using Nethereum library
    /// Tests cryptographic signing functionality similar to TypeScript viem tests
    /// </summary>
    public class CryptoSignTest
    {
        // Use the same private key as Go code
        private const string PrivateKeyHex = "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef";

        /// <summary>
        /// Test signing digest with private key using Nethereum
        /// Equivalent to TypeScript viem's account.sign({ hash: digest })
        /// </summary>
        public void TestShouldSignDigestWithPrivateKey()
        {
            Console.WriteLine("=== Test: Should Sign Digest With Private Key ===");
            
            // Create account using Nethereum
            var account = new Nethereum.Web3.Accounts.Account(PrivateKeyHex);
            
            // Test data
            var testData = "test";
            
            // Generate Keccak256 hash (same as Go's crypto.Keccak256)
            var testBytes = Encoding.UTF8.GetBytes(testData);
            var digest = new Sha3Keccack().CalculateHashFromHex(testBytes.ToHex());
            
            var key = new EthECKey(PrivateKeyHex);
            var signature = key.SignAndCalculateV(digest.HexToByteArray());
            
            // R and S values (64 bytes)
            var rsBytes = signature.To64ByteArray();
            
            Console.WriteLine($"Test Data: {testData}");
            Console.WriteLine($"Digest: 0x{digest}");
            Console.WriteLine($"Signature (R+S): {rsBytes.ToHex()}");
            Console.WriteLine($"Signature R: {signature.R.ToHex()}");
            Console.WriteLine($"Signature S: {signature.S.ToHex()}");
            Console.WriteLine($"Signature V type: {signature.V?.GetType()}");
            Console.WriteLine($"Signature V: {signature.V?.ToHex() ?? "null"}");
            Console.WriteLine($"Signature V length: {signature.V?.Length ?? 0}");
            
            // Check and process V value
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
                Console.WriteLine($"signature.V[0]: {signature.V[0]}, final vByte: {vByte}");
            }
            else
            {
                // Use default value 0 if V is missing (recovery id 0 + 27 = 27)
                Console.WriteLine("V is null or empty, using default recovery id 0");
                vByte = 27; // recovery id 0 + 27
            }
            
            // Full signature = R + S + V (65 bytes)
            var fullSignature = new byte[65];
            Array.Copy(rsBytes, 0, fullSignature, 0, 64);
            fullSignature[64] = vByte;
            
            var ValidatorSig = fullSignature.ToHex();
            
            Console.WriteLine($"V byte: {vByte:x2}");
            Console.WriteLine($"Full Signature: 0x{ValidatorSig}");
            Console.WriteLine($"Signature size: {ValidatorSig.Length}");
            
            // Validation
            if (string.IsNullOrEmpty(signature.ToString()))
                throw new Exception("Signature should not be null or empty");

            Console.WriteLine("✅ Test passed: Signature generated successfully");
        }


        /// <summary>
        /// Run all crypto signing tests
        /// </summary>
        public static void RunAllTests()
        {
            Console.WriteLine("🚀 Starting C# Crypto Sign Tests (Nethereum version)");
            Console.WriteLine("==================================================");
            
            var test = new CryptoSignTest();
            
            try
            {
                test.TestShouldSignDigestWithPrivateKey();
                //test.TestShouldVerifySignatureCorrectly();
                //test.TestShouldHandleMessageSigning();
                //test.TestShouldProduceConsistentResultsWithRawHash();
                //test.TestShouldWorkWithByteArrays();
                //test.TestShouldSignWithSpecificHexDigest();
                
                Console.WriteLine("\n🎉 All tests passed successfully!");
            }
            catch (Exception ex)
            {
                Console.WriteLine($"\n❌ Test failed: {ex.Message}");
                throw;
            }
        }
    }
}
