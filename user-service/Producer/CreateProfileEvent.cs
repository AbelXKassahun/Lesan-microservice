using System.Text.Json.Serialization;

namespace ProfileProducer;

/// <summary>
/// Represents the message payload expected by the Go consumer.
/// </summary>
public class CreateProfileEvent
{
    [JsonPropertyName("event")]
    public string Event { get; set; } = "create_profile";

    [JsonPropertyName("user_id")]
    public string UserId { get; set; }

    [JsonPropertyName("email")]
    public string Email { get; set; }
}