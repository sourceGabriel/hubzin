# DeskHub
## Technical Master Plan
Version: 0.1
Status: Draft
Owner: Gabriel Martins

---

# 1. Objetivo

DeskHub é um painel inteligente open source baseado em ESP32 capaz de centralizar informações pessoais, profissionais e de automação residencial em um único dispositivo físico.

Não é apenas um display.

O objetivo é criar uma plataforma modular para dashboards embarcados.

---

# 2. Princípios do Projeto

- Offline First
- Open Source
- Modular
- Extensível
- Atualização OTA
- Arquitetura desacoplada
- Hardware reproduzível
- Fácil manutenção
- Alta documentação
- Widgets independentes

---

# 3. Escopo

### MVP

✓ Hora

✓ Data

✓ Clima

✓ Agenda

✓ CPU do computador

✓ Spotify

✓ Troca de páginas

✓ Configuração WiFi

✓ OTA

---

Versão 1

Adicionar

GitHub

Home Assistant

Bambu Studio

Notificações

News

Todoist

Gmail

Modo foco

Pomodoro

---

Versão 2

Marketplace de Widgets

Bluetooth

Matter

Zigbee Gateway

Touch Screen

PCB própria

---

# 4. Arquitetura

Camadas

Usuário

↓

Display

↓

Widget Engine

↓

Core Firmware

↓

Communication Layer

↓

MQTT/REST

↓

Backend

↓

External APIs

---

# 5. Hardware

Inicial

ESP32-S3

Display IPS SPI

Encoder Rotativo

2 Botões

LED RGB

Buzzer

Sensor luminosidade

USB-C

Fonte 5V

---

Futuro

RTC

Sensor Temperatura

Sensor Umidade

Microfone

Alto Falante

BLE Beacon

Bateria

---

# 6. Firmware

Módulos

Boot Manager

Display Manager

Page Manager

Widget Manager

Theme Manager

Network Manager

OTA Manager

Storage Manager

MQTT Manager

REST Client

Logger

Scheduler

Power Manager

Input Manager

---

# 7. Widget Engine

Cada Widget possui

ID

Nome

Versão

Dependências

Tempo atualização

Prioridade

Cache

Layout

Eventos

Estado

---

Widgets nunca conhecem outros Widgets.

Toda comunicação acontece pelo Event Bus.

---

# 8. Widgets previstos

Clock

Calendar

Spotify

Weather

CPU

RAM

GPU

Printer

GitHub

Todoist

Home Assistant

News

Stock

Crypto

Network

Temperature

Humidity

Email

Notifications

Pomodoro

RSS

---

# 9. Layout

Home

Sistema

Música

Trabalho

Casa

Impressora

Mercado

Configurações

---

Cada página suporta

Grid

Lista

Cards

Gráfico

Texto

Imagem

---

# 10. Navegação

Encoder

Gira

↓

Troca página

Clique

↓

Seleciona

Clique longo

↓

Menu

Botão lateral

↓

Voltar

---

# 11. Comunicação

Firmware nunca consulta APIs públicas.

Sempre consulta Backend.

Backend agrega todas informações.

Benefícios

Menor consumo

Maior segurança

Menor uso memória

Cache

Escalabilidade

---

# 12. Backend

Serviços

Gateway

Dashboard

Weather

Calendar

Spotify

Github

Printer

Notification

Authentication

---

Banco

PostgreSQL

Redis

MQTT Broker

---

# 13. APIs

REST

Configuração

Status

OTA

Logs

Widgets

MQTT

Atualizações

Eventos

Notificações

Heartbeat

OTA

---

# 14. Atualização OTA

Servidor publica

Nova versão

↓

Firmware verifica assinatura

↓

Download

↓

Validação

↓

Instalação

↓

Rollback automático caso falhe

---

# 15. Segurança

HTTPS

JWT

TLS MQTT

OTA assinada

Tokens criptografados

Segredos fora firmware

Configuração protegida

---

# 16. Armazenamento

Flash

Configuração

Tema

Widgets ativos

Idioma

WiFi

Timezone

MQTT

Backend

---

# 17. Interface

Dark Theme

Light Theme

Auto Brightness

60 FPS quando possível

Animações

Fade

Slide

Loading

Toast

Modal

Popup

---

# 18. Design

Minimalista

Tipografia grande

Pouca informação por tela

Contraste elevado

Ícones vetoriais

Paleta neutra

---

# 19. Performance

Baseline mensurável (v1)

Boot

<3 segundos

FPS

>=30 FPS (alvo: 60 FPS quando possível)

RAM

Uso <70% (equivalente a RAM livre >30%)

