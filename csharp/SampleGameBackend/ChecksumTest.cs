using System;
using System.Text.RegularExpressions;
using Nethereum.Util;
using Nethereum.Hex.HexConvertors.Extensions;

namespace SampleGameBackend
{
    /// <summary>
    /// EVM Address validation utility and test class
    /// </summary>
    public class ChecksumTest
    {
        /// <summary>
        /// Validates if a string is a valid EVM address
        /// </summary>
        /// <param name="address">Address string to validate</param>
        /// <returns>True if valid EVM address, false otherwise</returns>
        public static bool IsValidEVMAddress(string address)
        {
            if (string.IsNullOrEmpty(address))
                return false;

            // Remove 0x prefix if present
            var cleanAddress = address.StartsWith("0x", StringComparison.OrdinalIgnoreCase) 
                ? address.Substring(2) 
                : address;

            // Check length (should be 40 hex characters)
            if (cleanAddress.Length != 40)
                return false;

            // Check if all characters are valid hexadecimal
            return Regex.IsMatch(cleanAddress, "^[0-9A-Fa-f]+$");
        }

        /// <summary>
        /// Validates if address has valid EIP-55 checksum using Nethereum
        /// </summary>
        /// <param name="address">Address string to validate</param>
        /// <returns>True if valid checksum address, false otherwise</returns>
        public static bool IsValidChecksumAddress(string address)
        {
            if (!IsValidEVMAddress(address))
                return false;

            try
            {
                // Use Nethereum's AddressUtil to validate checksum
                var addressUtil = new AddressUtil();
                
                // Ensure address has 0x prefix for Nethereum
                var formattedAddress = address.StartsWith("0x", StringComparison.OrdinalIgnoreCase) 
                    ? address 
                    : "0x" + address;

                // Check if the address is valid and has correct checksum
                return addressUtil.IsValidEthereumAddressHexFormat(formattedAddress) &&
                       addressUtil.IsChecksumAddress(formattedAddress);
            }
            catch
            {
                return false;
            }
        }

        /// <summary>
        /// Converts address to checksum format using Nethereum
        /// </summary>
        /// <param name="address">Address to convert</param>
        /// <returns>Checksum formatted address or null if invalid</returns>
        public static string? ToChecksumAddress(string address)
        {
            if (!IsValidEVMAddress(address))
                return null;

            try
            {
                var addressUtil = new AddressUtil();
                var formattedAddress = address.StartsWith("0x", StringComparison.OrdinalIgnoreCase) 
                    ? address 
                    : "0x" + address;

                return addressUtil.ConvertToChecksumAddress(formattedAddress);
            }
            catch
            {
                return null;
            }
        }

        /// <summary>
        /// Checks if two addresses are the same (case-insensitive)
        /// </summary>
        /// <param name="address1">First address</param>
        /// <param name="address2">Second address</param>
        /// <returns>True if addresses are the same</returns>
        public static bool AreAddressesEqual(string address1, string address2)
        {
            if (!IsValidEVMAddress(address1) || !IsValidEVMAddress(address2))
                return false;

            var clean1 = address1.StartsWith("0x", StringComparison.OrdinalIgnoreCase) 
                ? address1.Substring(2) 
                : address1;
            var clean2 = address2.StartsWith("0x", StringComparison.OrdinalIgnoreCase) 
                ? address2.Substring(2) 
                : address2;

            return string.Equals(clean1, clean2, StringComparison.OrdinalIgnoreCase);
        }

        /// <summary>
        /// Checks if address is the zero address
        /// </summary>
        /// <param name="address">Address to check</param>
        /// <returns>True if zero address</returns>
        public static bool IsZeroAddress(string address)
        {
            return AreAddressesEqual(address, "0x0000000000000000000000000000000000000000");
        }

        /// <summary>
        /// Manual test method to run all validation tests
        /// </summary>
        public static void RunTests()
        {
            Console.WriteLine("=== EVM Address Validation Tests ===\n");

            // Test basic address validation
            TestBasicAddressValidation();

            // Test checksum validation
            TestChecksumValidation();

            // Test address operations
            TestAddressOperations();

            Console.WriteLine("=== All Tests Completed ===");
        }

