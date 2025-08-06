## **1. Microserviço: Customers**

**Descrição**: Gerencia autenticação, perfis de usuário e autorização.

### Tecnologias:

- **Backend**: Go (Chi/REST + gRPC)
- **Banco de Dados**: PostgreSQL
- **Cache**: Redis
- **Autenticação**: JWT/OAuth2

### Requisitos Funcionais (RF):

- RF01: Cadastro de usuários com e-mail/senha ou OAuth2 (Google, Facebook).
- RF02: Autenticação via JWT com refresh tokens.
- RF03: Gerenciamento de perfis (atualizar nome, foto, endereço).
- RF04: Controle de roles (admin, usuário comum).
- RF05: Histórico de atividades do usuário (logins, alterações).
- RF06: Recuperação de senha via e-mail/SMS.

### Requisitos Não Funcionais (RNF):

- RNF01: Latência máxima de 200ms em 99% das requisições.
- RNF02: Autenticação com criptografia TLS 1.3.
- RNF03: Rate limiting (ex: 100 requisições/min por IP).
- RNF04: Cache de sessões em Redis (TTL de 1h).
- RNF05: Escalabilidade horizontal via Kubernetes.

---

## **Comunicação entre Serviços**

### gRPC:

- Comunicação interna (ex: Orders → Payments).
- Protocol Buffers para contratos rigorosos.

### REST (Chi):

- APIs públicas (ex: Frontend → Users).
- Documentação via Swagger/OpenAPI.

### GraphQL:

- Consultas complexas no frontend (ex: buscar pedidos com detalhes de produtos).

---

## **Tecnologias Complementares**

- **Autenticação**: JWT com [`golang-jwt`](https://github.com/golang-jwt/jwt).
- **Gateway de API**: [`grpc-gateway`](https://github.com/grpc-ecosystem/grpc-gateway).
- **Monitoramento**: Prometheus + Grafana.
- **Logs**: ELK Stack (Elasticsearch, Logstash, Kibana).
- **DevOps**: Docker Compose (local), Kubernetes (produção).

---

## **Desafios Técnicos**

1. **Concorrência**: Goroutines para processamento paralelo (ex: Recommendations).
2. **Performance**: Otimização de queries SQL (ex: evitar N+1 no Orders).
3. **Segurança**: Validação de inputs com [`validator`](https://github.com/go-playground/validator).
4. **Testes**: Testes de integração com `dockertest`.
