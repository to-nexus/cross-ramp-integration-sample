using System.Security.Cryptography;
using System.Text;
using System.Text.Json;

namespace SampleGameBackend.Services
{
    public class HmacService
    {
        // TODO: HMAC salt - In actual implementation, load from environment variables or configuration file
    private const string HMACSalt = "my_secret_salt_value_!@#$%^&*";

        public static string GenerateHmac(object data)
        {
            var jsonString = JsonSerializer.Serialize(data);
            var bodyBytes = Encoding.UTF8.GetBytes(jsonString);
            
            using var hmac = new HMACSHA256(Encoding.UTF8.GetBytes(HMACSalt));
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