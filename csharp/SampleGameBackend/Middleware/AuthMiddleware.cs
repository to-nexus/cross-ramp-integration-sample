using Microsoft.AspNetCore.Http;

namespace SampleGameBackend.Middleware
{
    public class AuthMiddleware
    {
        private readonly RequestDelegate _next;
        private readonly ILogger<AuthMiddleware> _logger;

        public AuthMiddleware(RequestDelegate next, ILogger<AuthMiddleware> logger)
        {
            _next = next;
            _logger = logger;
        }

        public async Task InvokeAsync(HttpContext context)
        {
            var authHeader = context.Request.Headers["Authorization"].FirstOrDefault();
            var dappAuth = context.Request.Headers["X-Dapp-Authorization"].FirstOrDefault();
            var sessionId = context.Request.Headers["X-Dapp-SessionID"].FirstOrDefault();

            _logger.LogInformation("AuthMiddleware: FullPath={FullPath}, authHeader={AuthHeader}, dappAuth={DappAuth}, sessionID={SessionId}", 
                context.Request.Path, authHeader, dappAuth, sessionId);

            context.Items["Authorization"] = authHeader;
            context.Items["X-Dapp-Authorization"] = dappAuth;
            context.Items["SessionId"] = sessionId;

            await _next(context);
        }
    }

    public static class AuthMiddlewareExtensions
    {
        public static IApplicationBuilder UseAuthMiddleware(this IApplicationBuilder builder)
        {
            return builder.UseMiddleware<AuthMiddleware>();
        }
    }
} 