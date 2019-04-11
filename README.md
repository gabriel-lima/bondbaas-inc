# Architecture
<img src="https://docs.google.com/drawings/d/e/2PACX-1vTGAeyLDzCKMvUlcnHbxA_OpakyjTFM658v1kPKVhz68xm_Sg3FiyITDWZU3nquaNuLSpTcEvumFjMv/pub?w=564&amp;h=840">

# Dependencies in Production
- [Postgres driver](https://github.com/lib/pq)

# Dependencies in Development
- [Fresh](https://github.com/gravityblast/fresh) to hot-reloading applications in Go

# Building application
`docker-compose build`

# Running application
`docker-compose up`

# Running tests
`docker-compose run app go test ./... -v`

# Acessing application container through sh
`docker-compose run app sh`

# Acessing database container
`docker-compose exec db psql -d bondbaas-db -U bondbaas`
