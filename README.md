# 🌦 Weather API (Go + Visual Crossing + Redis)

A simple, fast Weather API built in **Go** that proxies the [Visual Crossing Weather API](https://www.visualcrossing.com/weather-api), caches results in **Redis**, and supports rate limiting with `httprate`.

This project demonstrates:

* Fetching data from a **3rd-party API**
* Using **environment variables** for configuration
* Implementing **in-memory caching** with Redis
* Handling **HTTP requests/responses** in Go
* Adding **rate limiting** to prevent abuse

---

## 📌 Features

* **Proxy to Visual Crossing** — Hide your API key and add your own business logic
* **Redis caching** — Avoid repeated API calls for the same city/unit
* **Configurable TTL** — Cache duration set via `.env`
* **Rate limiting** — Defaults to `60 requests/min/IP`
* **Environment variables** — Keeps secrets out of code
* **Health check** — `/health` endpoint
* **Unit tests** — For cache, config, client, and handlers

---

## 🗂 Project structure

```
.
├── cmd/
│   └── server/          # Application entrypoint (main.go)
├── internal/
│   ├── cache/           # Redis caching logic
│   ├── config/          # Env var loading
│   ├── http/            # HTTP handlers & middleware
│   └── weather/         # Visual Crossing API client
├── .env.example         # Sample environment variables
├── docker-compose.yml   # Redis container setup
└── README.md
```

---

## 🚀 Getting Started

### 1️⃣ Clone the repo

```bash
git clone https://github.com/Khoa-Trinh/weather-api-api.git
cd weather-api
```

### 2️⃣ Install dependencies

```bash
go mod tidy
```

### 3️⃣ Set up `.env`

Copy `.env.example` → `.env` and fill in your values:

```dotenv
PORT=8080
VISUAL_CROSSING_API_KEY=your_visual_crossing_key
DEFAULT_UNIT_GROUP=metric

REDIS_ADDR=127.0.0.1:6379
REDIS_PASSWORD=
REDIS_DB=0
CACHE_TTL_SECONDS=43200
HTTP_TIMEOUT_SECONDS=10
```

You can get a free API key from: [Visual Crossing Sign Up](https://www.visualcrossing.com/weather/weather-data-services)

---

### 4️⃣ Start Redis (Docker)

```bash
docker compose up -d
```

Or install Redis locally:

* macOS: `brew install redis && brew services start redis`
* Ubuntu: `sudo apt install redis-server`
* Windows: [Memurai](https://www.memurai.com/) or [Redis on WSL2](https://redis.io/docs/latest/operate/oss_and_stack/install/install-redis/install-redis-on-windows/)

---

### 5️⃣ Run the server

```bash
go run ./cmd/server
```

You should see:

```
{"level":"INFO","msg":"starting server","port":"8080"}
```

---

## 📡 API Endpoints

### Health check

```bash
curl http://localhost:8080/health
```

**Response:**

```json
{"status":"ok"}
```

---

### Get weather data

```bash
curl "http://localhost:8080/weather?city=Hanoi&units=metric"
```

**Response:**

```json
{
  "queryCost": 1,
  "latitude": 21.0,
  "longitude": 105.85,
  "resolvedAddress": "Hanoi",
  "address": "Hanoi",
  "timezone": "Asia/Bangkok",
  "tzoffset": 7.0,
  "description": "Similar temperatures continuing with a chance of rain...",
  "days": [
    { "datetime": "2025-08-11", "temp": 30.5, "conditions": "Partially cloudy" }
  ]
}
```

**Headers:**

* `X-Cache: HIT` — Returned from Redis cache
* `X-Cache: MISS` — Freshly fetched from Visual Crossing

---

## 🧪 Running tests

```bash
go test ./...
```

Tests include:

* Config loading
* Cache behavior (using miniredis)
* Weather API client
* HTTP handlers

---

## 📜 License

MIT License © 2025 Khoa Trinh -- see [LICENSE](LICENSE) for details.

---

## 🔗 Links

* Visual Crossing: https://www.visualcrossing.com/
* Go Chi: https://github.com/go-chi/chi
* Redis Go client: https://github.com/redis/go-redis
* Roadmap Project: https://roadmap.sh/projects/weather-api-wrapper-service
* Related Roadmap Guide: https://roadmap.sh/projects/weather-api-wrapper-service