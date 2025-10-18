using System.Text;
using System.Text.Json;
using RabbitMQ.Client; // This should work now

namespace ProfileProducer{
    public class ProfileEventPublisher : IAsyncDisposable
    {
        private readonly IConnection _connection;
        private readonly IModel _channel;
        private const string QueueName = "create_profile";

        // Constructor is now private
        private ProfileEventPublisher(IConnection connection, IModel channel)
        {
            _connection = connection;
            _channel = channel;
        }

        // Public static factory method for async initialization
        public static async Task<ProfileEventPublisher> CreateAsync(string rabbitURL)
        {
            var factory = new ConnectionFactory() { Uri = new Uri(rabbitURL) };

            // 1. CreateConnection is now CreateConnectionAsync
            var connection = await factory.CreateConnectionAsync();

            // 2. CreateModel is now CreateModelAsync
            var channel = connection.CreateModel();

            // Idempotently declare the queue
            await channel.QueueDeclareAsync(queue: QueueName,
                                durable: true,
                                exclusive: false,
                                autoDelete: false,
                                arguments: null);

            return new ProfileEventPublisher(connection, channel);
        }

        public void PublishCreateProfileEvent(string userId, string email)
        {
            var messagePayload = new CreateProfileEvent
            {
                UserId = userId,
                Email = email
            };

            string jsonMessage = JsonSerializer.Serialize(messagePayload);
            var body = Encoding.UTF8.GetBytes(jsonMessage);

            var properties = _channel.CreateBasicProperties();
            properties.Persistent = true;
            properties.ContentType = "application/json";

            Console.WriteLine($"Publishing event for UserID: {userId}...");

            _channel.BasicPublish(exchange: "",
                                routingKey: QueueName,
                                basicProperties: properties,
                                body: body);

            Console.WriteLine($"✅ Message sent successfully.");
        }

        // Implement IAsyncDisposable
        public async ValueTask DisposeAsync()
        {
            // 3. Just dispose, no .Close() needed
            if (_channel != null)
            {
                await _channel.Dispose();
            }
            if (_connection != null)
            {
                await _connection.DisposeAsync();
            }
        }
    }
}
