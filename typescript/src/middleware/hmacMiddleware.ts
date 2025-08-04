import { Request, Response, NextFunction } from 'express';
import { HmacService } from '../services/hmacService';

export const hmacMiddleware = (req: Request, res: Response, next: NextFunction): void => {
  // Skip HMAC validation for GET requests
  if (req.method === 'GET') {
    next();
    return;
  }

  // Get HMAC signature from header
  // Express.js normalizes header names to lowercase
  const hmacSignature = req.headers['x-hmac-signature'] as string;

  // Read request body
  const requestBody = req.body ? JSON.stringify(req.body) : '';

  // Debug: Log the request body and HMAC signature
  console.log('Request Body:', requestBody);
  console.log('HMAC Signature:', hmacSignature);

  // Validate HMAC
  if (!HmacService.validateHmac(requestBody, hmacSignature)) {
    console.log('HMAC validation failed');
    res.status(401).json({
      success: false,
      errorCode: 'INVALID_HMAC_SIGNATURE',
      message: 'Invalid HMAC signature'
    });
    return;
  }

  console.log('HMAC validation passed');

  next();
}; 