# 🚀 Geradores Concorrentes em Go com Arquitetura Hexagonal

![Go Version](https://img.shields.io/badge/Go-1.22%2B-blue.svg)
![License](https://img.shields.io/badge/License-MIT-green.svg)
![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)

Este é um projeto de exemplo robusto que demonstra a implementação de geradores simultâneos em Go, seguindo as melhores práticas de engenharia de software, incluindo **Arquitetura Hexagonal (Ports & Adapters)**, princípios **SOLID**, e **Observabilidade** com OpenTelemetry.

O objetivo é servir como um guia prático para desenvolvedores Go que desejam construir aplicações concorrentes, testáveis, escaláveis e de fácil manutenção.

---

## ✨ Principais Características

* **Arquitetura Hexagonal (Ports & Adapters):** Separação clara entre a lógica de negócio (core) e as implementações externas (adaptadores), promovendo baixo acoplamento e alta testabilidade.
* **Princípios SOLID:** O código é estruturado seguindo os cinco princípios SOLID para criar um software mais compreensível, flexível e manutenível.
* **Padrões de Concorrência Avançados:** Utilização de canais, goroutines, `sync.WaitGroup` e `context` para gerenciar um pipeline de processamento de dados de forma segura e eficiente.
* **Observabilidade com OpenTelemetry:** Integração nativa de tracing para rastrear o ciclo de vida de cada "pedido" através das goroutines, desde a geração até o processamento final.
* **Graceful Shutdown:** A aplicação pode ser encerrada de forma limpa, garantindo que nenhum trabalho em andamento seja perdido.

---

## 🏛️ Visão Geral da Arquitetura

A aplicação segue o padrão de Arquitetura Hexagonal. A lógica de negócio principal (`core`) é completamente isolada e se comunica com o mundo exterior através de **Ports** (interfaces). As implementações concretas, como os geradores de números e o sistema de log, são **Adapters** que se conectam a essas portas.

```
+---------------------------------+      +--------------------+      +-------------------------+
|      Adapters (Mundo Externo)   |      |       Ports        |      |    Core / Domínio       |
|                                 |      |   (Interfaces Go)  |      |  (Lógica de Negócio)    |
|  - Geradores de Pedidos         |      |                    |      |                         |
|  - Logger (Stdout)              |----->| - OrderProcessor   |<-----| - ProcessingService     |
|  - (Poderia ser Kafka, DB, etc) |      |                    |      | - domain.Order          |
+---------------------------------+      +--------------------+      +-------------------------+
```

Essa estrutura permite trocar facilmente um adaptador por outro (por exemplo, substituir o logger de `Stdout` por um que envia para o `Datadog`) sem alterar uma única linha da lógica de negócio.

---

## 📁 Estrutura do Projeto

```
/concurrent-generators
├── /cmd                      # Entrypoints da aplicação
│   └── /main
├── /internal                 # Código privado do projeto
├── /adapters                 # Implementações concretas (o "como")
│   ├── /generator
│   └── /logger
│   └── /core                 # Lógica de negócio (o "o quê")
│       ├── /domain
│       ├── /ports
│       └── /services
├── go.mod
└── README.md
```

---

## 🚀 Como Executar

### Pré-requisitos

* [Go](https://go.dev/doc/install) (versão 1.22 ou superior)

### Passos

1.  **Clone o repositório:**
    ```bash
    git clone [https://github.com/br4tech/concurrent-generators.git](https://github.com/br4tech/concurrent-generators.git)
    ```
    
    ```bash
    cd concurrent-generators
    ```

2.  **Execute a aplicação:**
    O Go se encarregará de baixar as dependências (`go.mod`) e compilar o projeto.

    ```bash
    go run ./cmd/main.go
    ```

---

## 📊 Entendendo a Saída

Ao executar a aplicação, você observará duas coisas principais no seu terminal:

1.  **Logs de Processamento:** Mensagens em tempo real mostrando os geradores e workers sendo iniciados, processando pedidos de forma concorrente e finalizando de forma organizada.

    ```
    INFO: Iniciando a aplicação...
    INFO: Gerador 1 iniciado
    INFO: Worker 1 iniciado
    INFO: Worker 2 iniciado
    INFO: Gerador 2 iniciado
    ...
    INFO: Worker 1 recebeu o pedido: 5
    INFO: Worker 2 recebeu o pedido: 0
    INFO: Pedido processado com sucesso: ID 0, Valor 0
    ...
    INFO: Todos os geradores finalizaram. Canal de pedidos fechado.
    INFO: Worker 1 finalizando.
    INFO: Worker 2 finalizando.
    INFO: Aplicação finalizada com sucesso.
    ```

2.  **Dados de Trace (OpenTelemetry):** Após a finalização, o exportador de trace imprimirá no console um relatório detalhado (em formato JSON) de todos os spans criados. Isso permite visualizar a hierarquia, duração e metadados de cada operação.

    *(Exemplo de um span)*
    ```json
    {
      "Name": "ProcessOrder",
      "SpanContext": { ... },
      "Parent": { ... },
      "StartTime": "...",
      "EndTime": "...",
      "Attributes": [
        { "Key": "worker.id", "Value": { "Type": "INT64", "Value": 2 } },
        { "Key": "order.id", "Value": { "Type": "INT64", "Value": 10 } }
      ],
      ...
    }
    ```

---

## 💡 Conceitos Demonstrados

* **Pipeline Concorrente:** `generator.GenerateOrders` atua como o **produtor**, criando um canal de saída. O `processingSvc.StartWorkers` lança múltiplos **consumidores** (workers) que leem deste mesmo canal.
* **Fan-Out / Fan-In:** Embora não seja um fan-in explícito, o padrão onde múltiplos geradores (fan-out) escrevem em um único canal para ser consumido por múltiplos workers é uma variação poderosa desse conceito.
* **Injeção de Dependência:** Na função `main`, as dependências concretas (adaptadores) são criadas e injetadas no serviço do core, que depende apenas de interfaces. Isso é fundamental para a Arquitetura Hexagonal e o princípio D do SOLID.
* **Contexto para Cancelamento:** O `context.Context` é propagado por toda a pilha de chamadas para garantir que, ao receber um sinal de interrupção (Ctrl+C), todas as goroutines em andamento sejam notificadas para parar o que estão fazendo e encerrar de forma limpa.

---

## 🤝 Contribuições

Contribuições são sempre bem-vindas! Sinta-se à vontade para abrir uma *issue* para discutir novas funcionalidades ou relatar bugs. Se desejar contribuir com código, por favor, abra um *Pull Request*.
