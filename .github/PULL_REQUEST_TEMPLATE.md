name: Pull Request

description: Template padrão para PRs do DeskHub.

body:
  - type: textarea
    id: summary
    attributes:
      label: Resumo
      description: O que este PR altera?
      placeholder: Descreva as mudanças principais.
    validations:
      required: true

  - type: textarea
    id: motivation
    attributes:
      label: Motivação
      description: Por que esta mudança é necessária?
    validations:
      required: true

  - type: checkboxes
    id: checklist
    attributes:
      label: Checklist
      options:
        - label: Li e segui a arquitetura/documentação do projeto
        - label: Incluí ou atualizei documentação
        - label: Testei localmente
