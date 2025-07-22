// using Microsoft.AspNetCore.StaticFiles;

namespace UserService.Utils
{
    public class ProfileUtils
    {
        public string GetContentType(string path)
        {
            var provider = new Microsoft.AspNetCore.StaticFiles.FileExtensionContentTypeProvider();
            if (!provider.TryGetContentType(path, out var contentType))
            {
                contentType = "application/octet-stream"; // Fallback
            }
            return contentType;
        }
    }
}