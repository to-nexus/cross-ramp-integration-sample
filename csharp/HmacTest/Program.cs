using System.Security.Cryptography;
using System.Text;
using System.Text.Json;

namespace HmacTest
{
    public class Program
    {
        // TODO: HMAC salt - In actual implementation, load from environment variables or configuration file
        private const string Salt = "my_secret_salt_value_!@#$%^&*"; // hmac key

        public class Body
        {
            public int UserId { get; set; }
            public string Username { get; set; } = "";
            public string Email { get; set; } = "";
            public string Role { get; set; } = "";
            public int CreatedAt { get; set; }
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

            // Use hardcoded JSON string to match other languages
            var jsonString = "{\"userId\":1234,\"username\":\"홍길동\",\"email\":\"user@example.com\",\"role\":\"admin\",\"createdAt\":1234567890}";
            Console.WriteLine(jsonString);

            var bodyBytes = Encoding.UTF8.GetBytes(jsonString);
            using var hmac = new HMACSHA256(Encoding.UTF8.GetBytes(Salt));
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
