#/bin/bash!

pkill -9 redis-server

rm -rf data/redis/6380
rm -rf data/redis/6381
rm -rf data/redis/6382
rm -rf data/redis/6383
rm -rf data/redis/6384
rm -rf data/redis/6385

mkdir -p data/redis/6380
mkdir -p data/redis/6381
mkdir -p data/redis/6382
mkdir -p data/redis/6383
mkdir -p data/redis/6384
mkdir -p data/redis/6385

touch data/redis/6380/redis.log
touch data/redis/6381/redis.log
touch data/redis/6382/redis.log
touch data/redis/6383/redis.log
touch data/redis/6384/redis.log
touch data/redis/6385/redis.log

chmod 0755 data/redis/6380/redis.log
chmod 0755 data/redis/6381/redis.log
chmod 0755 data/redis/6382/redis.log
chmod 0755 data/redis/6383/redis.log
chmod 0755 data/redis/6384/redis.log
chmod 0755 data/redis/6385/redis.log

# redis-server config/redis-6380.conf
# redis-server config/redis-6381.conf
# redis-server config/redis-6382.conf
# redis-server config/redis-6383.conf
# redis-server config/redis-6384.conf
# redis-server config/redis-6385.conf

# redis-cli --cluster create 127.0.0.1:6380 127.0.0.1:6381 127.0.0.1:6382 127.0.0.1:6383 127.0.0.1:6384 127.0.0.1:6385 --cluster-replicas 1