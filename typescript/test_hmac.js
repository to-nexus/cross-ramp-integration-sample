const crypto = require('crypto');

const salt = "my_secret_salt_value_!@#$%^&*";

const requestBody = {
  userId: 1234,
  username: "홍길동",
  email: "user@example.com",
  role: "admin",
  createdAt: 1234567890
};

const jsonString = JSON.stringify(requestBody);
console.log('Request Body:', jsonString);

const hmac = crypto.createHmac('sha256', salt);
hmac.update(jsonString);
const signature = hmac.digest('hex');

console.log('HMAC Signature:', signature); 