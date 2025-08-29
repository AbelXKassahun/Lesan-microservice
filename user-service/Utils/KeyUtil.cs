using System.Security.Cryptography;
using System.Text.Json;
using Microsoft.IdentityModel.Tokens;

namespace UserService.Utils
{
    public class KeyUtils
    {
        public void generateKeys()
        {

            static string B64Url(byte[] bytes) => Base64UrlEncoder.Encode(bytes);

            var rsa = RSA.Create(2048);

            // Export PEM files
            var privatePem = new string(PemEncoding.Write("PRIVATE KEY", rsa.ExportPkcs8PrivateKey()));
            var publicPem = new string(PemEncoding.Write("PUBLIC KEY", rsa.ExportSubjectPublicKeyInfo()));
            Directory.CreateDirectory("keys");
            File.WriteAllText(Path.Combine("keys", "private.pem"), privatePem);
            File.WriteAllText(Path.Combine("keys", "public.pem"), publicPem);

            // Build JWKS (with a stable kid)
            var p = rsa.ExportParameters(false);
            var kid = Guid.NewGuid().ToString("N"); // or compute RFC7638 thumbprint later
            var jwk = new
            {
                kty = "RSA",
                use = "sig",
                alg = "RS256",
                kid = kid,
                n = B64Url(p.Modulus!),
                e = B64Url(p.Exponent!)
            };
            var jwks = new { keys = new[] { jwk } };
            File.WriteAllText(Path.Combine("keys", "jwks.json"), JsonSerializer.Serialize(jwks, new JsonSerializerOptions { WriteIndented = true }));

            Console.WriteLine("Generated keys/private.pem, keys/public.pem, keys/jwks.json");
            Console.WriteLine($"kid: {kid}");
        }

    }
}