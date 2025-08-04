using SampleGameBackend.Database;
using SampleGameBackend.Models;

namespace SampleGameBackend.Services
{
    public class ExchangeService
    {
        private readonly DatabaseService _databaseService;
        private readonly ILogger<ExchangeService> _logger;

        public ExchangeService(DatabaseService databaseService, ILogger<ExchangeService> logger)
        {
            _databaseService = databaseService;
            _logger = logger;
        }

        public void ProcessExchangeResult(string sessionId, List<PairAsset> outputs, uint receiptStatus)
        {
            // Skip processing if receipt status is not 0x1
            if (receiptStatus != 1)
            {
                _logger.LogInformation("ProcessExchangeResult: sessionID={SessionId}, receiptStatus={ReceiptStatus}, action=skipped", 
                    sessionId, receiptStatus);
                return;
            }

            // Skip processing if output is empty
            if (outputs.Count == 0)
            {
                _logger.LogInformation("ProcessExchangeResult: sessionID={SessionId}, outputs=empty, action=skipped", sessionId);
                return;
            }

            // Process asset increase
            try
            {
                _databaseService.AddAssets(sessionId, outputs);
                _logger.LogInformation("ProcessExchangeResult: sessionID={SessionId}, outputs={Outputs}, action=assets_added", 
                    sessionId, outputs);
            }
            catch (Exception ex)
            {
                _logger.LogError(ex, "Failed to add assets for sessionId: {SessionId}", sessionId);
                throw;
            }
        }
    }
} 