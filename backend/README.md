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

## Regras de validação de payload

- Endpoints de escrita (`/v1/auth/token`, `/heartbeat`, `/events`) rejeitam campos desconhecidos.
- Endpoints de escrita rejeitam payload com múltiplos objetos JSON no mesmo body.
- Endpoints de escrita exigem `Content-Type: application/json`.
- `heartbeat` exige `firmwareVersion` não vazio e `cpuUsage`/`ramUsage` no intervalo `0..100`.
- `events` exige `schemaVersion` e `type` não vazios.
