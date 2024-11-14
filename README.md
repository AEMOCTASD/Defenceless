
# Defenceless

- This is a bad app. Bad coding. A lot of vulnerabilities.
- You might have to troubleshoot one of the docker files.

## Instalation

```bash
docker-compose up --build
```

## Usage

```bash
curl -X POST -d '{"value":"Oh, NOO!!!!"}' -H "Content-Type: application/json" http://localhost:8080/add
curl http://localhost:8080/get/1
```
