# Changelog

Todas as mudanças relevantes deste projeto serão documentadas aqui.

O formato segue [Keep a Changelog](https://keepachangelog.com/pt-BR/1.1.0/) e [SemVer](https://semver.org/lang/pt-BR/).

## [0.3.0] - 2026-07-21
### Added
- Implementação funcional inicial em Go para backend e firmware simulador.
- Endpoints REST v1 de auth, config, status, OTA, heartbeat, eventos e snapshots.
- Runtime de firmware com boot sequence, máquina de estados, widgets MVP, navegação e fallback offline por cache.
- Testes automatizados para contratos, backend e firmware (`go test ./...`).

## [0.2.0] - 2026-07-21
### Added
- Consolidação arquitetural do `docs/TECHNICAL_MASTER_PLAN.md` com seções 27-50.
- Definição de baseline único de NFRs mensuráveis.
- Documento de contratos v1 em `docs/API_MQTT_CONTRACTS_V1.md`.
- Critérios de aceite MVP em `docs/MVP_ACCEPTANCE_CRITERIA.md`.
- Plano de validação em `docs/VALIDATION_PLAN.md`.
- Backlog técnico priorizado em `docs/IMPLEMENTATION_BACKLOG.md`.

## [0.1.0] - 2026-07-21
### Added
- Estrutura inicial do repositório DeskHub.
- README, licença MIT e documentação técnica principal.
