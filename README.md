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
