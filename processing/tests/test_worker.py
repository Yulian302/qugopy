import redis
from redis.exceptions import RedisError
from os import getenv


def test_redis_connection_health():
    REDIS_HOST = getenv("REDIS_HOST", "127.0.0.1")
    REDIS_PORT = int(getenv("REDIS_PORT", "6379"))
    try:
        rdb = redis.Redis(REDIS_HOST, REDIS_PORT, db=0)
        rdb.ping()
    except RedisError as e:
        print(f"❌ Redis connection error: {e}")
        raise
    finally:
        rdb.close()
