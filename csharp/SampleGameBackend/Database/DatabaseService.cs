using SampleGameBackend.Models;
using System.Collections.Concurrent;

namespace SampleGameBackend.Database
{
    public class DatabaseService
    {
        private static readonly ConcurrentDictionary<string, SessionAssets> _sessionAssets = new();
        private static readonly ConcurrentDictionary<string, string> _uuidMappings = new();

        public SessionAssets GetOrCreateSessionAssets(string sessionId)
        {
            if (_sessionAssets.TryGetValue(sessionId, out var existingAssets))
            {
                return existingAssets;
            }

            // Create new session assets
            var sessionAssets = new SessionAssets
            {
                SessionId = sessionId,
                Assets = GenerateRandomAssets(),
                CreatedAt = DateTime.UtcNow.ToString("yyyy-MM-ddTHH:mm:ssZ"),
                UpdatedAt = DateTime.UtcNow.ToString("yyyy-MM-ddTHH:mm:ssZ")
            };

            _sessionAssets.TryAdd(sessionId, sessionAssets);
            return sessionAssets;
        }

        private Dictionary<string, string> GenerateRandomAssets()
        {
            var random = new Random();
            var assets = new Dictionary<string, string>();
            var baseAmount = 100000000;

            // Generate asset_money randomly (1000 ~ 5000)
            var moneyAmount = random.Next(baseAmount) + 1000;
            assets["asset_money"] = moneyAmount.ToString();

            // Generate asset_gold randomly (500 ~ 3000)
            var goldAmount = random.Next(baseAmount) + 500;
            assets["asset_gold"] = goldAmount.ToString();

            // Generate item_gem randomly (500 ~ 3000)
            var gemAmount = random.Next(baseAmount) + 500;
            assets["item_gem"] = gemAmount.ToString();

            // Generate item_banana randomly (500 ~ 3000)
            var bananaAmount = random.Next(baseAmount) + 500;
            assets["item_banana"] = bananaAmount.ToString();

            // Generate asset_silver randomly (500 ~ 3000)
            var silverAmount = random.Next(baseAmount) + 500;
            assets["asset_silver"] = silverAmount.ToString();

            // Generate item_apple randomly (500 ~ 3000)
            var appleAmount = random.Next(baseAmount) + 500;
            assets["item_apple"] = appleAmount.ToString();

            // Generate item_fish randomly (500 ~ 3000)
            var fishAmount = random.Next(baseAmount) + 500;
            assets["item_fish"] = fishAmount.ToString();

            // Generate item_branch randomly (500 ~ 3000)
            var branchAmount = random.Next(baseAmount) + 500;
            assets["item_branch"] = branchAmount.ToString();

            // Generate item_horn randomly (500 ~ 3000)
            var hornAmount = random.Next(baseAmount) + 500;
            assets["item_horn"] = hornAmount.ToString();

            // Generate item_maple randomly (500 ~ 3000)
            var mapleAmount = random.Next(baseAmount) + 500;
            assets["item_maple"] = mapleAmount.ToString();

            return assets;
        }

        public void CheckAndDeductAssets(string sessionId, List<ValidateAsset> fromAssets)
        {
            var sessionAssets = GetOrCreateSessionAssets(sessionId);

            // Validate and deduct balance for each asset
            foreach (var asset in fromAssets)
            {
                if (!sessionAssets.Assets.ContainsKey(asset.Id))
                {
                    throw new Exception($"Asset {asset.Id} not found in session");
                }

                var currentBalance = sessionAssets.Assets[asset.Id];
                if (!int.TryParse(currentBalance, out var currentAmount))
                {
                    throw new Exception($"Invalid balance format for asset {asset.Id}");
                }

                // Validate balance
                if (currentAmount < asset.Amount)
                {
                    throw new Exception($"Insufficient balance for asset {asset.Id}: required {asset.Amount}, available {currentAmount}");
                }

                // Deduct
                var newBalance = currentAmount - asset.Amount;
                sessionAssets.Assets[asset.Id] = newBalance.ToString();
            }

            // Set update time
            sessionAssets.UpdatedAt = DateTime.UtcNow.ToString("yyyy-MM-ddTHH:mm:ssZ");
        }

        public void AddAssets(string sessionId, List<PairAsset> assets)
        {
            var sessionAssets = GetOrCreateSessionAssets(sessionId);

            // Increase balance for each asset
            foreach (var asset in assets)
            {
                if (sessionAssets.Assets.ContainsKey(asset.AssetId))
                {
                    // Add to existing balance
                    var currentBalance = sessionAssets.Assets[asset.AssetId];
                    if (ulong.TryParse(currentBalance, out var currentAmount))
                    {
                        var newBalance = currentAmount + asset.Amount;
                        sessionAssets.Assets[asset.AssetId] = newBalance.ToString();
                    }
                }
                else
                {
                    // Create new asset if it doesn't exist
                    sessionAssets.Assets[asset.AssetId] = asset.Amount.ToString();
                }
            }

            // Set update time
            sessionAssets.UpdatedAt = DateTime.UtcNow.ToString("yyyy-MM-ddTHH:mm:ssZ");
        }

        public void StoreUuidMapping(string uuid, string sessionId)
        {
            _uuidMappings.TryAdd(uuid, sessionId);
            Console.WriteLine($"StoreUUIDMapping: uuid={uuid}, sessionID={sessionId}, action=committed");
        }

        public string GetSessionIdByUuid(string uuid)
        {
            if (_uuidMappings.TryGetValue(uuid, out var sessionId))
            {
                Console.WriteLine($"GetSessionIDByUUID: uuid={uuid}, sessionID={sessionId}");
                return sessionId;
            }

            Console.WriteLine($"GetSessionIDByUUID: warning=UUID not found, uuid={uuid}");
            throw new Exception($"UUID mapping not found: {uuid}");
        }
    }
} 