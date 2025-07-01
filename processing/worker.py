
import json
import redis
from os import getenv
import time
from pydantic import BaseModel


from dotenv import load_dotenv


class Task(BaseModel):
    id: str
    type: str
    payload: str


class Worker:
    def __init__(self):
        pass

    def process_task(self, task: Task):
        print(f"✅ Processing task {task['id']} - Payload: {task['payload']}")

    def run(self):
        while True:
            task_data = rdb.brpop("task_queue", timeout=5)
            if task_data:
                _, raw = task_data
                task: Task = json.loads(raw)
                self.process_task(task)
            else:
                print("No task! Sleeping...")
                time.sleep(1)


if __name__ == "__main__":
    load_dotenv("../.env")

    rdb = redis.Redis(getenv("REDIS_HOST"), getenv("REDIS_PORT"), db=0)
    w = Worker()
    w.run()
