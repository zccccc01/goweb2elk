# goweb2elk

Use docker to run goweb and elk(elasticsearch, logstash, kibana)

This is a rough template

# Usage

```bash
docker-compose up --build -d
curl localhost:8080/ # test goweb
curl localhost:8080/error
curl -X GET "localhost:9200/_cat/indices?v" # test elasticsearch if see .ds-logs-generic-default-xxx then success
```

### Create your visual chart

Use browser to visit `http://localhost:5601`

1. `Management` -> `Stack Management`
2. `Kibana` -> `Data Views`
3. `Create data view`
4. `Name` -> `logs`
5. `Index pattern` -> `logs*`
6. `Time field` -> `@timestamp`
7. `Visualize` -> `Create visualization` -> `Lens`