CPU

Uso nominal <40% (picos controlados <60%)

Tempo de atualização de widget

SLO <100 ms (alvo operacional <50 ms)

Latência de input

<50 ms

Reconexão WiFi

<10 segundos

---

# 20. Consumo

Tela desligada

Modo baixo consumo

WiFi Sleep

Atualização dinâmica

Brilho automático

---

# 21. Logs

INFO

WARNING

ERROR

DEBUG

Persistência circular

Exportação futura

---

# 22. Roadmap

Fase 1

ESP32

Display

Hora

Clima

OTA

---

Fase 2

Backend

MQTT

Widgets

---

Fase 3

Spotify

GitHub

Home Assistant

---

Fase 4

Caixa 3D

PCB

---

Fase 5

Marketplace

---

# 23. Estrutura do Repositório

docs/

firmware/

backend/

frontend/

hardware/

pcb/

3d/

docker/

.github/

assets/

---

# 24. Documentação

Vision

PRD

SRS

Architecture

ADR

Hardware

Electrical

Firmware

Backend

Testing

Deployment

Assembly

Maintenance

Roadmap

---

# 25. Critérios de Qualidade

Sem dependências acopladas

Documentação antes do código

Cobertura de testes

Versionamento semântico

Arquitetura limpa

Baixo consumo

Fácil expansão

---

# 26. Futuro

Aplicativo Android

Aplicativo iOS

Dashboard Web

Plugin SDK

Marketplace

Widgets comunitários

Integração Matter

Assistente de Voz

IA Local

Suporte Raspberry Pi

Versão Linux

---

---

# 27. Arquitetura do Sistema

## 27.1 Visão Geral

DeskHub será dividido em quatro grandes domínios independentes.

```text
+------------------------------------------------------+
|                    Usuário                           |
+-------------------------+----------------------------+
                          |
                  Interface Física
                          |
+------------------------------------------------------+
|                     Firmware                         |
+------------------------------------------------------+
| Core | Widgets | Display | Input | Network | Storage |
+------------------------------------------------------+
                          |
                    MQTT / HTTPS
                          |
+------------------------------------------------------+
|                     Backend                          |
+------------------------------------------------------+
| API | Cache | Integrations | Auth | Notification     |
+------------------------------------------------------+
                          |
                  Serviços Externos
```

Cada camada possui responsabilidade única.

Toda comunicação deve ocorrer através de contratos bem definidos.

---

# 28. Arquitetura do Firmware

O firmware será dividido em três níveis.

## Core

Responsável pelo funcionamento do dispositivo.

Responsabilidades

- Inicialização
- Boot
- Watchdog
- Gerenciamento de memória
- Sistema de eventos
- Scheduler
- OTA
- Storage
- Configuração

## Services

Serviços compartilhados

- Display
- WiFi
- MQTT
- HTTP
- Filesystem
- RTC
- Power
- Logger
- Animation
- Theme
- Input

## Features

Tudo que o usuário vê

- Widgets
- Páginas
- Menus
- Notificações
- Popup
- Overlay

---

# 29. Boot Sequence

```text
Power On
↓
POST
↓
GPIO Check
↓
Filesystem
↓
Configuration
↓
Display Init
↓
Load Theme
↓
Load Widgets
↓
WiFi
↓
MQTT
↓
Backend Sync
↓
Ready
```

Tempo máximo esperado: <3 segundos.

---

# 30. Máquina de Estados

Estados possíveis

- BOOTING
- IDLE
- SYNCING
- UPDATING
- SETUP
- SLEEP
- ERROR
- RECOVERY

Toda transição deverá possuir timeout.

---

# 31. Sistema de Widgets

Todo widget deverá obedecer ao mesmo contrato.

Estados

- Created
- Loading
- Ready
- Paused
- Updating
- Hidden
- Destroyed

Cada widget possui

- UUID
- Nome
- Descrição
- Versão
- Autor
- Dependências
- Permissões
- Tempo de atualização
- Prioridade
- Memória utilizada
- Estado
- Eventos publicados
- Eventos escutados

---

# 32. Ciclo de Vida dos Widgets

```text
Create
↓
Initialize
↓
Load Resources
↓
First Render
↓
Visible
↓
Update
↓
Sleep
↓
Wake
↓
Destroy
```

Atualizações devem ocorrer em background.

---

# 33. Sistema de Eventos

Arquitetura Publish/Subscribe.

Nenhum widget pode chamar outro diretamente.

Eventos previstos

