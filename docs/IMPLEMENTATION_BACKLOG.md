# DeskHub OS - Implementation Backlog (Priorizado)

Status: Draft

## Fase A - Core Firmware Foundation
1. Boot manager + state machine
2. Storage abstraction + config schema versioning
3. Event bus + scheduler priorities
4. Display/input base
5. Network manager (WiFi + reconnect)

## Fase B - Backend Gateway Foundation
1. Auth service de dispositivo
2. Config endpoint v1
3. Status endpoint v1
4. OTA metadata endpoint v1
5. MQTT broker topics v1 e heartbeat

## Fase C - Widget Runtime MVP
1. Widget contract e lifecycle
2. Page/layout engine mínimo
3. Widgets MVP: Clock/Date/Weather/Agenda/CPU/Spotify
4. Cache strategy por widget
5. Error handling e modo offline

## Fase D - OTA + Hardening
1. Assinatura e validação OTA
2. Rollback automático
3. Watchdog por módulo
4. Logs estruturados por nível
5. Testes de soak e estabilidade

## Fase E - Release Gates
1. Rodar plano de validação completo
2. Verificar NFRs de baseline
3. Atualizar documentação oficial
4. Gerar release notes e tag semântica
