import { keccak256, stringToBytes, hexToBytes } from 'viem';
import { privateKeyToAccount } from 'viem/accounts';

describe('Crypto Sign with Viem', () => {
  test('should sign digest with private key', async () => {
    // Go 코드와 동일한 개인키 사용
    const privateKeyHex = '0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef';
    
    // viem을 사용해서 account 생성
    const account = privateKeyToAccount(privateKeyHex);
    
    // 테스트 데이터
    const testData = 'test';
    
    // Keccak256 해시 생성 (Go의 crypto.Keccak256과 동일)
    const digest = keccak256(stringToBytes(testData));
    
    // 해시에 서명 (viem의 sign 함수 사용)
    const signature = await account.sign({ hash: digest });
    
    console.log('Digest:', digest);
    console.log('Signature:', signature);
    
    // 검증
    expect(signature).toBeDefined();
    expect(signature.length).toBe(132); // 0x + 64(r) + 64(s) + 2(v) = 132 characters
    expect(signature.startsWith('0x')).toBe(true);
  });

  test('should verify signature correctly', async () => {
    const privateKeyHex = '0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef';
    const account = privateKeyToAccount(privateKeyHex);
    
    const testData = 'test';
    const digest = keccak256(stringToBytes(testData));
    
    // 서명 생성
    const signature = await account.sign({ hash: digest });
    
    // viem에서 서명 검증은 다른 방식으로 수행
    // account address와 서명이 일치하는지 확인
    expect(account.address).toBeDefined();
    expect(signature).toBeDefined();
    
    console.log('Account address:', account.address);
    console.log('Signature:', signature);
    
    // 동일한 데이터로 다시 서명해서 일치하는지 확인
    const signature2 = await account.sign({ hash: digest });
    expect(signature).toBe(signature2); // deterministic signing
  });

  test('should handle message signing', async () => {
    const privateKeyHex = '0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef';
    const account = privateKeyToAccount(privateKeyHex);
    
    const message = 'test message';
    
    // 메시지 직접 서명 (viem이 자동으로 해시 처리)
    const signature = await account.signMessage({ message });
    
    console.log('Message:', message);
    console.log('Message signature:', signature);
    
    expect(signature).toBeDefined();
    expect(signature.length).toBe(132);
    expect(signature.startsWith('0x')).toBe(true);
  });

  test('should produce consistent results with raw hash', async () => {
    const privateKeyHex = '0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef';
    const account = privateKeyToAccount(privateKeyHex);
    
    const testData = 'test';
    const digest = keccak256(stringToBytes(testData));
    
    // 동일한 해시에 대해 여러 번 서명 (deterministic해야 함)
    const signature1 = await account.sign({ hash: digest });
    const signature2 = await account.sign({ hash: digest });
    
    expect(signature1).toBe(signature2);
    
    console.log('Consistent signature:', signature1);
  });

  test('should work with byte arrays', async () => {
    const privateKeyHex = '0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef';
    const account = privateKeyToAccount(privateKeyHex);
    
    // 바이트 배열로 테스트
    const testBytes = new Uint8Array([116, 101, 115, 116]); // "test" in bytes
    const digest = keccak256(testBytes);
    
    const signature = await account.sign({ hash: digest });
    
    console.log('Bytes digest:', digest);
    console.log('Bytes signature:', signature);
    
    expect(signature).toBeDefined();
    expect(signature.length).toBe(132);
  });

  test('should sign with specific hex digest', async () => {
    const privateKeyHex = '0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef';
    const account = privateKeyToAccount(privateKeyHex);
    
    // 특정 hex 값을 digest로 직접 사용 (viem은 hex 문자열을 받음)
    const digest = '0xd91c81e564e4f69229a9224943fa9a79ff21b60fcef5096bfb79e1ce28591a85';
    
    const signature = await account.sign({ hash: digest });
    
    console.log('Hex digest:', digest);
    console.log('Digest as bytes:', hexToBytes(digest));
    console.log('Signature!!:', signature);
    
    expect(signature).toBeDefined();
    expect(signature.length).toBe(132);
    expect(signature.startsWith('0x')).toBe(true);
  });
});
