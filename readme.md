# 🏛️ Guia Estrutural: Ports & Adapters (Hexagonal) e DDD

Este documento descreve a arquitetura do projeto `todo-ddd`, focando na separação entre a lógica de negócio (Domínio) e os detalhes técnicos (Infraestrutura), utilizando interfaces para garantir a Inversão de Dependência.

## 1. O Domínio (Núcleo Independente)

O pacote `pkg/domain` é o coração da aplicação. Ele não tem dependências de frameworks ou bancos de dados; ele apenas define a essência do negócio.

| Elemento                | Localização                | Função                                                                                                                                                       |
| :---------------------- | :------------------------- | :----------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **Entidade**            | `pkg/domain/task.go`       | Representa a _Task_. Contém a lógica de modificação do estado (ex: `Update`, `Delete`, métodos `SetStatus`).                                                 |
| **Value Objects (VOs)** | `pkg/domain/valueobject`   | Garantem a validade e a tipagem de valores importantes. Exemplos: `Priority` (1, 2, 3) e `Status` (e.g., "new", "completed").                                |
| **Ports (Interfaces)**  | `pkg/domain/repository.go` | A interface `TaskRepository` é a _Port_ (Porta) de saída. Ela define o contrato de que o Domínio precisa para persistir dados, sem especificar a tecnologia. |

## 2. Camada de Aplicação (Use Cases)

O pacote `pkg/usecase/task` contém a lógica de orquestração e as regras de aplicação.

- **Finalidade:** Coordena as entidades do Domínio e utiliza as _Ports_ (Interfaces) de repositório para realizar operações complexas (Transações, Ações Múltiplas).
- **Implementação:** Cada caso de uso (ex: `CreateTaskUseCase`, `DeleteTaskUseCase`) possui uma injeção da interface `domain.TaskRepository`.
- **Inversão de Dependência:** O _Usecase_ depende da interface `TaskRepository` definida no Domínio (seta para dentro). Isso mantém o Usecase isolado da infraestrutura.

## 3. Infraestrutura (Adapters/Adaptadores)

Os pacotes `internal/adapters` contêm as implementações concretas que se "adaptam" às _Ports_ do Domínio e às necessidades externas.

### 3.1 Adapter de Persistência (Banco de Dados)

- **Localização:** `internal/adapters/db/sqlite`.
- **Função:** O struct `SQLiteTaskRepository` **implementa explicitamente** todos os métodos da interface `domain.TaskRepository` (Ex: `Save`, `Update`, `List`).
- **Detalhe Técnico:** É responsável por toda a interação com o banco de dados SQLite, incluindo as _queries_ SQL e a gestão da conexão (`InitDB`).

### 3.2 Adapter de API (Gin/HTTP)

- **Localização:** `internal/adapters/api/handler` e `internal/adapters/api/router.go`.
- **Função:** O `TaskHandler` recebe os Usecases injetados e atua como uma **Porta de Entrada** (Input Port) da arquitetura:
  - Recebe requisições HTTP.
  - Converte JSON (`CreateTaskRequest`) para a entrada do Usecase (`CreateTaskInput`).
  - Converte a saída do Usecase para a resposta HTTP (`TaskResponse`).

## 🧩 4. Casos de Uso Agregadores e Transações (Onboarding)

O caso de uso SetupOnboardingUseCase (pkg/usecase/setup/setup.go) é um Application Service agregador.
Enquanto os casos de uso de User e Task lidam com operações individuais, o Onboarding coordena ambos em uma única operação transacional.

### 🧠 4.1 Conceito de Agregador

Em DDD, um Agregador de Casos de Uso é um serviço que:
- combina várias operações de aplicação/domínio,
- garante consistência entre agregados (ex: User e Task),
- e aplica regras de orquestração e atomicidade.

Neste projeto, o SetupOnboardingUseCase:
1. Cria um novo usuário (UserAggregate);
2. Cria uma tarefa de boas-vindas (TaskAggregate);
3. Faz commit apenas se ambas as operações forem bem-sucedidas.

## 5. O Ponto de Arranque (Glue Code)

O `cmd/main.go` é o único local que "cola" todas as peças, realizando a Injeção de Dependência final (Service Locator/Container).

1.  Cria o Adaptador de DB concreto: `repo := sqlite.NewSQLiteTaskRepository(db)`.
2.  Injeta esse Adaptador nas estruturas de Usecase: `listUC := &usecase.ListTaskUseCase{TaskRepo: repo}`.
3.  Injeta os Usecases no Adaptador de API (Handler): `taskHandler := &handler.TaskHandler{ListUC: listUC, ...}`.
4.  Inicia o servidor, expondo a API.

| Princípio             | Resumo no Código                                                                                                                                     |
| :-------------------- | :--------------------------------------------------------------------------------------------------------------------------------------------------- |
| **Teste Fácil**       | Para testar um Usecase, basta injetar um **Mock** que implemente `domain.TaskRepository`.                                                            |
| **Acoplamento Baixo** | O Domínio não será afetado se a API mudar de Gin para Fiber, ou o DB mudar de SQLite para Postgres, pois a dependência é sempre na Interface (Port). |
