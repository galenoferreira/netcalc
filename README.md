

# 🌐 CIDRCalc

Ferramenta completa para cálculo de sub-redes IPv4. Disponível como API REST, interface web e com monitoramento via Prometheus e Grafana.

## 📦 Componentes

- 🖥️ **Backend:** API REST em Go
- 🌐 **Frontend:** Next.js (React)
- 📊 **Monitoramento:** Prometheus + Grafana
- 🐳 **Orquestração:** Docker Compose

---

## 🚀 Como rodar o projeto

### Pré-requisitos:
- Docker
- Docker Compose

### ⚙️ Subindo tudo:

```bash
docker-compose up --build
```

- API disponível em: [http://localhost:8080](http://localhost:8080)
- Frontend em: [http://localhost:3000](http://localhost:3000)
- Prometheus em: [http://localhost:9090](http://localhost:9090)
- Grafana em: [http://localhost:3001](http://localhost:3001)

---

## 🔗 Endpoints Backend

### `POST /calculate`
Calcula informações de sub-rede.

**Request JSON:**
```json
{
  "ip": "192.168.1.10",
  "cidr": 24
}
```

**Response JSON:**
```json
{
  "ip": "192.168.1.10",
  "cidr": 24,
  "mask": "255.255.255.0",
  "network": "192.168.1.0",
  "broadcast": "192.168.1.255",
  "first_ip": "192.168.1.1",
  "last_ip": "192.168.1.254",
  "hosts": 254
}
```

### `GET /metrics`
Exporta métricas Prometheus.

---

## 📊 Monitoramento

- **Prometheus:** Coleta métricas da API no endpoint `/metrics`.
- **Grafana:** Dashboards para acompanhar:
  - Total de requisições
  - Tempo médio das requisições
  - Uptime da API

---

## 🗂️ Estrutura de pastas

```
.
├── backend/           # API em Go
│   ├── Dockerfile
│   └── main.go
├── frontend/          # Interface Next.js
│   ├── Dockerfile
│   └── package.json
├── monitoring/        # Configuração do Prometheus
│   └── prometheus.yml
├── docker-compose.yml # Orquestração dos containers
└── README.md          # Este arquivo
```

---

## 🐳 Docker Compose

- `backend` → Porta `8080`
- `frontend` → Porta `3000`
- `prometheus` → Porta `9090`
- `grafana` → Porta `3001`

---

## 🛠️ Scripts úteis

Subir o ambiente:

```bash
docker-compose up --build
```

Derrubar o ambiente:

```bash
docker-compose down
```

---

## 📜 Licença

Este projeto está licenciado sob a licença MIT.

---

## ✨ Contribuição

Contribuições são bem-vindas! Crie uma issue, envie sugestões ou abra um pull request. 💪