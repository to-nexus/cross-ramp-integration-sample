import * as crypto from 'crypto';
import { HMAC_SALT } from '../types';

export class HmacService {
  /**
   * Base64 URL decode function
   * @param str - Base64 URL encoded string
   * @returns Decoded buffer
   */
  private static base64UrlDecode(str: string): Buffer {
    // Convert URL safe base64 to standard base64
    str = str.replace(/-/g, '+').replace(/_/g, '/');
    // Add padding (if needed)
    while (str.length % 4) {
      str += '=';
    }
    return Buffer.from(str, 'base64');
  }

  /**
   * Generate HMAC signature for given data
   * @param data - Data to sign
   * @returns HMAC signature as hex string
   */
  static generateHmac(data: any): string {
    const jsonString = JSON.stringify(data);
    const bodyBytes = Buffer.from(jsonString, 'utf8');
    
    // Use Base64 URL decoding as per guide
    const saltBytes = this.base64UrlDecode(HMAC_SALT);
    const hmac = crypto.createHmac('sha256', saltBytes);
    hmac.update(bodyBytes);
    
    return hmac.digest('hex');
  }

  /**
   * Validate HMAC signature
   * @param requestBody - Request body as string
   * @param hmacSignature - HMAC signature from header
   * @returns true if signature is valid
   */
  static validateHmac(requestBody: string, hmacSignature: string): boolean {
    if (!hmacSignature) {
      return false;
    }

    const calculatedHmac = this.generateHmacFromString(requestBody);
    return calculatedHmac.toLowerCase() === hmacSignature.toLowerCase();
  }

  /**
   * Generate HMAC signature for request body
   * @param requestBody - Request body as string
   * @returns HMAC signature as hex string
   */
  static generateHmacFromString(requestBody: string): string {
    const bodyBytes = Buffer.from(requestBody, 'utf8');
    
    // Use Base64 URL decoding as per guide
    const saltBytes = this.base64UrlDecode(HMAC_SALT);
    const hmac = crypto.createHmac('sha256', saltBytes);
    hmac.update(bodyBytes);
    
    return hmac.digest('hex');
  }
} 