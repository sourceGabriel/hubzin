# Backend

Serviços de agregação e APIs para o DeskHub.

## Execução local

```bash
go run ./cmd/deskhub-backend
```

## Endpoints v1 implementados

- `POST /v1/auth/token`
- `GET /v1/devices/{deviceId}/config`
- `GET /v1/devices/{deviceId}/status`
- `GET /v1/devices/{deviceId}/ota/latest`
- `POST /v1/devices/{deviceId}/heartbeat`
- `POST /v1/devices/{deviceId}/events`
- `GET /v1/widgets/{widgetId}/snapshot`

## Autenticação de dispositivo (simulador local)

`POST /v1/auth/token` valida `deviceId` e `deviceSecret` contra credenciais em memória.

- `dev1` / `x`
- `demo-device` / `local-only`
