using Microsoft.AspNetCore.Mvc;
using Microsoft.AspNetCore.Identity;
using System.Security.Claims;
using Microsoft.AspNetCore.Authorization;
using Microsoft.AspNetCore.Authentication.JwtBearer;

using TicketGrpc;

using UserService.Models;
using UserService.Utils;
using Profile.Dto;


namespace Profile.Controller
{
    [ApiController]
    // [Route("/api/v{version:apiVersion}/user/[controller]")]
    [Route("/api/user/[controller]")]
    // [ApiVersion("2.0")]
    [Authorize(AuthenticationSchemes = JwtBearerDefaults.AuthenticationScheme)]
    public class ProfileController : ControllerBase
    {
        private readonly UserManager<ApplicationUser> _userManager;
        private readonly TicketCleanupService.TicketCleanupServiceClient _ticketClient;
        private readonly ProfileUtils _profileUtils;

        public ProfileController(UserManager<ApplicationUser> userManager, TicketCleanupService.TicketCleanupServiceClient ticketClient, ProfileUtils profileUtils)
        {
            _userManager = userManager;
            _ticketClient = ticketClient;
            _profileUtils = profileUtils;
        }

        [HttpGet("me")]
        public async Task<IActionResult> GetOwnProfile()
        {
            bool img404 = false;

            var userId = User.FindFirstValue(ClaimTypes.NameIdentifier);
            var user = await _userManager.FindByIdAsync(userId);
            if (user == null) return NotFound();

            var userImage = user.ProfileImagePath;
            if (userImage == null || !System.IO.File.Exists(userImage))
                img404 = true;

            var fileBytes = await System.IO.File.ReadAllBytesAsync(userImage);
            var contentType = _profileUtils.GetContentType(userImage);

            return Ok(new UserProfileDto
            {
                Id = user.Id,
                Email = user.Email,
                UserName = user.UserName == user.Email ? "" : user.UserName,
                FullName = string.IsNullOrWhiteSpace(user.FullName) ? "" : user.FullName,
                ProfileImage = img404 ? Convert.ToBase64String(fileBytes) : "Image not found", // base64 image 
                ProfileImageType = img404 ? contentType : "-",
            });
        }

        [HttpGet("{id}")]
        public async Task<IActionResult> GetProfileById(string id)
        {
            bool img404 = false;
            var user = await _userManager.FindByIdAsync(id);
            if (user == null) return NotFound();

            var userImage = user.ProfileImagePath;
            if (userImage == null || !System.IO.File.Exists(userImage))
                img404 = true;

            var fileBytes = await System.IO.File.ReadAllBytesAsync(userImage);
            var contentType = _profileUtils.GetContentType(userImage);

            return Ok(new UserProfileDto
            {
                Id = user.Id,
                Email = user.Email,
                UserName = user.UserName == user.Email ? "" : user.UserName,
                FullName = string.IsNullOrWhiteSpace(user.FullName) ? "" : user.FullName,
                ProfileImage = img404 ? Convert.ToBase64String(fileBytes) : "Image not found", // base64 image 
                ProfileImageType = img404 ? contentType : "-",
            });
        }

        [HttpPut("update")]
        public async Task<IActionResult> UpdateProfileInfo([FromBody] UpdateProfileDto dto) // previously was FromBody
        {
            var userId = User.FindFirstValue(ClaimTypes.NameIdentifier);
            var user = await _userManager.FindByIdAsync(userId);
            if (user == null) return NotFound();

            bool updated = false;

            if (!string.IsNullOrWhiteSpace(dto.FullName) && dto.FullName != user.FullName)
            {
                user.FullName = dto.FullName;
                updated = true;
            }

            if (!string.IsNullOrWhiteSpace(dto.UserName) && user.UserName == user.Email)
            {
                user.UserName = dto.UserName;
                updated = true;
            }

            if (!updated)
                return BadRequest("No changes provided or values are empty.");

            var result = await _userManager.UpdateAsync(user);
            if (!result.Succeeded)
                return BadRequest(result.Errors);

            return Ok(new UserProfileDto
            {
                Id = user.Id,
                Email = user.Email,
                UserName = user.UserName == user.Email ? "" : user.UserName,
                FullName = string.IsNullOrWhiteSpace(user.FullName) ? "" : user.FullName,
            });
        }

        [HttpPost]
        public async Task<IActionResult> UpdateProfileImage([FromForm] IFormFile image)
        {
            var userId = User.FindFirstValue(ClaimTypes.NameIdentifier);
            var user = await _userManager.FindByIdAsync(userId);
            if (user == null) return NotFound();

            if (image == null || image.Length == 0)
                return BadRequest("No file uploaded.");

            var parentDirectory = Directory.GetParent(Directory.GetCurrentDirectory()).FullName;
            var uploadPath = Path.Combine(parentDirectory, "uploads");
            // if (!Directory.Exists(uploadPath)) // no need to check for this 
            //     Directory.CreateDirectory(uploadPath);

            var fileName = $"{Guid.NewGuid()}_{image.FileName}";
            var filePath = Path.Combine(uploadPath, fileName);

            if (user.ProfileImagePath == filePath)
            {
                return BadRequest("No changes were provided");
            }

            using (var stream = new FileStream(filePath, FileMode.Create))
            {
                await image.CopyToAsync(stream);
            }

            // Save to Database
            user.ProfileImagePath = filePath;
            var result = await _userManager.UpdateAsync(user);
            if (!result.Succeeded)
                return BadRequest(result.Errors);
            return Ok("File uploaded successfully"); // return Ok(new { message = "File uploaded successfully" });
        }

        [HttpDelete("delete")]
        public async Task<IActionResult> DeleteAccount()
        {
            var userId = User.FindFirstValue(ClaimTypes.NameIdentifier);
            var user = await _userManager.FindByIdAsync(userId);
            if (user == null) return NotFound();

            // cleanup message through the broker for all services/subscribers

            var result = await _userManager.DeleteAsync(user);
            if (!result.Succeeded)
                return BadRequest(result.Errors);

            DeleteUserImage(user.ProfileImagePath);

            return Ok(new
            {
                Message = "Account deleted successfully",
            });
        }

        public void DeleteUserImage(string? userImage)
        {
            try
            {
                if (System.IO.File.Exists(userImage))
                {
                    System.IO.File.Delete(userImage);
                }

                Console.WriteLine($"Image deleted, {userImage}");
            }
            catch (Exception ex)
            {
                Console.WriteLine($"Error deleting image: {ex.Message}");
            }
        }
    }
}