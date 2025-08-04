# Sample Game Backend - C# .NET Version

이 프로젝트는 Ramp Protocol과의 통합을 위한 게임 백엔드 API의 C# .NET 구현체입니다.


## 🛠️ Installation

### Prerequisites
- .NET 8.0 SDK
- Visual Studio 2022 또는 VS Code

### Setup
1. 프로젝트 디렉토리로 이동:
```bash
cd csharp/SampleGameBackend
```

2. 의존성 복원:
```bash
dotnet restore
```

3. 애플리케이션 실행:
```bash
dotnet run
```

## 📁 Project Structure

```
SampleGameBackend/
├── Controllers/          # API 컨트롤러
│   ├── AssetsController.cs
│   ├── ValidationController.cs
│   └── ExchangeController.cs
├── Models/              # 데이터 모델
│   ├── Asset.cs
│   ├── Exchange.cs
│   └── Response.cs
├── Services/            # 비즈니스 로직
│   ├── KeystoreService.cs
│   ├── ValidationService.cs
│   └── ExchangeService.cs
├── Database/            # 데이터베이스 서비스
│   └── DatabaseService.cs
├── Middleware/          # 미들웨어
│   └── AuthMiddleware.cs
├── Program.cs           # 애플리케이션 진입점
├── appsettings.json     # 설정 파일
└── SampleGameBackend.csproj
```

## 🔧 Configuration

### Environment Variables
- `ASPNETCORE_ENVIRONMENT`: 실행 환경 (Development/Production)
- `ASPNETCORE_URLS`: 서버 URL (기본값: http://localhost:8080)


## DEBUG
```
export DOTNET_ROOT="/opt/homebrew/opt/dotnet@8/libexec" && export PATH="/opt/homebrew/opt/dotnet@8/bin:$PATH" && dotnet run --urls "http://localhost:8080"
```