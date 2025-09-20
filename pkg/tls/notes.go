package tls

// Only TLS 1.3 (must enforce with MinVersion: tls.VersionTLS13)

// Certificates
// - PROD: use trusted system CAs
// - DEV: generate self-signed certificate with SAN:
//   openssl req -x509 -newkey rsa:4096 -sha256 -days 365 -nodes \
//     -keyout server.key -out server.crt \
//     -subj "/CN=localhost" \
//     -addext "subjectAltName=DNS:localhost,IP:127.0.0.1"