        private static void TestBasicAddressValidation()
        {
            Console.WriteLine("--- Basic Address Validation Tests ---");

            var testCases = new[]
            {
                // Valid addresses
                new { Address = "0x742C4D61c2a093537BDBdea15Ec35b3C4094d0e3", Expected = true, Name = "Valid address with 0x prefix" },
                new { Address = "742C4D61c2a093537BDBdea15Ec35b3C4094d0e3", Expected = true, Name = "Valid address without 0x prefix" },
                new { Address = "0x742c4d61c2a093537bdbdea15ec35b3c4094d0e3", Expected = true, Name = "Valid address all lowercase" },
                new { Address = "0x742C4D61C2A093537BDBDEA15EC35B3C4094D0E3", Expected = true, Name = "Valid address all uppercase" },
                new { Address = "0x0000000000000000000000000000000000000000", Expected = true, Name = "Zero address" },

                // Invalid addresses
                new { Address = "", Expected = false, Name = "Empty string" },
                new { Address = "0x742C4D61c2a093537BDBdea15Ec35b3C4094d0e", Expected = false, Name = "Too short" },
                new { Address = "0x742C4D61c2a093537BDBdea15Ec35b3C4094d0e33a", Expected = false, Name = "Too long" },
                new { Address = "0x742C4D61c2a093537BDBdea15Ec35b3C4094d0g3", Expected = false, Name = "Invalid hex characters" },
                new { Address = "0x", Expected = false, Name = "Only 0x prefix" },
                new { Address = "0x742C4D61c2a093537BDBdXa15Ec35b3C4094d0e3", Expected = false, Name = "Invalid characters in middle" }
            };

            foreach (var test in testCases)
            {
                var result = IsValidEVMAddress(test.Address);
                var status = result == test.Expected ? "✅ PASS" : "❌ FAIL";
                Console.WriteLine($"{status} {test.Name}: {test.Address} -> {result}");
            }

            Console.WriteLine();
        }

        private static void TestChecksumValidation()
        {
            Console.WriteLine("--- Checksum Validation Tests ---");

            var testCases = new[]
            {
                // Valid checksum addresses (EIP-55)
                new { Address = "0x5aAeb6053F3E94C9b9A09f33669435E7Ef1BeAed", Expected = true, Name = "Valid checksum address" },
                new { Address = "0xfB6916095ca1df60bB79Ce92cE3Ea74c37c5d359", Expected = true, Name = "Valid checksum address 2" },
                new { Address = "0xdbF03B407c01E7cD3CBea99509d93f8DDDC8C6FB", Expected = true, Name = "Valid checksum address 3" },
                new { Address = "0xD1220A0cf47c7B9Be7A2E6BA89F429762e7b9aDb", Expected = true, Name = "Valid checksum address 4" },

                // Invalid checksum (wrong case)
                new { Address = "0x5aaeb6053f3e94c9b9a09f33669435e7ef1beaed", Expected = false, Name = "Invalid checksum - all lowercase" },
                new { Address = "0x5AAEB6053F3E94C9B9A09F33669435E7EF1BEAED", Expected = false, Name = "Invalid checksum - all uppercase" },
                new { Address = "0x5aAeb6053f3e94c9b9a09f33669435e7ef1beaed", Expected = false, Name = "Invalid checksum - mixed wrong case" },

                // Invalid addresses
                new { Address = "invalid-address", Expected = false, Name = "Invalid address format" },
                new { Address = "", Expected = false, Name = "Empty string" }
            };

            foreach (var test in testCases)
            {
                var result = IsValidChecksumAddress(test.Address);
                var status = result == test.Expected ? "✅ PASS" : "❌ FAIL";
                Console.WriteLine($"{status} {test.Name}: {test.Address} -> {result}");
            }

            Console.WriteLine();
        }

        private static void TestAddressOperations()
        {
            Console.WriteLine("--- Address Operations Tests ---");

            // Test address conversion
            var originalAddress = "0x5aAeb6053F3E94C9b9A09f33669435E7Ef1BeAed";
            var checksumAddress = ToChecksumAddress(originalAddress.ToLower());
            var conversionTest = checksumAddress == originalAddress;
            Console.WriteLine($"{(conversionTest ? "✅ PASS" : "❌ FAIL")} Address conversion: {originalAddress.ToLower()} -> {checksumAddress}");

            // Test address equality
            var addr1 = "0x5aAeb6053F3E94C9b9A09f33669435E7Ef1BeAed";
            var addr2 = "0x5aaeb6053f3e94c9b9a09f33669435e7ef1beaed";
            var equalityTest = AreAddressesEqual(addr1, addr2);
            Console.WriteLine($"{(equalityTest ? "✅ PASS" : "❌ FAIL")} Address equality (case-insensitive): {equalityTest}");

            // Test zero address
            var zeroAddr = "0x0000000000000000000000000000000000000000";
            var nonZeroAddr = "0x5aAeb6053F3E94C9b9A09f33669435E7Ef1BeAed";
            var zeroTest1 = IsZeroAddress(zeroAddr);
            var zeroTest2 = !IsZeroAddress(nonZeroAddr);
            Console.WriteLine($"{(zeroTest1 ? "✅ PASS" : "❌ FAIL")} Zero address detection: {zeroAddr} -> {zeroTest1}");
            Console.WriteLine($"{(zeroTest2 ? "✅ PASS" : "❌ FAIL")} Non-zero address detection: {nonZeroAddr} -> {!IsZeroAddress(nonZeroAddr)}");

            Console.WriteLine();
        }
    }
}
