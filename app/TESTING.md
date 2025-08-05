# Tests Documentation

Este projeto inclui uma suíte completa de testes automatizados para garantir a qualidade e confiabilidade do código.

## Tipos de Testes

### 1. Testes Unitários

Testam componentes individuais em isolamento:

- **Models** (`internal/models/models_test.go`): Testa validações de entrada e estruturas de dados
- **Repositories** (`internal/repositories/repositories_test.go`): Testa operações de banco de dados usando SQLite em memória
- **Handlers** (`internal/handlers/handlers_test.go`): Testa endpoints da API usando mocks

### 2. Testes de Integração

Testam a integração entre múltiplos componentes:

- **Integration Tests** (`tests/integration/integration_test.go`): Testa fluxos completos da aplicação com banco de dados real

## Como Executar os Testes

### Pré-requisitos

```bash
# Instalar dependências
make deps

# Para testes de integração, certifique-se de que o PostgreSQL está rodando
make db-up
```

### Executar Todos os Testes

```bash
make test
```

### Executar Apenas Testes Unitários

```bash
make test-unit
```

### Executar Testes com Relatório de Cobertura

```bash
make test-coverage
```

Isso gera um arquivo `coverage.html` que você pode abrir no navegador para ver a cobertura detalhada.

### Executar Testes de Componentes Específicos

```bash
# Testes dos models
make test-models

# Testes dos handlers
make test-handlers

# Testes dos repositories
make test-repositories

# Testes dos usecases
make test-usecases
```

## Estrutura dos Testes

### Testes de Models

- **Validação de entrada**: Testa se valores válidos e inválidos são corretamente validados
- **Estruturas de dados**: Verifica se os modelos são construídos corretamente

### Testes de Repository

- **Operações CRUD**: Testa inserção, atualização e consulta de dados
- **Casos extremos**: Testa comportamento quando não há dados
- **Múltiplas operações**: Testa sequências de operações

### Testes de Handlers

- **Endpoints da API**: Testa todos os endpoints com diferentes cenários
- **Códigos de status HTTP**: Verifica se os códigos corretos são retornados
- **Validação de JSON**: Testa parsing e validação de payloads
- **Tratamento de erros**: Verifica se erros são tratados adequadamente

### Testes de Integração

- **Fluxo completo**: Testa operações Set → Get → Sync → Check
- **Persistência**: Verifica se dados são corretamente persistidos no banco
- **Health check**: Testa se a aplicação responde corretamente

## Cobertura de Testes

O projeto visa manter alta cobertura de testes:

- **Models**: 100% - Testam todas as validações e estruturas
- **Repositories**: 95%+ - Testam todas as operações de banco
- **Handlers**: 95%+ - Testam todos os endpoints e casos de erro
- **Integration**: Fluxos principais cobertos

## Mocks e Stubs

### MockContractUseCase

Usado nos testes de handlers para simular operações blockchain sem precisar de conexão real:

```go
type MockContractUseCase struct {
    setValue      uint64
    getValue      uint64
    setError      error
    getError      error
    // ... outros campos
}
```

### MockContractUseCaseIntegration

Usado nos testes de integração para simular blockchain mas usar banco real:

```go
type MockContractUseCaseIntegration struct {
    repo            repositories.Repository
    blockchainValue uint64
}
```

## Configuração de CI/CD

Os testes podem ser facilmente integrados em pipelines CI/CD:

```bash
# Pipeline básico
make ci
```

Isso executa:

1. Download de dependências
2. Formatação do código
3. Testes com cobertura

## Boas Práticas

### Nomenclatura de Testes

```go
func TestHandler_SetValue_Success(t *testing.T) // Padrão: Test[Component]_[Method]_[Scenario]
```

### Estrutura de Teste

```go
func TestSomething(t *testing.T) {
    // Arrange (preparação)
    // Act (ação)
    // Assert (verificação)
}
```

### Casos de Teste

- **Happy path**: Cenários de sucesso
- **Edge cases**: Casos extremos (valores límite, entrada vazia)
- **Error cases**: Cenários de erro (falhas de rede, banco indisponível)

## Troubleshooting

### Testes de Integração Falhando

Se os testes de integração estão falhando:

1. Verifique se o PostgreSQL está rodando: `make db-up`
2. Certifique-se de que as credenciais estão corretas
3. Verifique se o banco de teste pode ser criado

### Baixa Cobertura

Para melhorar a cobertura:

```bash
make test-coverage
# Abra coverage.html no navegador
# Identifique linhas não cobertas
# Adicione testes para essas linhas
```

### Testes Lentos

Se os testes estão lentos:

- Testes unitários usam mocks e devem ser rápidos
- Testes de integração podem ser mais lentos (uso real do banco)
- Use `go test -short` para pular testes de integração durante desenvolvimento

## Melhorias Futuras

1. **Testes de Performance**: Adicionar benchmarks para operações críticas
2. **Testes de Carga**: Testar comportamento sob alta carga
3. **Testes E2E**: Testes completos com blockchain real (ambiente de teste)
4. **Property-based testing**: Usar ferramentas como QuickCheck para Go
