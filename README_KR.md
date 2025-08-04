# Sample Game Backend API

이 프로젝트는 Language별 Ramp 연동 가이드를 위한 샘플 백엔드입니다.

자산정보 조회, 검증, 온체인 결과 수신에 대한 후처리 내용을 담고 있고 

그 외 처리는 간소화되어 있습니다. (예를 들어 자산 정보 조회시 메모리 상에서 임시 생성해서 결과 노출되도록 코드가 작성되어 있습니다. 실제 구현 시 미리 구현되어 있는 DB 상의 재화가 노출되도록 구현해야 합니다.)

API의 요청, 응답에 주목하고 고려사항은 TODO 코멘트를 참고해 주세요.

## ⚠️ 중요 안내사항

**이 코드는 Ramp 연동을 위한 참고용 구현체로만 제공됩니다.**

- 이 샘플 코드는 Ramp 프로토콜과의 연동 패턴을 보여주기 위해 설계되었습니다
- **각 게임사별 요구사항에 맞게 수정이 필요할 수 있습니다**
- **이 코드 사용으로 인한 문제나 손해에 대한 책임은 사용자에게 있습니다**
- 이 구현체는 프로덕션 환경에서 사용할 준비가 되지 않았으며, 배포 전 철저한 테스트와 커스터마이징이 필요합니다
- 작성자는 이 샘플 코드 사용으로 인한 결과에 대해 책임지지 않습니다

## 지원 언어

이 프로젝트는 다음 언어로 구현되어 있습니다:

- **Go** (Go 1.23.0) - Golang 구현체
- **C#** (.NET 8.0) - ASP.NET Core 기반
- **C++** (C++17) - CMake 기반

## 주요 기능
- CROSSRAMP 연동을 위한 가상의 게임서버입니다.
- 인게임 자산 조회 API
- 유저의 주문 검증 API
- 온체인 결과 후처리 API

## API 스펙

### 1. 자산 정보 조회 API

#### 요청
```bash
curl -X GET "http://localhost:8080/api/assets?language=ko" \
  -H "Authorization: Bearer <CROSS_AUTH_JWT>" \
  -H "X-Dapp-Authorization: Bearer <DAPP_ACCESS_TOKEN>" \
  -H "X-Dapp-SessionID: <DAPP_SESSION_ID>"
```

#### 응답 (status=200)
```json
{
  "success": true,
  "errorCode": null,
  "data": {
    "v1": {
      "player_id": "C1",            // 고유 키
      "name": "playerName_C1",
      "wallet_address": "0xaaaa",   // 매핑된 지갑 주소
      "server": "test",
      "assets": [
        {
          "id": "asset_money",      // ramp console에 등록된 asset id
          "balance": "2000"
        },
        {
          "id": "item_gem",
          "balance": "1500"
        }
      ]
    },
    "guide": {
      "Authorization": "Bearer <token>",
      "X-Dapp-Authorization": "Bearer <token>",
      "X-Dapp-SessionID": "<session_id>",
      "message": "guide 필드는 요청 시 헤더 정보를 표기합니다. 올바르게 게임사와 프로토콜을 맞췄는지 확인하는 용도이고 게임사에는 제공되지 않습니다. ramp frontend 개발자 참고용입니다.",
      "session_info": {
        "created_at": "2024-01-01T00:00:00Z",
        "updated_at": "2024-01-01T00:00:00Z"
      }
    }
  }
}
```

### 2. 주문 검증 API

#### 요청
```bash
curl -X POST "http://localhost:8080/api/validate" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <CROSS_AUTH_JWT>" \
  -H "X-Dapp-Authorization: Bearer <DAPP_ACCESS_TOKEN>" \
  -H "X-Dapp-SessionID: <DAPP_SESSION_ID>" \
  -d '{
    "user_sig": "0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
    "user_address": "0xB777C937fa1afC99606aFa85c5b83cFe7f82BabD",
    "project_id": "acjviwejsi",
    "intent": {
      "method": "mint",
      "from": [
        { "type": "asset", "id": "asset_money", "amount": 1000 }
      ],
      "to": [
        { "type": "asset", "id": "item_gem", "amount": 1000 }
      ]
    }
  }'
```

