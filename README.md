<center>
<h1>Kipa</h1>
<p>A secure and serverless key-value store</p>
</center>

## Uses 
- Go
- Vue.js
- MongoDB
- Redis Queues
- Sendgrid
- Docker
- Swagger

## Requirements
- Go (>1.17)
- Node.js 
- Docker

## Using Kipa
> Fill in the environment variables in the `.env` file

1. To run Kipa using Docker compose 
```bash
$ make run-docker-compose
```

2. To run unit tests
```bash
$ make run-unit-tests
```

3. To run e2e tests
```bash
$ make run-server-e2e-tests
```
Web client: http://localhost:3000 
Backend: http://localhost:5050

## Documentation
Swagger documentation is available at http://localhost:5050/docs/index.html 