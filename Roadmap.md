# Roadmap de Atividades para o Projeto de Serviço de Táxi

Este roadmap segue a abordagem BDD, garantindo que para cada funcionalidade:
1. Primeiro definimos os comportamentos esperados (feature files)
2. Desenvolvemos testes unitários
3. Implementamos o código necessário
4. Validamos com testes E2E
5. Refinamos conforme necessário

Dessa forma, o desenvolvimento é guiado pelos comportamentos esperados do sistema, melhorando a qualidade e garantindo que as expectativas dos stakeholders sejam atendidas.

## 1. Configuração Inicial

- [ ] Configurar estrutura base do backend (Go + Fiber)
- [ ] Configurar Godog para BDD no backend
- [ ] Configurar banco de dados PostgreSQL
- [ ] Inicializar projeto frontend (Svelte + TypeScript)
- [ ] Configurar Cypress com cucumber-preprocessor
- [ ] Definir estrutura de diretórios para ambos os projetos

## 2. Especificação de Comportamentos (Feature Files)

- [ ] Definir cenários para autenticação do motorista
- [ ] Definir cenários para gerenciamento de perfil do motorista
- [ ] Definir cenários para o fluxo de corridas
- [ ] Definir cenários para avaliações e histórico
- [ ] Definir cenários para notificações
- [ ] Revisar cenários com stakeholders

## 3. Backend - Desenvolvimento baseado em BDD

### Módulo de Autenticação
- [ ] Escrever feature files em Gherkin (Godog)
- [ ] Implementar testes unitários para autenticação
- [ ] Desenvolver modelo e repositório de motorista
- [ ] Implementar serviço de autenticação (registro, login, recuperação)
- [ ] Criar controladores e rotas de API
- [ ] Implementar middleware de autenticação
- [ ] Validar com testes E2E

### Módulo de Gerenciamento de Motoristas
- [ ] Escrever feature files em Gherkin
- [ ] Implementar testes unitários para validações
- [ ] Desenvolver modelo completo de motorista
- [ ] Implementar serviço de validação (idade, CPF, CNH)
- [ ] Criar CRUD completo de motoristas
- [ ] Desenvolver vinculação de veículo
- [ ] Validar com testes E2E

### Módulo de Corridas
- [ ] Escrever feature files para fluxo de corridas
- [ ] Implementar testes unitários para regras de corrida
- [ ] Desenvolver modelo de corridas
- [ ] Implementar serviço gerador de corridas mock
- [ ] Criar endpoints de listagem e aceitação
- [ ] Desenvolver lógica de cancelamento e finalização
- [ ] Implementar cálculo de tempo estimado
- [ ] Validar com testes E2E

### Módulo de Avaliações e Histórico
- [ ] Escrever feature files para avaliações
- [ ] Implementar testes unitários para regras de avaliação
- [ ] Desenvolver modelos de avaliação e histórico
- [ ] Criar serviços de histórico e filtros
- [ ] Implementar endpoints de API
- [ ] Validar com testes E2E

### Módulo de Notificações
- [ ] Escrever feature files para notificações
- [ ] Implementar testes unitários para regras de notificação
- [ ] Desenvolver sistema de notificações
- [ ] Criar endpoints de configuração
- [ ] Implementar lógica de envio de alertas
- [ ] Validar com testes E2E

## 4. Frontend - Desenvolvimento baseado em BDD

### Módulo de Autenticação
- [ ] Escrever cenários Cucumber para frontend
- [ ] Implementar testes unitários de componentes
- [ ] Desenvolver tela de login e registro
- [ ] Implementar fluxo de recuperação de senha
- [ ] Criar serviço de gerenciamento de tokens
- [ ] Validar com testes E2E Cypress

### Dashboard do Motorista
- [ ] Escrever cenários Cucumber
- [ ] Implementar testes unitários
- [ ] Desenvolver componentes da página principal
- [ ] Criar componente de status online/offline
- [ ] Implementar visualização de estatísticas
- [ ] Validar com testes E2E

### Gestão de Corridas
- [ ] Escrever cenários Cucumber
- [ ] Implementar testes unitários
- [ ] Desenvolver tela de corridas disponíveis
- [ ] Criar interface para corrida em andamento
- [ ] Implementar controles de aceitação/cancelamento
- [ ] Desenvolver visualização de tempo e rota
- [ ] Validar com testes E2E

### Histórico e Avaliações
- [ ] Escrever cenários Cucumber
- [ ] Implementar testes unitários
- [ ] Desenvolver página de histórico
- [ ] Criar componente de avaliação
- [ ] Implementar filtros e visualização detalhada
- [ ] Validar com testes E2E

### Notificações e Configurações
- [ ] Escrever cenários Cucumber
- [ ] Implementar testes unitários
- [ ] Desenvolver componente de notificações
- [ ] Criar página de configurações
- [ ] Implementar integração com backend
- [ ] Validar com testes E2E

## 5. Integração e Validação Final

- [ ] Implementar integração completa entre frontend e backend
- [ ] Executar testes E2E do fluxo completo
- [ ] Corrigir bugs e inconsistências
- [ ] Realizar testes de aceitação com stakeholders
- [ ] Refinar componentes e otimizar desempenho

## 6. Finalização

- [ ] Otimizar consultas ao banco de dados
- [ ] Verificar compatibilidade cross-browser
- [ ] Preparar documentação final
- [ ] Realizar testes de carga (opcional)
- [ ] Preparar ambiente para demonstração

