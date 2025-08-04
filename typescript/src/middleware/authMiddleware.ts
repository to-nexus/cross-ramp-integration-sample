import { Request, Response, NextFunction } from 'express';

export interface AuthenticatedRequest extends Request {
  sessionId?: string;
}

export const authMiddleware = (req: AuthenticatedRequest, res: Response, next: NextFunction): void => {
  // TODO: In actual implementation, validate JWT tokens properly
  // For now, we'll do basic header validation
  
  const crossAuthToken = req.headers.authorization;
  // Express.js normalizes header names to lowercase
  const dappAuthToken = req.headers['x-dapp-authorization'] as string;
  const sessionId = req.headers['x-dapp-sessionid'] as string;

  // Check if required headers are present
  if (!crossAuthToken || !dappAuthToken || !sessionId) {
    res.status(401).json({
      success: false,
      errorCode: 'INVALID_SESSION_ID',
      message: 'Missing required authentication headers'
    });
    return;
  }

  // TODO: Validate JWT tokens
  // For demo purposes, we'll accept any non-empty values
  if (!crossAuthToken.startsWith('Bearer ') || 
      !dappAuthToken.startsWith('Bearer ') || 
      sessionId.trim() === '') {
    res.status(401).json({
      success: false,
      errorCode: 'INVALID_SESSION_ID',
      message: 'Invalid authentication headers'
    });
    return;
  }

  // Store session ID for later use
  req.sessionId = sessionId;
  next();
}; 