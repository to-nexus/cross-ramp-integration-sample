using Microsoft.AspNetCore.Mvc;
using SampleGameBackend.Database;
using SampleGameBackend.Models;
using SampleGameBackend.Services;

namespace SampleGameBackend.Controllers
{
    [ApiController]
    [Route("api")]
    public class ExchangeController : ControllerBase
    {
        private readonly ExchangeService _exchangeService;
        private readonly DatabaseService _databaseService;
        private readonly ILogger<ExchangeController> _logger;

        public ExchangeController(ExchangeService exchangeService, DatabaseService databaseService, ILogger<ExchangeController> logger)
        {
            _exchangeService = exchangeService;
            _databaseService = databaseService;
            _logger = logger;
        }

        [HttpPost("result")]
        public IActionResult ExchangeResult([FromBody] ExchangeResp req)
        {
            try
            {
                // Request validation
                if (req == null)
                {
                    _logger.LogError("Failed to bind request body");
                    return BadRequest(new { error = "Failed to bind request body" });
                }

                // Log request body
                _logger.LogInformation("ResultHandler: requestBody={RequestBody}", req);

                if (req.Intent.Outputs.Count > 0)
                {
                    // Get SessionID by UUID
                    string sessionId;
                    try
                    {
                        sessionId = _databaseService.GetSessionIdByUuid(req.Uuid);
                    }
                    catch (Exception ex)
                    {
                        _logger.LogError(ex, "Failed to get session ID by UUID: {Uuid}", req.Uuid);
                        return BadRequest(new { error = "Invalid UUID or session not found" });
                    }

                    // Process exchange result
                    try
                    {
                        // Note: In C# version, we'll use a simple receipt status
                        // In real implementation, you would extract this from the receipt
                        var receiptStatus = 1u; // Assuming success
                        _exchangeService.ProcessExchangeResult(sessionId, req.Intent.Outputs, receiptStatus);
                    }
                    catch (Exception ex)
                    {
                        _logger.LogError(ex, "Failed to process exchange result");
                        return StatusCode(500, new { error = "Failed to process exchange result" });
                    }
                }

                return Ok(new { success = true });
            }
            catch (Exception ex)
            {
                _logger.LogError(ex, "ExchangeResult error");
                return StatusCode(500, new { error = "Internal server error" });
            }
        }
    }
} 