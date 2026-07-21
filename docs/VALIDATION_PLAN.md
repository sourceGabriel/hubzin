# DeskHub OS - Validation Plan v1

Status: Draft

## 1. Objetivo

Validar funcionalidade, estabilidade, desempenho, resiliência e segurança antes de avanço de fase.

## 2. Camadas de validação

### Firmware
- Boot sequence
- Máquina de estados
- Widget lifecycle
- Navegação e input
- Persistência e recovery

### Backend
- Contratos REST/MQTT
- Autenticação
- Cache
- Normalização de dados

### Integração
- Sincronização backend -> firmware
- Entrega de eventos MQTT
- Consistência de configuração remota

### OTA
- Fluxo completo de atualização
- Assinatura e checksum
- Rollback em falha induzida

### Resiliência Offline
- Queda de WiFi
- Queda de backend
- Uso de cache com aviso discreto

## 3. Métricas obrigatórias

- Boot <3 s
- Reconexão WiFi <10 s
- Input <50 ms
- Atualização de widget <100 ms (alvo <50 ms)
- FPS >=30
- RAM <70%
- CPU nominal <40% (pico <60%)

## 4. Critério de aprovação

Cada fase só avança quando:
- todos os testes críticos passam;
- métricas obrigatórias são atingidas;
- regressões críticas são zero;
- documentação foi atualizada.
