const crypto = require('crypto');

const salt = "my_secret_salt_value_!@#$%^&*";

const requestBody = {
  user_sig: "0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
  user_address: "0xB777C937fa1afC99606aFa85c5b83cFe7f82BabD",
  project_id: "acjviwejsi",
  digest: "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
  uuid: "test-uuid-123",
  intent: {
    method: "mint",
    from: [
      { type: "asset", id: "asset_money", amount: 1000 }
    ],
    to: [
      { type: "asset", id: "item_gem", amount: 1000 }
    ]
  }
};

const jsonString = JSON.stringify(requestBody);
console.log('Request Body:', jsonString);

const hmac = crypto.createHmac('sha256', salt);
hmac.update(jsonString);
const signature = hmac.digest('hex');

console.log('HMAC Signature:', signature); 