#### 응답 (status=200)
```json
{
  "success": true,
  "errorCode": null,
  "data": {
    "userSig": "0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
    "validatorSig": "0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"
  }
}
```

### 3. 교환 결과 처리 API

#### 요청
```bash
curl -X POST "http://localhost:8080/api/result" \
  -H "Content-Type: application/json" \
  -d '{
    "uuid": "exchange-uuid-123",
    "intent": {
      "outputs": [
        { "asset_id": "item_gem", "amount": 1000 }
      ]
    },
    "receipt": {
      "status": 1
    }
  }'
```

#### 응답 (status=200)
```json
{
  "success": true
}
```

## 설치 및 실행

### Go 버전 (권장)

#### 필수 요구사항
- Go 1.23.0 이상
- Git

#### 1. 저장소 클론
```bash
git clone <repository-url>
cd sample-game-backend/golang
```

#### 2. 의존성 설치
```bash
go mod tidy
```

#### 3. 서버 실행
```bash
go run main.go
```

서버가 포트 8080에서 시작됩니다. 세션별 자산 정보는 **go-memdb** 인메모리 데이터베이스에 저장됩니다.

### C# 버전

#### 필수 요구사항
- .NET 8.0 SDK
- Git

#### 1. 저장소 클론
```bash
git clone <repository-url>
cd sample-game-backend/csharp/SampleGameBackend
```

#### 2. 의존성 복원
```bash
dotnet restore
```

#### 3. 서버 실행
```bash
dotnet run
```

### C++ 버전

#### 필수 요구사항
- CMake 3.16 이상
- C++17 호환 컴파일러
- nlohmann/json 라이브러리
- OpenSSL
- Git

#### 1. 저장소 클론
```bash
git clone <repository-url>
cd sample-game-backend/cpp
```

#### 2. 빌드
```bash
mkdir build && cd build
cmake ..
make
```

#### 3. 서버 실행
```bash
./sample_game_backend_cpp
```

## API 엔드포인트

| 메서드 | 엔드포인트 | 설명 |
|--------|------------|------|
| GET | `/api/assets` | 자산 정보 조회 |
| POST | `/api/validate` | 주문 검증 |
| POST | `/api/result` | 교환 결과 처리 |
| GET | `/health` | 헬스체크 |


## 에러 코드

| 에러 코드 | 설명 |
|-----------|------|
| INVALID_REQUEST | 잘못된 요청 형식 또는 필수 필드 누락 |
| INVALID_SESSION_ID | 세션 ID 누락 또는 잘못됨 |
| DB_ERROR | 데이터베이스 작업 실패 |
| INVALID_INTENT | 잘못된 인텐트 구조 또는 메서드 |
| UUID_MAPPING_FAILED | UUID 매핑 저장 실패 |
| INSUFFICIENT_BALANCE | 작업에 대한 자산 잔액 부족 |
| SIGNATURE_GENERATION_FAILED | 검증자 서명 생성 실패 |

## 테스트

### 수동 테스트 예제

#### 자산 정보 조회
```bash
curl -X GET "http://localhost:8080/api/assets?language=ko" \
  -H "Authorization: Bearer test_cross_auth_jwt_token" \
  -H "X-Dapp-Authorization: Bearer test_dapp_access_token" \
  -H "X-Dapp-SessionID: test_session_id"
```

#### 주문 검증
```bash
curl -X POST "http://localhost:8080/api/validate" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer test_cross_auth_jwt_token" \
  -H "X-Dapp-Authorization: Bearer test_dapp_access_token" \
  -H "X-Dapp-SessionID: test_session_id" \
  -d '{
    "user_sig": "0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
    "user_address": "0xB777C937fa1afC99606aFa85c5b83cFe7f82BabD",
    "project_id": "acjviwejsi",
    "intent": {
      "method": "mint",
      "from": [
        { "type": "asset", "id": "asset_money", "amount": 1000 }
      ],
      "to": [
        { "type": "asset", "id": "item_gem", "amount": 1000 }
      ]
    }
  }'
```

#### 헬스체크
```bash
curl -X GET "http://localhost:8080/health"
```

## 프로젝트 구조

```
sample-game-backend/
├── golang/                 # Go 구현체 (메인)
├── csharp/               # C# 구현체
├── cpp/                  # C++ 구현체
```
