using System.Security.Cryptography;
using System.Text;
using System.Text.Json;

namespace HmacTest
{
    public class Program
    {
        private const string Salt = "my_secret_salt_value_!@#$%^&*"; // hmac key

        public class Body
        {
            public int UserId { get; set; }
            public string Username { get; set; } = "";
            public string Email { get; set; } = "";
            public string Role { get; set; } = "";
            public int CreatedAt { get; set; }
        }

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

        public static void TestSha256()
        {
            var body = new Body
            {
                UserId = 1234,
                Username = "홍길동",
                Email = "user@example.com",
                Role = "admin",
                CreatedAt = 1234567890
            };

            var bodyBytes = JsonSerializer.SerializeToUtf8Bytes(body);
            Console.WriteLine(Encoding.UTF8.GetString(bodyBytes));

            // Use Base64 URL decoding as per guide
            var saltBytes = Base64UrlDecode(Salt);
            using var hmac = new HMACSHA256(saltBytes);
            var hashBytes = hmac.ComputeHash(bodyBytes);
            var hashString = Convert.ToHexString(hashBytes).ToLower();

            Console.WriteLine($"hashString: {hashString}"); // expected X-HMAC-Signature: f96cf60394f6b8ad3c6de2d5b2b1d1a540f9529082a8eb9cee405bfbdd9f37a1
        }

        public static void Main(string[] args)
        {
            TestSha256();
        }
    }
}
