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
        [JsonPropertyName("project_id")]
        public string ProjectId { get; set; } = string.Empty;

        [JsonPropertyName("pair_id")]
        public uint PairId { get; set; }

        [JsonPropertyName("token")]
        public PairToken Token { get; set; } = new();

        [JsonPropertyName("materials")]
        public List<PairAsset> Materials { get; set; } = new();

        [JsonPropertyName("outputs")]
        public List<PairAsset> Outputs { get; set; } = new();
    }
} 