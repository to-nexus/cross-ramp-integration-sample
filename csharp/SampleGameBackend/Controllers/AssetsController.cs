using Microsoft.AspNetCore.Mvc;
using SampleGameBackend.Database;
using SampleGameBackend.Models;
using SampleGameBackend.Middleware;

namespace SampleGameBackend.Controllers
{
    [ApiController]
    [Route("api")]
    public class AssetsController : ControllerBase
    {
        private readonly DatabaseService _databaseService;
        private readonly ILogger<AssetsController> _logger;

        public AssetsController(DatabaseService databaseService, ILogger<AssetsController> logger)
        {
            _databaseService = databaseService;
            _logger = logger;
        }

        [HttpGet("assets")]
        public IActionResult GetAssets([FromQuery] string? language)
        {
            try
            {
                // Validate session ID
                var sessionId = HttpContext.Items["SessionId"] as string;
                if (string.IsNullOrEmpty(sessionId))
                {
                    return BadRequest(new { success = false, errorCode = "INVALID_SESSION" });
                }

                // Validate language parameter
                if (string.IsNullOrEmpty(language))
                {
                    language = "ko"; // Default value
                }

                // Get or create session-specific asset information
                var sessionAssets = _databaseService.GetOrCreateSessionAssets(sessionId);

                // Convert to Asset struct
                var assets = new List<Asset>();
                foreach (var kvp in sessionAssets.Assets)
                {
                    assets.Add(new Asset
                    {
                        Id = kvp.Key,
                        Balance = kvp.Value
                    });
                }

                var v1Data = new V1Data
                {
                    PlayerId = sessionId,
                    Name = $"playerName_{sessionId}",
                    WalletAddress = "0xaaaa",
                    Server = "test",
                    Assets = assets
                };

                // Parse session information
                var createdAt = DateTime.Parse(sessionAssets.CreatedAt);
                var updatedAt = DateTime.Parse(sessionAssets.UpdatedAt);

                var guide = new
                {
                    Authorization = HttpContext.Request.Headers["Authorization"].FirstOrDefault(),
                    DappAuth = HttpContext.Request.Headers["X-Dapp-Authorization"].FirstOrDefault(),
                    SessionID = sessionId,
                    Message = "The guide field displays header information at request time. It is used to verify that the game company and protocol are correctly matched and is not provided to the game company. For ramp frontend developer reference.",
                    SessionInfo = new
                    {
                        CreatedAt = createdAt,
                        UpdatedAt = updatedAt
                    }
                };

                var response = new Response
                {
                    Success = true,
                    ErrorCode = null,
                    Data = new ResponseData
                    {
                        V1 = v1Data,
                        Guide = guide
                    }
                };

                return Ok(response);
            }
            catch (Exception ex)
            {
                _logger.LogError(ex, "GetAssets error for sessionId: {SessionId}", 
                    HttpContext.Items["SessionId"] as string);
                return StatusCode(500, new { success = false, errorCode = "DB_ERROR" });
            }
        }
    }
} 