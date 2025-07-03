
import json
import redis
from os import getenv
from typing import Optional
import time
import grpc
from pydantic import BaseModel
from datetime import datetime

import task_pb2
import task_pb2_grpc


class Task(BaseModel):
    type: str
    payload: str
    priority: int
    # deadline: Optional[datetime] = None
    # recurring: Optional[bool] = False


class IntTask(BaseModel):
    id: str
    task: Task


class Worker:
    def __init__(self, rdb=None, is_local=True):
        self.rdb = rdb
        self.is_local = is_local
        if is_local:
            channel = grpc.insecure_channel("localhost:50051")
            self.stub = task_pb2_grpc.TaskServiceStub(channel)

    def process_task(self, int_task: IntTask):
        print(
            f"✅ Processing task {int_task.id} - Payload: {int_task.task.payload}")

    def run(self):
        while True:
            if self.is_local:
                try:
                    task: IntTask = self.stub.GetTask(
                        task_pb2.Empty(), timeout=2)
                    self.process_task(task)
                    continue
                except grpc.RpcError as e:
                    if e.code() == grpc.StatusCode.NOT_FOUND:
                        print("No task in queue")
                    else:
                        raise
                    time.sleep(1)
                    continue
            try:
                task_data = rdb.brpop("task_queue", timeout=5)

                if task_data:
                    _, raw = task_data
                    task: Task = json.loads(raw)
                    self.process_task(task)
                else:
                    print("No task! Sleeping...")
                    time.sleep(1)
            except Exception as e:
                print(f"❌ Error processing task: {e}")
                time.sleep(1)


if __name__ == "__main__":
    MODE = getenv("MODE", "local").lower()
    is_local = MODE != "redis"

    w = Worker()

    if is_local:
        worker = Worker(is_local=True)
    else:
        REDIS_HOST = getenv("REDIS_HOST", "127.0.0.1")
        REDIS_PORT = int(getenv("REDIS_PORT", "6379"))
        rdb = redis.Redis(REDIS_HOST, REDIS_PORT, db=0)
        worker = Worker(rdb=rdb, is_local=False)

    print(f"🚀 Starting worker in {'LOCAL' if is_local else 'REDIS'} mode...")
    worker.run()
