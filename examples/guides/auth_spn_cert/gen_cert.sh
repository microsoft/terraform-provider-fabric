# Generates RSA 4096-bit private key with AES-256 encryption
openssl genpkey -algorithm RSA -pkeyopt rsa_keygen_bits:4096 -aes256 -pass pass:'YourPrivateKeyPassword' -out private.key

# Generates a self-signed certificate valid for 1 year
openssl req -subj '/CN=myclientcert/O=Contoso Inc./ST=WA/C=US' -x509 -sha256 -days 365 -passin pass:'YourPrivateKeyPassword' -key private.key -out client.crt

# Generates a PKCS12 bundle from a private key and a certificate
openssl pkcs12 -export -passin pass:'YourPrivateKeyPassword' -password pass:'YourBundlePassword' -inkey private.key -in client.crt -out bundle.pfx
