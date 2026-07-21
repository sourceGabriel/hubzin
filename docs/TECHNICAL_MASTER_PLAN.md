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

Boot

<3 segundos

Tela

60 FPS

RAM livre

>30%

Uso CPU

<40%

Tempo atualização Widget

<50 ms

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

# Filosofia

DeskHub não deve ser apenas um projeto de ESP32.

Deve ser uma plataforma para painéis inteligentes, onde hardware, firmware e backend evoluam de forma independente, permitindo novos dispositivos, novos widgets e novas integrações sem necessidade de reescrever o núcleo do sistema.
