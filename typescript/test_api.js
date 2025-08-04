const crypto = require('crypto');

// HMAC salt
const salt = "my_secret_salt_value_!@#$%^&*";

// Test request body
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

// Generate HMAC signature
function generateHmac(data) {
  const jsonString = JSON.stringify(data);
  const hmac = crypto.createHmac('sha256', salt);
  hmac.update(jsonString);
  return hmac.digest('hex');
}

const hmacSignature = generateHmac(requestBody);
console.log('HMAC Signature:', hmacSignature);

// Test validation API
import fetch from 'node-fetch';

async function testValidationAPI() {
  try {
    const response = await fetch('http://localhost:8080/api/validate', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': 'Bearer test_cross_auth_jwt_token',
        'X-Dapp-Authorization': 'Bearer test_dapp_access_token',
        'X-Dapp-SessionID': 'test_session_id',
        'X-HMAC-Signature': hmacSignature
      },
      body: JSON.stringify(requestBody)
    });

    const result = await response.json();
    console.log('Validation API Response:', JSON.stringify(result, null, 2));
  } catch (error) {
    console.error('Error:', error);
  }
}

testValidationAPI(); 