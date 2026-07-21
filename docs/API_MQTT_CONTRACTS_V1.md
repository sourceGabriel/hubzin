# DeskHub OS - API e MQTT Contracts v1

Status: Draft

## 1. Princípios de Contrato

- Firmware consome apenas Backend.
- Todo payload possui versão (`schemaVersion`).
- Contratos são retrocompatíveis dentro da mesma major version.
- Erros são padronizados com código e mensagem.

## 2. REST v1

### 2.1 Autenticação

- `POST /v1/auth/token`
- Entrada: credenciais de dispositivo
- Saída: `accessToken`, `expiresAt`, `refreshToken`

### 2.2 Configuração de Dispositivo

- `GET /v1/devices/{deviceId}/config`
- Retorna: tema, páginas, widgets ativos, integrações habilitadas e intervalos de atualização.

### 2.3 Estado do Sistema

- `GET /v1/devices/{deviceId}/status`
- Retorna: conectividade, versão firmware, uso de recursos e saúde do dispositivo.

### 2.4 OTA

- `GET /v1/devices/{deviceId}/ota/latest`
- Retorna: versão, checksum, assinatura, URL de download e política de rollout.

### 2.5 Widgets

- `GET /v1/widgets/{widgetId}/snapshot`
- Retorna estado normalizado para renderização imediata.

## 3. MQTT v1

## 3.1 Convenção de tópicos

- `deskhub/v1/devices/{deviceId}/heartbeat`
- `deskhub/v1/devices/{deviceId}/events`
- `deskhub/v1/devices/{deviceId}/notifications`
- `deskhub/v1/devices/{deviceId}/widgets/{widgetId}/update`
- `deskhub/v1/devices/{deviceId}/commands`

## 3.2 Payload base

```json
{
  "schemaVersion": "1.0",
  "timestamp": "2026-07-21T00:00:00Z",
  "deviceId": "string",
  "type": "string",
  "data": {}
}
```

## 3.3 Heartbeat mínimo

- `uptimeSec`
- `firmwareVersion`
- `wifiConnected`
- `mqttConnected`
- `cpuUsage`
- `ramUsage`

## 4. Segurança

- REST via HTTPS obrigatório.
- MQTT via TLS obrigatório.
- Tokens nunca persistidos em texto plano.
- Segredos fora do firmware e fora do repositório.

## 5. Versionamento

- Mudanças incompatíveis exigem `v2` de endpoint/tópico.
- Mudanças compatíveis incrementam minor de `schemaVersion`.

## 6. Modelo de Erro

```json
{
  "error": {
    "code": "CONFIG_NOT_FOUND",
    "message": "Configuration not found for device",
    "traceId": "string"
  }
}
```
