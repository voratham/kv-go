# Distribute datastore


```sh
# run dev
 go run main.go --db-location=$PWD/my.db --config-file=$PWD/sharding.toml --shard=Bangkok

# run production ready 3 cluster
./launch.sh 
 
# mock data
curl http://localhost:8080/set\?key\=a\&value\=1
curl http://localhost:8080/set\?key\=b\&value\=2
curl http://localhost:8080/set\?key\=c\&value\=3
curl http://localhost:8080/set\?key\=d\&value\=4

# check data
curl http://localhost:8080/get\?key\=a
curl http://localhost:8080/get\?key\=b
curl http://localhost:8080/get\?key\=c
curl http://localhost:8080/get\?key\=d
```

