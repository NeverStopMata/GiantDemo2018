docker run -d -p 6379:6379 -v "$PWD":/data --restart always --name mope_redis redis redis-server /data/redis.conf
