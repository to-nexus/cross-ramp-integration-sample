import { ethers } from 'ethers';

describe('Crypto Sign', () => {
  test('should sign digest with private key', () => {
    // Go 코드와 동일한 개인키 사용
    const privateKeyHex = '1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef';
    
    // ethers를 사용해서 개인키 생성
    const wallet = new ethers.Wallet('0x' + privateKeyHex);
    
    // 테스트 데이터
    const testData = 'test';
    
    // Keccak256 해시 생성 (Go의 crypto.Keccak256과 동일)
    const digest = ethers.keccak256(ethers.toUtf8Bytes(testData));
    
    // 해시에 서명
    const signature = wallet.signingKey.sign(digest);
    const signature2 = wallet.signMessageSync(digest);
    
    // 서명을 compact format (r + s + v)으로 변환
    const compactSignature = ethers.Signature.from(signature).serialized;
    
    

    console.log('Signature:', compactSignature);
    console.log('Signature2:', signature2);
    // 검증
    expect(compactSignature).toBeDefined();
    expect(compactSignature.length).toBe(132); // 0x + 64(r) + 64(s) + 2(v) = 132 characters
    expect(compactSignature.startsWith('0x')).toBe(true);
  });

  test('should verify signature correctly', () => {
    const privateKeyHex = '1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef';
    const wallet = new ethers.Wallet('0x' + privateKeyHex);
    
    const testData = 'test';
    const digest = ethers.keccak256(ethers.toUtf8Bytes(testData));
    
    // 서명 생성
    const signature = wallet.signingKey.sign(digest);
    
    // 서명 검증
    const recoveredAddress = ethers.recoverAddress(digest, signature);
    
    expect(recoveredAddress).toBe(wallet.address);
    console.log('Wallet address:', wallet.address);
    console.log('Recovered address:', recoveredAddress);
  });

  test('should produce consistent results', () => {
    const privateKeyHex = '1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef';
    const wallet = new ethers.Wallet('0x' + privateKeyHex);
    
    const testData = 'test';
    const digest = ethers.keccak256(ethers.toUtf8Bytes(testData));
    
    // 동일한 입력에 대해 여러 번 서명 (deterministic해야 함)
    const signature1 = wallet.signingKey.sign(digest);
    const signature2 = wallet.signingKey.sign(digest);
    
    expect(signature1.r).toBe(signature2.r);
    expect(signature1.s).toBe(signature2.s);
    expect(signature1.v).toBe(signature2.v);
  });
});
