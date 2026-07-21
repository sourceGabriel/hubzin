# DeskHub OS - MVP Acceptance Criteria

Status: Draft

## Escopo MVP

- Hora
- Data
- Clima
- Agenda
- CPU do computador
- Spotify
- Troca de páginas
- Configuração WiFi
- OTA

## Definition of Done por item

### Hora e Data
- Renderização em tela inicial
- Atualização em intervalo correto
- Sem travamento por 24h de execução contínua

### Clima
- Dados vindos do backend
- Cache local em falha de rede
- Tempo de atualização conforme configuração

### Agenda
- Dados normalizados pelo backend
- Exibição ordenada por horário
- Fallback em modo offline com última sincronização

### CPU do computador
- Telemetria recebida via backend/MQTT
- Exibição com atualização periódica
- Tratamento de valor ausente sem crash

### Spotify
- Atualização por evento
- Exibição de estado (tocando/pausado)
- UI não bloqueante durante mudança de faixa

### Troca de páginas
- Navegação por encoder/botões
- Tempo de resposta de input <50 ms
- Sem perda de estado crítico do widget

### Configuração WiFi
- Fluxo de setup inicial funcional
- Persistência de credenciais
- Reconexão automática <10 s após retorno da rede

### OTA
- Verificação de assinatura
- Download e validação
- Rollback automático em falha

## Gate de aceite do MVP

- Boot <3 s
- FPS mínimo 30
- RAM <70%
- CPU nominal <40% (picos <60%)
- Disponibilidade funcional >=99% em teste de soak
