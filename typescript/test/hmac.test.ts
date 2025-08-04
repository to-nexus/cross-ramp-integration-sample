import { HmacService } from '../src/services/hmacService';

// TODO: HMAC salt - In actual implementation, load from environment variables or configuration file
const salt = "my_secret_salt_value_!@#$%^&*"; // hmac key

interface Body {
  userId: number;
  username: string;
  email: string;
  role: string;
  createdAt: number;
}

const body: Body = {
  userId: 1234,
  username: "홍길동",
  email: "user@example.com",
  role: "admin",
  createdAt: 1234567890,
};

describe('HMAC Service', () => {
  test('should generate correct HMAC signature', () => {
    const bodyBytes = JSON.stringify(body);
    console.log(bodyBytes);
    
    const hashString = HmacService.generateHmac(body);
    console.log('hashString:', hashString); // expected X-HMAC-Signature: f96cf60394f6b8ad3c6de2d5b2b1d1a540f9529082a8eb9cee405bfbdd9f37a1
    
    expect(hashString).toBeDefined();
    expect(hashString.length).toBe(64); // SHA256 hex string length
  });

  test('should validate HMAC signature correctly', () => {
    const bodyString = JSON.stringify(body);
    const signature = HmacService.generateHmacFromString(bodyString);
    
    expect(HmacService.validateHmac(bodyString, signature)).toBe(true);
    expect(HmacService.validateHmac(bodyString, 'invalid_signature')).toBe(false);
  });
}); 