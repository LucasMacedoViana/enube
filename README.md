# Enube - Importador de Dados e API com Go + PostgreSQL

Este projeto tem como objetivo a construção de uma aplicação em **Go (Golang)** para:

✅ Importar grandes volumes de dados de um arquivo Excel para um banco de dados **PostgreSQL**,  
✅ Expor uma **API REST** com autenticação via **JWT** e endpoints de consulta,  
✅ (Opcional) Fornecer um **Dashboard em React** com indicadores, agrupamentos e métricas de cobrança.

---

## 🏗️ Arquitetura

- **Linguagem:** Golang 1.20+
- **Framework HTTP:** [Fiber](https://github.com/gofiber/fiber)
- **ORM:** GORM
- **Banco de Dados:** PostgreSQL
- **Autenticação:** JWT
- **Documentação:** Swagger (via `swag`)
- **Validação:** Go Validator
- **Importador:** Leitura e persistência de `.xlsx` usando `excelize`
- **Segurança:** Hash de senha com bcrypt
- **Estrutura:** Padrões **SOLID**, **Clean Architecture**, e separação em camadas `domain`, `infra`, `handler`, `middleware`, `utils`.

---

## 📦 Funcionalidades

### 1. Importador de Planilha Excel
- Lê arquivos `.xlsx` com grande volume de dados
- Normaliza os dados em tabelas:
    - Parceiros
    - Clientes
    - Produtos
    - Assinaturas
    - Medidores
    - Itens de faturamento
- Cria os registros em lote usando `CreateInBatches`
- Valida e evita duplicações com `FirstOrCreate`
- Registra progresso e erros em arquivos `.log`

### 2. API RESTful
- `/api/login` - Autenticação de usuários com JWT
- `/api/users` - Cadastro e listagem de usuários
- `/api/dashboard/summary` - Totais de faturamento e quantidade
- `/api/dashboard/by-category` - Agrupamento por categoria
- (em expansão) Endpoints para agrupamentos por cliente, recurso e mês

### 3. Segurança
- Autenticação JWT
- Middleware de proteção para rotas privadas
- Hash seguro de senhas
- Swagger com suporte a Bearer Token

---

## 🚀 Como executar

### Pré-requisitos
- Docker + Docker Compose
- Go 1.20+ (se for rodar localmente)

### Subir a aplicação com Docker
```bash
docker-compose up --build
```

### Para rodar localmente
1. Instale as dependências:
   ```bash
   go mod tidy
   ```
2. Configure as variáveis de ambiente no arquivo `.env`:
   ```dotenv
    JWT_SECRET=secret
    DATABASE_DSN=postgres://root:root@localhost:5432/enube
    ADMIN_NAME=admin
    ADMIN_PASSWORD=123456
    PORT=3000
    ```
3. Execute a aplicação:
4. ```bash
   go run cmd/main.go
   ```
5. Acesse a API em `http://localhost:3000/swagger/index.html` para ver a documentação Swagger.
---
## 👥 Usuários padrão
{
  "name": "admin",
  "password": "123456"
}
---
## 🎲 Banco de Dados
![imagem](./ENUBBE.png)

--

