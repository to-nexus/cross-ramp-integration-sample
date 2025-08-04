using System.Text.Json.Serialization;

namespace SampleGameBackend.Models
{
    public class V1Data
    {
        [JsonPropertyName("player_id")]
        public string PlayerId { get; set; } = string.Empty;

        [JsonPropertyName("name")]
        public string Name { get; set; } = string.Empty;

        [JsonPropertyName("wallet_address")]
        public string WalletAddress { get; set; } = string.Empty;

        [JsonPropertyName("server")]
        public string Server { get; set; } = string.Empty;

        [JsonPropertyName("assets")]
        public List<Asset> Assets { get; set; } = new();
    }

    public class Response
    {
        [JsonPropertyName("success")]
        public bool Success { get; set; }

        [JsonPropertyName("errorCode")]
        public string? ErrorCode { get; set; }

        [JsonPropertyName("data")]
        public ResponseData Data { get; set; } = new();
    }

    public class ResponseData
    {
        [JsonPropertyName("v1")]
        public V1Data V1 { get; set; } = new();

        [JsonPropertyName("guide")]
        public object Guide { get; set; } = new();
    }

    public class ValidateRequest
    {
        [JsonPropertyName("uuid")]
        public string Uuid { get; set; } = string.Empty;

        [JsonPropertyName("user_sig")]
        public string UserSig { get; set; } = string.Empty;

        [JsonPropertyName("user_address")]
        public string UserAddress { get; set; } = string.Empty;

        [JsonPropertyName("project_id")]
        public string ProjectId { get; set; } = string.Empty;

        [JsonPropertyName("digest")]
        public string Digest { get; set; } = string.Empty;

        [JsonPropertyName("intent")]
        public ValidateIntent Intent { get; set; } = new();
    }

    public class ValidateIntent
    {
        [JsonPropertyName("method")]
        public string Method { get; set; } = string.Empty;

        [JsonPropertyName("from")]
        public List<ValidateAsset> From { get; set; } = new();

        [JsonPropertyName("to")]
        public List<ValidateAsset> To { get; set; } = new();
    }

    public class ValidateAsset
    {
        [JsonPropertyName("type")]
        public string Type { get; set; } = string.Empty;

        [JsonPropertyName("id")]
        public string Id { get; set; } = string.Empty;

        [JsonPropertyName("amount")]
        public int Amount { get; set; }
    }

    public class ValidateResponse
    {
        [JsonPropertyName("success")]
        public bool Success { get; set; }

        [JsonPropertyName("errorCode")]
        public string? ErrorCode { get; set; }

        [JsonPropertyName("data")]
        public ValidateResponseData Data { get; set; } = new();
    }

    public class ValidateResponseData
    {
        [JsonPropertyName("userSig")]
        public string UserSig { get; set; } = string.Empty;

        [JsonPropertyName("validatorSig")]
        public string ValidatorSig { get; set; } = string.Empty;
    }
} 