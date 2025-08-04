import * as crypto from 'crypto';
import { HMAC_SALT } from '../types';

export class HmacService {
  /**
   * Generate HMAC signature for given data
   * @param data - Data to sign
   * @returns HMAC signature as hex string
   */
  static generateHmac(data: any): string {
    const jsonString = JSON.stringify(data);
    const bodyBytes = Buffer.from(jsonString, 'utf8');
    
    const hmac = crypto.createHmac('sha256', HMAC_SALT);
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
    
    const hmac = crypto.createHmac('sha256', HMAC_SALT);
    hmac.update(bodyBytes);
    
    return hmac.digest('hex');
  }
} 