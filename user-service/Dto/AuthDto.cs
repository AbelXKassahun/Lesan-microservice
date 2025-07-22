using System.ComponentModel.DataAnnotations;

namespace Auth.Dto
{
    public class AuthCredentialsDto
    {
        [Required]
        [EmailAddress]
        public required string Email { get; set; }

        [Required]
        public required string Password { get; set; }
    }

    public class ClaimsDto
    {
        public required string Id { get; set; }
        public required string Email { get; set; }
        public required string Role { get; set; } = "User";
    }

    public class TokenRefreshRequest
    {
        public required string AccessToken { get; set; }
        public required string RefreshToken { get; set; }
    }
}