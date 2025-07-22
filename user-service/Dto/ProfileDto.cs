namespace Profile.Dto
{
    public class UserProfileDto
    {
        public required string Id { get; set; }
        public required string Email { get; set; }
        public string? FullName { get; set; }
        public string? UserName { get; set; }
        public string? ProfileImage { get; set; }
        public string? ProfileImageType { get; set; }
    }

    public class UpdateProfileDto
    {
        public string? FullName { get; set; }
        public string? UserName { get; set; } // updated once
        // public string? Email { get; set; }  // doesnt make sense to update this
    }

}