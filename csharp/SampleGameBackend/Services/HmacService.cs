using System.Security.Cryptography;
using System.Text;
using System.Text.Json;

namespace SampleGameBackend.Services
{
    public class HmacService
    {
        private const string HMACSalt = "my_secret_salt_value_!@#$%^&*";

        // Base64 URL decoding function (as per guide specification)
        private static byte[] Base64UrlDecode(string str)
        {
            // Convert URL safe base64 to standard base64
            str = str.Replace('-', '+').Replace('_', '/');
            // Add padding (if needed)
            while (str.Length % 4 != 0)
            {
                str += '=';
            }
            return Convert.FromBase64String(str);
        }

        public static string GenerateHmac(string requestBody)
        {
            var bodyBytes = Encoding.UTF8.GetBytes(requestBody);
            
            // Use Base64 URL decoding as per guide
            var saltBytes = Base64UrlDecode(HMACSalt);
            using var hmac = new HMACSHA256(saltBytes);
            var hashBytes = hmac.ComputeHash(bodyBytes);
            return Convert.ToHexString(hashBytes).ToLower();
        }

        public static bool ValidateHmac(string requestBody, string hmacSignature)
        {
            if (string.IsNullOrEmpty(hmacSignature))
                return false;

            var calculatedHmac = GenerateHmac(requestBody);
            return string.Equals(calculatedHmac, hmacSignature, StringComparison.OrdinalIgnoreCase);
        }
    }
} 