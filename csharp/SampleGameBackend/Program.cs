using SampleGameBackend.Database;
using SampleGameBackend.Middleware;
using SampleGameBackend.Services;

var builder = WebApplication.CreateBuilder(args);

// Add services to the container.
builder.Services.AddControllers();
builder.Services.AddEndpointsApiExplorer();
builder.Services.AddSwaggerGen();

// Add CORS
builder.Services.AddCors(options =>
{
    options.AddDefaultPolicy(policy =>
    {
        policy.AllowAnyOrigin()
              .AllowAnyMethod()
              .AllowAnyHeader();
    });
});

// Register services
builder.Services.AddSingleton<DatabaseService>();
builder.Services.AddSingleton<KeystoreService>();
builder.Services.AddScoped<ValidationService>();
builder.Services.AddScoped<ExchangeService>();

var app = builder.Build();

// Configure the HTTP request pipeline.
if (app.Environment.IsDevelopment())
{
    app.UseSwagger();
    app.UseSwaggerUI();
}

app.UseHttpsRedirection();
app.UseCors();
app.UseAuthMiddleware();
app.UseMiddleware<HmacMiddleware>();

app.MapControllers();

// Health check endpoint
app.MapGet("/health", () => Results.Ok(new { status = "healthy" }));

// 서버 시작 시 KeystoreService 초기화하여 로그 출력
var keystoreService = app.Services.GetRequiredService<KeystoreService>();
Console.WriteLine("🔑 KeystoreService initialized successfully");

Console.WriteLine("🚀 Server started on port 8080");
Console.WriteLine("📡 API endpoint: http://localhost:8080/api/assets?language=ko");
Console.WriteLine("🔐 Order validation API: http://localhost:8080/api/validate");
Console.WriteLine("💚 Health check: http://localhost:8080/health");
Console.WriteLine("💾 Session-specific asset information is stored in memory");

app.Run(); 