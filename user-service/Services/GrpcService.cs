using Grpc.Core;
using UserGrpc;

using UsersService.Data;

namespace UserService.service
{
    public class UserServiceImpl : UserGrpc.UserService.UserServiceBase
    {
        private readonly AppDbContext _context;

        public UserServiceImpl(AppDbContext context)
        {
            _context = context;
        }

        public override async Task<UserProfileDto> GetProfileById(UserRequest request, ServerCallContext context)
        {
            // var user = await _context.Users.FindAsync(Guid.Parse(request.UserId));
            var user = await _context.Users.FindAsync(request.UserId);
            if (user == null)
            {
                throw new RpcException(new Status(StatusCode.NotFound, "User not found"));
            }

            return new UserProfileDto
            {
                Id = user.Id.ToString(),
                Email = user.Email,
            };
        }
    }
}