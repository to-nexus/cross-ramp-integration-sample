using SampleGameBackend.Database;
using SampleGameBackend.Models;

namespace SampleGameBackend.Services
{
    public class ValidationService
    {
        private readonly KeystoreService _keystoreService;
        private readonly DatabaseService _databaseService;

        public ValidationService(KeystoreService keystoreService, DatabaseService databaseService)
        {
            _keystoreService = keystoreService;
            _databaseService = databaseService;
        }

        public bool ValidateIntent(ValidateIntent intent)
        {
            // Validate allowed methods
            var allowedMethods = new HashSet<string> { "mint", "transfer", "burn", "burn-permit" };

            if (!allowedMethods.Contains(intent.Method))
            {
                return false;
            }

            // Special validation for mint method
            if (intent.Method == "mint")
            {
                // Must have at least one from item
                if (intent.From.Count == 0)
                {
                    return false;
                }

                // All from items must be asset type
                foreach (var from in intent.From)
                {
                    if (from.Type != "asset")
                    {
                        return false;
                    }
                }
            }

            return true;
        }

        public string GenerateValidatorSignature(string userSig, string digest)
        {
            try
            {
                var signature = _keystoreService.Sign(digest);
                return signature;
            }
            catch (Exception ex)
            {
                throw new Exception($"Failed to generate validator signature: {ex.Message}");
            }
        }

        public void ValidateAndProcessMint(string sessionId, List<ValidateAsset> fromAssets)
        {
            // Asset balance validation and deduction
            _databaseService.CheckAndDeductAssets(sessionId, fromAssets);
        }
    }
} 