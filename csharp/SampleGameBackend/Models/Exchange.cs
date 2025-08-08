using System.Text.Json.Serialization;

namespace SampleGameBackend.Models
{
    public class ExchangeResultRequest
    {
        [JsonPropertyName("uuid")]
        public string Uuid { get; set; } = string.Empty;

        [JsonPropertyName("tx_hash")]
        public string TxHash { get; set; } = string.Empty;

        [JsonPropertyName("receipt")]
        public object Receipt { get; set; } = new();

        [JsonPropertyName("intent")]
        public ExchangeIntent Intent { get; set; } = new();
    }

    public class PairToken
    {
        [JsonPropertyName("id")]
        public string TokenId { get; set; } = string.Empty;

        [JsonPropertyName("amount")]
        public uint Amount { get; set; }
    }

    public class PairAsset
    {
        [JsonPropertyName("id")]
        public string AssetId { get; set; } = string.Empty;

        [JsonPropertyName("amount")]
        public uint Amount { get; set; }
    }

    public class ExchangeIntent
    {
        [JsonPropertyName("type")]
        public string Type { get; set; } = string.Empty;

        [JsonPropertyName("method")]
        public string Method { get; set; } = string.Empty;

        [JsonPropertyName("from")]
        public List<PairAsset> From { get; set; } = new();

        [JsonPropertyName("to")]
        public List<PairAsset> To { get; set; } = new();
    }
} 