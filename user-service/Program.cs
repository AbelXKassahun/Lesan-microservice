using Microsoft.AspNetCore.Mvc;
// using Microsoft.AspNetCore.Mvc.Versioning;
using Microsoft.EntityFrameworkCore;
using Microsoft.AspNetCore.Identity;
using Microsoft.AspNetCore.Authentication.JwtBearer;
using Microsoft.IdentityModel.Tokens;
using System.Text;
using Microsoft.AspNetCore.Server.Kestrel.Core;

using TicketGrpc;

using UsersService.Data;
using UserService.Models;
using UserService.Utils;
using UserService.service;
using StackExchange.Redis;



var builder = WebApplication.CreateBuilder(args);

// this is run only once
// KeyUtils keyUtils = new KeyUtils();
// keyUtils.generateKeys();

builder.Services.AddControllers();

// Add services to the container.
// Learn more about configuring Swagger/OpenAPI at https://aka.ms/aspnetcore/swashbuckle
builder.Services.AddEndpointsApiExplorer();
builder.Services.AddSwaggerGen();
// Load .env file
DotNetEnv.Env.Load();
var redisConnectionString = Environment.GetEnvironmentVariable("RedisConnectionString");

var host = Environment.GetEnvironmentVariable("POSTGRES_HOST");
var port = Environment.GetEnvironmentVariable("POSTGRES_PORT");
var db = Environment.GetEnvironmentVariable("POSTGRES_DB");
var user = Environment.GetEnvironmentVariable("POSTGRES_USER");
var password = Environment.GetEnvironmentVariable("POSTGRES_PASSWORD");

var connectionString = $"Host={host};Port={port};Database={db};Username={user};Password={password}";

builder.Services.AddGrpc();

builder.Services.AddGrpcClient<TicketCleanupService.TicketCleanupServiceClient>(o =>
{
    o.Address = new Uri("http://ticket-service:5144");
});

builder.WebHost.ConfigureKestrel(options =>
{
    // For gRPC: HTTP/2 only
    options.ListenAnyIP(5168, listenOptions =>
    {
        listenOptions.Protocols = HttpProtocols.Http2;
    });

    // For REST: HTTP/1.1
    options.ListenAnyIP(5080, listenOptions =>
    {
        listenOptions.Protocols = HttpProtocols.Http1;
    });
});

// Register DbContext
builder.Services.AddDbContext<AppDbContext>(options =>
    options.UseNpgsql(connectionString));

// Redis
builder.Services.AddSingleton<IConnectionMultiplexer>(
    _ => ConnectionMultiplexer.Connect(redisConnectionString)); 

builder.Services.AddScoped<IRedisCacheService, RedisCacheService>();

builder.Services.AddApiVersioning(options =>
{
    options.AssumeDefaultVersionWhenUnspecified = true;
    options.DefaultApiVersion = new ApiVersion(2, 0);
    options.ReportApiVersions = true;
});

builder.Services.AddIdentity<ApplicationUser, IdentityRole>(options =>
{
    options.User.RequireUniqueEmail = true;
    options.SignIn.RequireConfirmedEmail = false;
})
.AddEntityFrameworkStores<AppDbContext>()
.AddDefaultTokenProviders();

DotNetEnv.Env.Load("../");
string secretKey = Environment.GetEnvironmentVariable("JWT_SECRET_KEY");
var key = Encoding.UTF8.GetBytes(secretKey);

builder.Services.AddAuthentication(options =>
{
    options.DefaultAuthenticateScheme = JwtBearerDefaults.AuthenticationScheme;
    options.DefaultChallengeScheme = JwtBearerDefaults.AuthenticationScheme;
})
.AddJwtBearer(options =>
{
    options.TokenValidationParameters = new TokenValidationParameters
    {
        ValidateIssuer = true,
        ValidateAudience = true,
        ValidateLifetime = true,
        ValidateIssuerSigningKey = true,
        ValidIssuer = "users-service",
        ValidAudience = "ticket-users",
        IssuerSigningKey = new SymmetricSecurityKey(key)
    };
});
builder.Services.AddAuthorization();

// password rules
builder.Services.Configure<IdentityOptions>(options =>
{
    options.Password.RequiredLength = 12;
    options.Password.RequireDigit = true;
    options.Password.RequireUppercase = true;
    options.Password.RequireNonAlphanumeric = false;
});

// DI
builder.Services.AddSingleton<ProfileUtils>();

var app = builder.Build();

app.MapGrpcService<UserServiceImpl>();
// app.UseRouting();
app.MapControllers();

// foreach (var endpoint in app.Services.GetRequiredService<EndpointDataSource>().Endpoints)
// {
//     Console.WriteLine(endpoint.DisplayName);
// }
// Run EF migrations automatically on startup
// using (var scope = app.Services.CreateScope())
// {
//     var dbContext = scope.ServiceProvider.GetRequiredService<AppDbContext>();
//     dbContext.Database.Migrate();
// }

// Configure the HTTP request pipeline.
if (app.Environment.IsDevelopment())
{
    app.UseSwagger();
    app.UseSwaggerUI();
}

app.UseHttpsRedirection();

app.UseAuthentication();
app.UseAuthorization();

app.Run();