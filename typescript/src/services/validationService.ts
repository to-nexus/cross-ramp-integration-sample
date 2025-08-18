import { ethers } from 'ethers';
import { ValidateRequestIntent, Asset } from '../types';
import { MemoryDatabase } from '../database/memoryDb';

export class ValidationService {
  private db: MemoryDatabase;

  constructor(db: MemoryDatabase) {
    this.db = db;
  }

  /**
   * Validate intent structure
   * @param intent - Intent to validate
   * @returns true if valid
   */
  validateIntent(intent: ValidateRequestIntent): boolean {
    if (!intent || !intent.method) {
      return false;
    }

    const validMethods = ['mint', 'transfer', 'burn'];
    if (!validMethods.includes(intent.method)) {
      return false;
    }

    // Validate from assets
    if (!Array.isArray(intent.from) || intent.from.length === 0) {
      return false;
    }

    for (const fromAsset of intent.from) {
      if (!fromAsset.type || !fromAsset.id || typeof fromAsset.amount !== 'number') {
        return false;
      }
    }

    // Validate to assets
    if (!Array.isArray(intent.to) || intent.to.length === 0) {
      return false;
    }

    for (const toAsset of intent.to) {
      if (!toAsset.type || !toAsset.id || typeof toAsset.amount !== 'number') {
        return false;
      }
    }

    return true;
  }

  /**
   * Validate and process mint operation
   * @param sessionId - Session ID
   * @param fromAssets - Assets to deduct
   * @throws Error if insufficient balance
   */
  validateAndProcessMint(sessionId: string, fromAssets: Array<{ type: string; id: string; amount: number }>): void {
    const sessionAssets = this.db.getSessionAssets(sessionId);
    if (!sessionAssets) {
      throw new Error('Session not found');
    }

    // Check if user has sufficient balance for each asset
    for (const fromAsset of fromAssets) {
      const userAsset = sessionAssets.assets.find(asset => asset.id === fromAsset.id);
      if (!userAsset) {
        throw new Error(`Asset ${fromAsset.id} not found`);
      }

      const userBalance = parseInt(userAsset.balance);
      if (userBalance < fromAsset.amount) {
        throw new Error(`Insufficient balance for asset ${fromAsset.id}`);
      }
    }

    // Deduct assets (in real implementation, this would be a transaction)
    const updatedAssets = sessionAssets.assets.map(asset => {
      const fromAsset = fromAssets.find(fa => fa.id === asset.id);
      if (fromAsset) {
        const currentBalance = parseInt(asset.balance);
        const newBalance = currentBalance - fromAsset.amount;
        return { ...asset, balance: newBalance.toString() };
      }
      return asset;
    });

    this.db.updateSessionAssets(sessionId, updatedAssets);
  }

  /**
   * Generate validator signature
   * @param userSig - User signature
   * @param digest - Digest hash (already keccak256 hashed)
   * @returns Validator signature
   */
  async generateValidatorSignature(userSig: string, digest: string): Promise<string> {
    try {
      // TODO: In actual implementation, use validator's private key
      // For now, we'll create a mock signature
      const mockPrivateKey = "0x1234567890123456789012345678901234567890123456789012345678901234";
      const wallet = new ethers.Wallet(mockPrivateKey);
      
      // digest is already keccak256 hashed, so we can sign it directly
      const digestBytes = ethers.getBytes(digest);
      
      const signature = await wallet.signingKey.sign(digestBytes);
      const compactSignature = ethers.Signature.from(signature).serialized;
      return compactSignature;
    } catch (error) {
      throw new Error('Failed to generate validator signature');
    }
  }

  /**
   * Process exchange result
   * @param sessionId - Session ID
   * @param outputs - Output assets
   * @param receiptStatus - Receipt status
   */
  processExchangeResult(sessionId: string, outputs: Array<{ asset_id: string; amount: number }>, receiptStatus: number): void {
    if (receiptStatus !== 1) {
      throw new Error('Transaction failed');
    }

    const sessionAssets = this.db.getSessionAssets(sessionId);
    if (!sessionAssets) {
      throw new Error('Session not found');
    }

    // Add output assets to user's inventory
    const updatedAssets = [...sessionAssets.assets];
    
    for (const output of outputs) {
      const existingAsset = updatedAssets.find(asset => asset.id === output.asset_id);
      if (existingAsset) {
        const currentBalance = parseInt(existingAsset.balance);
        const newBalance = currentBalance + output.amount;
        existingAsset.balance = newBalance.toString();
      } else {
        updatedAssets.push({
          id: output.asset_id,
          balance: output.amount.toString()
        });
      }
    }

    this.db.updateSessionAssets(sessionId, updatedAssets);
  }
} 