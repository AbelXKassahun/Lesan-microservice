using Microsoft.AspNetCore.Mvc;
using Microsoft.AspNetCore.Identity;

using UserService.Models;
using Auth.Dto;
using System.Security.Claims;
using UserService.Utils;

namespace Auth.Controller
{
    [ApiController]
    // [Route("/api/v{version:apiVersion}/user/[controller]")]
    [Route("/api/user/[controller]")]
    // [ApiVersion("2.0")]
    public class AuthController : ControllerBase
    {
        private readonly UserManager<ApplicationUser> _userManager;
        private readonly SignInManager<ApplicationUser> _signInManager;

        public AuthController(SignInManager<ApplicationUser> signInManager, UserManager<ApplicationUser> userManager)
        {
            _signInManager = signInManager;
            _userManager = userManager;
        }

        [HttpPost("sign-up")]
        // [MapToApiVersion("2.0")]
        public async Task<IActionResult> Signup([FromBody] AuthCredentialsDto model)
        {
            var user = new ApplicationUser
            {
                Email = model.Email,
                UserName = model.Email, // this is becuase UserName is required but i wouldnt prompt the user for one during signup
            };


            var result = await _userManager.CreateAsync(user, model.Password);

            if (!result.Succeeded)
            {
                var errors = result.Errors.Select(e => new { e.Code, e.Description });
                return BadRequest(new { errors });
            }

            // create acess token
            AuthUtils util = new AuthUtils();
            string jwt = util.GenerateAccessToken(new ClaimsDto
            {
                Id = user.Id,
                Email = model.Email,
                Role = "User",
            });

            // create refresh token and store it 
            var refreshToken = util.GenerateRefreshToken();
            user.RefreshToken = refreshToken;
            user.RefreshTokenExpiryTime = DateTime.UtcNow.AddDays(7);
            result = await _userManager.UpdateAsync(user);
            if (!result.Succeeded)
            {
                var errors = result.Errors.Select(e => new { e.Code, e.Description });
                return BadRequest(new { errors });
            }

            return Ok(new
            {
                Message = "User created successfully",
                UserId = user.Id,
                AccessToken = jwt,
                RefreshToken = refreshToken,
            });
        }

        [HttpPost("sign-in")]
        // [MapToApiVersion("2.0")]
        public async Task<IActionResult> SignIn([FromBody] AuthCredentialsDto model)
        {
            var user = await _userManager.FindByEmailAsync(model.Email);
            if (user == null)
                return Unauthorized("Invalid email or password.");

            // validate credentials
            var result = await _signInManager.CheckPasswordSignInAsync(user, model.Password, lockoutOnFailure: false);

            if (!result.Succeeded)
                return Unauthorized("Invalid email or password.");

            // create token
            AuthUtils util = new AuthUtils();
            string jwt = util.GenerateAccessToken(new ClaimsDto
            {
                Id = user.Id,
                Email = model.Email,
                Role = "User",
            });
            Console.WriteLine($"@@@ user-id {user.Id}");

            // create refresh token and store it 
            var refreshToken = util.GenerateRefreshToken();
            user.RefreshToken = refreshToken;
            user.RefreshTokenExpiryTime = DateTime.UtcNow.AddDays(7);
            var reslt = await _userManager.UpdateAsync(user);
            if (!reslt.Succeeded)
            {
                var errors = reslt.Errors.Select(e => new { e.Code, e.Description });
                return BadRequest(new { errors });
            }

            return Ok(new
            {
                Message = "Login successful",
                userId = user.Id,
                AccessToken = jwt,
                RefreshToken = refreshToken,
            });
        }


        [HttpPost("refresh-token")]
        // [MapToApiVersion("2.0")]
        public async Task<IActionResult> Refresh([FromBody] TokenRefreshRequest model)
        {
            AuthUtils util = new AuthUtils();

            if (string.IsNullOrWhiteSpace(model.AccessToken) || string.IsNullOrWhiteSpace(model.RefreshToken))
                return BadRequest("Access token or refresh token is missing.");

            var principal = util.GetPrincipalFromExpiredToken(model.AccessToken);
            if (principal == null)
                return BadRequest("Invalid access token");

            var userId = principal.FindFirstValue(ClaimTypes.NameIdentifier);
            var user = await _userManager.FindByIdAsync(userId);
            if (user == null || user.RefreshToken != model.RefreshToken || user.RefreshTokenExpiryTime <= DateTime.UtcNow)
                return BadRequest("Invalid refresh token");

            var newAccessToken = util.GenerateAccessToken(new ClaimsDto
            {
                Id = user.Id,
                Email = user.Email,
                Role = "User",
            });

            var newRefreshToken = util.GenerateRefreshToken();
            user.RefreshToken = newRefreshToken;
            user.RefreshTokenExpiryTime = DateTime.UtcNow.AddDays(7);
            await _userManager.UpdateAsync(user);

            return Ok(new
            {
                AccessToken = newAccessToken,
                RefreshToken = newRefreshToken
            });
        }

    }
}