using Microsoft.AspNetCore.Mvc;
using SampleGameBackend.Database;
using SampleGameBackend.Models;
using SampleGameBackend.Services;
using System.Text.Json;

namespace SampleGameBackend.Controllers
{
    [ApiController]
    [Route("api")]
    public class ValidationController : ControllerBase
    {
        private readonly ValidationService _validationService;
        private readonly DatabaseService _databaseService;
        private readonly ILogger<ValidationController> _logger;

        public ValidationController(ValidationService validationService, DatabaseService databaseService, ILogger<ValidationController> logger)
        {
            _validationService = validationService;
            _databaseService = databaseService;
            _logger = logger;
        }

        [HttpPost("validate")]
        public IActionResult ValidateUserAction([FromBody] ValidateRequest req)
        {
            // TODO: We need a defense logic to prevent duplicate UUIDs in requests.
            try
            {
                // Request validation
                if (req == null || string.IsNullOrEmpty(req.Uuid) || string.IsNullOrEmpty(req.UserSig) || 
                    string.IsNullOrEmpty(req.UserAddress) || string.IsNullOrEmpty(req.ProjectId) || 
                    string.IsNullOrEmpty(req.Digest) || req.Intent == null)
                {
                    return BadRequest(new ValidateResponse
                    {
                        Success = false,
                        ErrorCode = "INVALID_REQUEST"
                    });
                }

                // Intent validation
                if (!_validationService.ValidateIntent(req.Intent))
                {
                    return BadRequest(new ValidateResponse
                    {
                        Success = false,
                        ErrorCode = "INVALID_INTENT"
                    });
                }

                // Get session ID
                var sessionId = HttpContext.Items["SessionId"] as string;
                if (string.IsNullOrEmpty(sessionId))
                {
                    return BadRequest(new ValidateResponse
                    {
                        Success = false,
                        ErrorCode = "INVALID_SESSION"
                    });
                }

                // Store UUID and SessionID mapping
                try
                {
                    _databaseService.StoreUuidMapping(req.Uuid, sessionId);
                }
                catch (Exception ex)
                {
                    _logger.LogError(ex, "Failed to store UUID mapping");
                    return StatusCode(500, new ValidateResponse
                    {
                        Success = false,
                        ErrorCode = "UUID_MAPPING_FAILED"
                    });
                }

                var requestJson = JsonSerializer.Serialize(req);
                _logger.LogInformation("ValidateUserActionHandler: sessionID={SessionId}, uuid={Uuid}, req={Request}", 
                    sessionId, req.Uuid, requestJson);

                // For mint method, validate and deduct assets
                if (req.Intent.Method == "mint" || req.Intent.Method == "transfer")
                {
                    try
                    {
                        _validationService.ValidateAndProcessMint(sessionId, req.Intent.From);
                    }
                    catch (Exception ex)
                    {
                        _logger.LogError(ex, "Insufficient balance for sessionId: {SessionId}", sessionId);
                        return BadRequest(new ValidateResponse
                        {
                            Success = false,
                            ErrorCode = "INSUFFICIENT_BALANCE"
                        });
                    }
                }

                // Generate validator signature
                string validatorSig;
                try
                {
                    validatorSig = _validationService.GenerateValidatorSignature(req.UserSig, req.Digest);
                }
                catch (Exception ex)
                {
                    _logger.LogError(ex, "GenerateValidatorSignature failed");
                    return StatusCode(500, new ValidateResponse
                    {
                        Success = false,
                        ErrorCode = "SIGNATURE_GENERATION"
                    });
                }

                _logger.LogInformation("validateUserActionHandler: validatorSig={ValidatorSig}, userSig={UserSig}, digest={Digest}", 
                    validatorSig, req.UserSig, req.Digest);

                // Success response
                var response = new ValidateResponse
                {
                    Success = true,
                    ErrorCode = null,
                    Data = new ValidateResponseData
                    {
                        UserSig = req.UserSig,
                        ValidatorSig = validatorSig
                    }
                };

                return Ok(response);
            }
            catch (Exception ex)
            {
                _logger.LogError(ex, "ValidateUserAction error");
                return StatusCode(500, new ValidateResponse
                {
                    Success = false,
                    ErrorCode = "INTERNAL_ERROR"
                });
            }
        }
    }
} 