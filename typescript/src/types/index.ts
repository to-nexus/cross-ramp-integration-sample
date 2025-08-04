// TODO: HMAC salt - In actual implementation, load from environment variables or configuration file
export const HMAC_SALT = "my_secret_salt_value_!@#$%^&*";

export interface Asset {
  id: string;
  balance: string;
}

export interface PlayerV1 {
  player_id: string;
  name: string;
  wallet_address: string;
  server: string;
  assets: Asset[];
}

export interface Guide {
  Authorization: string;
  "X-Dapp-Authorization": string;
  "X-Dapp-SessionID": string;
  message: string;
  session_info: {
    created_at: string;
    updated_at: string;
  };
}

export interface AssetsResponse {
  success: boolean;
  errorCode: string | null;
  data: {
    v1: PlayerV1;
    guide: Guide;
  };
}

export interface ValidateRequestIntent {
  method: string;
  from: Array<{
    type: string;
    id: string;
    amount: number;
  }>;
  to: Array<{
    type: string;
    id: string;
    amount: number;
  }>;
}

export interface ValidateRequest {
  uuid: string;
  user_sig: string;
  user_address: string;
  project_id: string;
  digest: string;
  intent: ValidateRequestIntent;
}

export interface ValidateResponseData {
  userSig: string;
  validatorSig: string;
}

export interface ValidateResponse {
  success: boolean;
  errorCode: string | null;
  data: ValidateResponseData;
}

export interface ExchangeResultIntent {
  outputs: Array<{
    asset_id: string;
    amount: number;
  }>;
}

export interface ExchangeResultRequest {
  uuid: string;
  intent: ExchangeResultIntent;
  receipt: {
    status: number;
  };
}

export interface ExchangeResultResponse {
  success: boolean;
}

export interface SessionAssets {
  sessionId: string;
  assets: Asset[];
  createdAt: Date;
  updatedAt: Date;
}

export interface UuidMapping {
  uuid: string;
  sessionId: string;
  createdAt: Date;
}

export interface ErrorResponse {
  success: false;
  errorCode: string;
  message?: string;
}

export type ErrorCode = 
  | "INVALID_REQUEST"
  | "INVALID_SESSION_ID"
  | "DB_ERROR"
  | "INVALID_INTENT"
  | "UUID_MAPPING_FAILED"
  | "INSUFFICIENT_BALANCE"
  | "SIGNATURE_GENERATION_FAILED"
  | "INVALID_HMAC_SIGNATURE"; 