- TimeChanged
- MinuteChanged
- HourChanged
- WeatherUpdated
- NotificationReceived
- SpotifyChanged
- PageChanged
- ThemeChanged
- BrightnessChanged
- WifiConnected
- WifiDisconnected
- MQTTConnected
- MQTTDisconnected
- SleepMode
- WakeUp
- OTAStarted
- OTAFinished

---

# 34. Sistema de Páginas

Páginas base

- Home
- Trabalho
- Música
- Casa
- Sistema
- Impressora
- Financeiro
- Configuração

Cada página possui

- Nome
- Ícone
- Widgets
- Layout
- Prioridade
- Permissões
- Animação de entrada
- Animação de saída

---

# 35. Sistema de Layout

Tipos

- Grid
- Columns
- Rows
- Cards
- Full Screen
- Scrollable
- Dynamic

O Layout Engine calcula automaticamente o posicionamento.

---

# 36. Design Responsivo

Suporte esperado a múltiplos displays

- 128x64 OLED
- 240x240 TFT
- 320x240 TFT
- 480x320 IPS
- 800x480 (futuro)

Nenhum widget deve assumir resolução fixa.

---

# 37. Sistema de Temas

Tema define

- Cores
- Ícones
- Fontes
- Espaçamento
- Animações
- Bordas
- Sombras
- Gradientes

Tema pode ser alterado sem reinicialização.

---

# 38. Sistema de Fontes

Categorias

- Tiny
- Small
- Normal
- Large
- Title
- Display

Widgets usam categoria e não tamanho absoluto.

---

# 39. Sistema de Ícones

Biblioteca própria.

Vetorial quando possível com fallback bitmap.

Categorias

- System
- Weather
- Music
- Notification
- Network
- Printer
- Calendar
- Settings
- Home
- IoT

---

# 40. Sistema de Navegação

Entradas

- Encoder
- Botões
- Touch (futuro)
- BLE Remote (futuro)
- Gestos (futuro)

Nenhuma ação crítica pode ocorrer em um único clique.

---

# 41. Configuração Inicial

Fluxo

```text
Primeiro Boot
↓
Modo AP
↓
Portal Captive
↓
Idioma
↓
WiFi
↓
Timezone
↓
Backend
↓
MQTT
↓
Tema
↓
Download Configuração
↓
Finalização
```

Após configuração inicial, o AP deve ser desativado.

---

# 42. Sistema de Configuração

Categorias

- Rede
- Tela
- Som
- Widgets
- Integrações
- Sistema
- Desenvolvedor
- Backup

Toda configuração possui

- Valor padrão
- Tipo
- Validação
- Descrição
- Persistência
- Versão

---

# 43. Persistência

Categorias

- Configuração
- Preferências
- Cache
- Estado
- Sessão
- Logs

Usar camada de abstração de storage.

---

# 44. Estratégia de Cache

Cache local com TTL independente por domínio (clima, agenda, notícias, imagens, ícones, tempo e Spotify).

---

# 45. Sincronização

Sincronização nunca bloqueia renderização.

Pode ocorrer em:

- Inicialização
- Mudança de WiFi
- Intervalo
- Evento
- Solicitação do usuário

---

# 46. Tratamento de Erros

Categorias

- Hardware
- Rede
- Backend
- Widget
- Renderização
- OTA
- Storage

Cada erro deve possuir

- Código
- Descrição
- Origem
- Nível
- Solução sugerida
- Possibilidade de recuperação

---

# 47. Recovery

Caso backend indisponível

- Entrar em modo offline
- Utilizar cache
- Exibir aviso discreto
- Tentar reconectar

Nunca reiniciar automaticamente por perda de conexão.

---

# 48. Watchdog

Monitorar

- Loop principal
- Render
- WiFi
- MQTT
- Scheduler
- OTA

Tempo máximo configurável por módulo.

---

# 49. Scheduler

Prioridades

- Realtime
- High
- Normal
- Low
- Background

Input possui prioridade máxima.
Renderização possui prioridade superior a widgets.

---

# 50. Objetivos de Qualidade (Baseline de Aceite)

- Disponibilidade: >=99%
- Boot: <3 segundos
- Reconexão WiFi: <10 segundos
- Atualização de widget: SLO <100 ms (alvo <50 ms)
- Uso de RAM: <70%
- Uso de CPU: nominal <40%, pico <60%
- FPS: mínimo 30, alvo 60 quando possível
- Resposta de input: <50 ms

---

# Filosofia

DeskHub não deve ser apenas um projeto de ESP32.

Deve ser uma plataforma para painéis inteligentes, onde hardware, firmware e backend evoluam de forma independente, permitindo novos dispositivos, novos widgets e novas integrações sem necessidade de reescrever o núcleo do sistema.
