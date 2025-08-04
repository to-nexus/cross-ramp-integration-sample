import { Response } from 'express';
import { AuthenticatedRequest } from '../middleware/authMiddleware';
import { ValidateRequest, ValidateResponse } from '../types';
import { ValidationService } from '../services/validationService';
import { MemoryDatabase } from '../database/memoryDb';

export class ValidationController {
  private validationService: ValidationService;
  private db: MemoryDatabase;

  constructor(validationService: ValidationService, db: MemoryDatabase) {
    this.validationService = validationService;
    this.db = db;
  }

  /**
   * Validate user action
   * @param req - Request object
   * @param res - Response object
   */
  async validateUserAction(req: AuthenticatedRequest, res: Response): Promise<void> {
    try {
      const sessionId = req.sessionId;
      const requestBody = req.body as ValidateRequest;

      // Request validation
      if (!requestBody || !requestBody.uuid || !requestBody.user_sig || 
          !requestBody.user_address || !requestBody.project_id || 
          !requestBody.digest || !requestBody.intent) {
        res.status(400).json({
          success: false,
          errorCode: 'INVALID_REQUEST',
          message: 'Missing required fields'
        });
        return;
      }

      // Intent validation
      if (!this.validationService.validateIntent(requestBody.intent)) {
        res.status(400).json({
          success: false,
          errorCode: 'INVALID_INTENT',
          message: 'Invalid intent structure'
        });
        return;
      }

      // Get session ID
      if (!sessionId) {
        res.status(400).json({
          success: false,
          errorCode: 'INVALID_SESSION_ID',
          message: 'Session ID is required'
        });
        return;
      }

      // TODO: We need a defense logic to prevent duplicate UUIDs in requests.
      // Store UUID and SessionID mapping
      try {
        this.db.storeUuidMapping(requestBody.uuid, sessionId);
      } catch (error) {
        console.error('Failed to store UUID mapping:', error);
        res.status(500).json({
          success: false,
          errorCode: 'UUID_MAPPING_FAILED',
          message: 'Failed to store UUID mapping'
        });
        return;
      }

      console.log('ValidateUserActionHandler:', {
        sessionID: sessionId,
        uuid: requestBody.uuid,
        req: JSON.stringify(requestBody)
      });

      // For mint method, validate and deduct assets
      if (requestBody.intent.method === 'mint' || requestBody.intent.method === 'transfer') {
        try {
          this.validationService.validateAndProcessMint(sessionId, requestBody.intent.from);
        } catch (error) {
          console.error('Insufficient balance for sessionId:', sessionId, error);
          res.status(400).json({
            success: false,
            errorCode: 'INSUFFICIENT_BALANCE',
            message: 'Insufficient balance for operation'
          });
          return;
        }
      }

      // Generate validator signature
      let validatorSig: string;
      try {
        validatorSig = await this.validationService.generateValidatorSignature(requestBody.user_sig, requestBody.digest);
      } catch (error) {
        console.error('GenerateValidatorSignature failed:', error);
        res.status(500).json({
          success: false,
          errorCode: 'SIGNATURE_GENERATION_FAILED',
          message: 'Failed to generate validator signature'
        });
        return;
      }

      console.log('validateUserActionHandler:', {
        validatorSig,
        userSig: requestBody.user_sig,
        digest: requestBody.digest
      });

      // Success response
      const response: ValidateResponse = {
        success: true,
        errorCode: null,
        data: {
          userSig: requestBody.user_sig,
          validatorSig
        }
      };

      res.json(response);
    } catch (error) {
      console.error('ValidateUserAction error:', error);
      res.status(500).json({
        success: false,
        errorCode: 'INTERNAL_ERROR',
        message: 'Internal server error'
      });
    }
  }
} 