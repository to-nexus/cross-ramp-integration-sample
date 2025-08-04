using Microsoft.AspNetCore.Http;
using SampleGameBackend.Services;
using System.Text;

namespace SampleGameBackend.Middleware
{
    public class HmacMiddleware
    {
        private readonly RequestDelegate _next;

        public HmacMiddleware(RequestDelegate next)
        {
            _next = next;
        }

        public async Task InvokeAsync(HttpContext context)
        {
            // Skip HMAC validation for GET requests
            if (context.Request.Method == "GET")
            {
                await _next(context);
                return;
            }

            // Read request body
            context.Request.EnableBuffering();
            var body = await new StreamReader(context.Request.Body).ReadToEndAsync();
            context.Request.Body.Position = 0;

            // Get HMAC signature from header
            var hmacSignature = context.Request.Headers["X-HMAC-Signature"].FirstOrDefault();

            // Validate HMAC
            if (!HmacService.ValidateHmac(body, hmacSignature))
            {
                context.Response.StatusCode = 401;
                await context.Response.WriteAsJsonAsync(new
                {
                    success = false,
                    errorCode = "INVALID_HMAC_SIGNATURE",
                    message = "Invalid HMAC signature"
                });
                return;
            }

            await _next(context);
        }
    }
} 