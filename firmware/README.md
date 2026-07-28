# Firmware

Código embarcado do DeskHub (ESP32-S3).

## Simulador funcional local

```bash
go run ./cmd/deskhub-firmware
```

## Funcionalidades implementadas no simulador

- Boot sequence e máquina de estados
- Setup/sync/idle/recovery
- Widgets MVP (clock, date, weather, agenda, cpu, spotify)
- Navegação de páginas
- Cache local com fallback offline
- Heartbeat periódico para backend
