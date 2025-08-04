using System.Text.Json.Serialization;

namespace SampleGameBackend.Models
{
    public class Asset
    {
        [JsonPropertyName("id")]
        public string Id { get; set; } = string.Empty;

        [JsonPropertyName("balance")]
        public string Balance { get; set; } = string.Empty;
    }

    public class SessionAssets
    {
        [JsonPropertyName("session_id")]
        public string SessionId { get; set; } = string.Empty;

        [JsonPropertyName("assets")]
        public Dictionary<string, string> Assets { get; set; } = new();

        [JsonPropertyName("created_at")]
        public string CreatedAt { get; set; } = string.Empty;

        [JsonPropertyName("updated_at")]
        public string UpdatedAt { get; set; } = string.Empty;
    }
} 