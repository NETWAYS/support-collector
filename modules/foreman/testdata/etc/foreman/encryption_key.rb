# Be sure to restart your server when you modify this file.

# Your encryption key for encrypting and decrypting database fields.
# If you change this key, all encrypted data will NOT be able to be decrypted by Foreman!
# Make sure the key is at least 32 bytes such as SecureRandom.hex(20)

# You can use `rake security:generate_encryption_key` to regenerate this file.

module EncryptionKey
  ENCRYPTION_KEY = '7ee82cb4cb031bc3415be013d8bfc28c'
end