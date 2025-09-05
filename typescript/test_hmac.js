const crypto = require('crypto');

// Base64 URL decoding function (as per guide specification)
function base64UrlDecode(str) {
  // Convert URL safe base64 to standard base64
  str = str.replace(/-/g, '+').replace(/_/g, '/');
  // Add padding (if needed)
  while (str.length % 4) {
    str += '=';
  }
  return Buffer.from(str, 'base64');
}

function hmacSha256(data, salt) {
  return crypto.createHmac('sha256', base64UrlDecode(salt)).update(data).digest('hex');
}

// Define JSON object request body
const requestBody = {
  userId: 1234,
  username: '홍길동',
  email: 'user@example.com',
  role: 'admin',
  createdAt: 1234567890,
};

const salt = 'my_secret_salt_value_!@#$%^&*'; // hmac key
const jsonString = JSON.stringify(requestBody);

console.log('JSON string:', jsonString);
// Output result
console.log('HMAC-SHA256:', hmacSha256(jsonString, salt)); 
// expected X-HMAC-Signature: f96cf60394f6b8ad3c6de2d5b2b1d1a540f9529082a8eb9cee405bfbdd9f37a1 