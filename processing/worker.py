
import json
import sys
from typing import Any, Dict, Union
import redis
import time
import grpc
import signal
import logging
from google.protobuf.empty_pb2 import Empty
from os import getenv, path
from pydantic import BaseModel, field_validator
from dotenv import load_dotenv

import task_pb2_grpc
from handlers.image_processor import handle_task


def shutdown_handler(signum, frame):
    print("👋 Received shutdown signal", flush=True)
    sys.exit(0)


signal.signal(signal.SIGINT, shutdown_handler)
signal.signal(signal.SIGTERM, shutdown_handler)


class Task(BaseModel):
    type: str
    payload: Union[bytes, Dict[str, Any]]
    priority: int
    # deadline: Optional[datetime] = None
    # recurring: Optional[bool] = False

    @field_validator('payload', mode='before')
    def convert_payload(cls, v):
        if isinstance(v, dict):
            return json.dumps(v).encode('utf-8')  # Convert dict to bytes
        return v


class IntTask(BaseModel):
    id: str
    task: Task


def wait_for_grpc_ready(channel):
    for attempt in range(5):
        try:
            grpc.channel_ready_future(channel).result(timeout=3)
            return True
        except grpc.FutureTimeoutError:
            logging.warning("Waiting for gRPC server...")
            time.sleep(1)
    return False


class Worker:
    def __init__(self, rdb=None, is_local=True):
        self.rdb = rdb
        self.is_local = is_local
        if is_local:
            channel = grpc.insecure_channel("localhost:50051")
            if not wait_for_grpc_ready(channel):
                print("❌ gRPC server never became ready", flush=True)
                sys.exit(1)
            self.stub = task_pb2_grpc.TaskServiceStub(channel)

    def process_task(self, int_task: IntTask):
        task_type = int_task.task.type
        if task_type == "process_image":
            logging.info(handle_task(payload=int_task.task.payload))
        else:
            return

    def run(self):
        while True:
            if self.is_local:
                try:
                    task: IntTask = self.stub.GetPythonTask(
                        Empty(), timeout=5)
                    self.process_task(task)
                except grpc.RpcError as e:
                    if e.code() == grpc.StatusCode.NOT_FOUND:
                        logging.info("No task in queue")
                        time.sleep(0.2)
                    elif e.code() == grpc.StatusCode.UNAVAILABLE:
                        print("Server unavailable, retrying...")
                        channel = grpc.insecure_channel("localhost:50051")
                        self.stub = task_pb2_grpc.TaskServiceStub(channel)
                        time.sleep(1)
            else:
                try:
                    task_data = rdb.zpopmin("python_queue", 1)
                    if task_data:
                        raw, _ = task_data[0]
                        task_dict = json.loads(raw)
                        task = IntTask(**task_dict)
                        self.process_task(task)
                    else:
                        logging.info("No task! Sleeping...")
                        time.sleep(0.2)
                except Exception as e:
                    logging.error(f"❌ Error processing task: {e}")
                    time.sleep(1)


if __name__ == "__main__":
    load_dotenv(path.abspath(path.join(path.dirname(__file__), "..", ".env")))
    is_production = getenv("IS_PRODUCTION", "false").lower() == "true"
    if is_production:
        logging.basicConfig(
            level=logging.INFO,
            format='%(asctime)s - %(levelname)s - %(message)s',
            handlers=[
                logging.FileHandler('workers.log'),
            ]
        )
    else:
        logging.basicConfig(
            level=logging.DEBUG,
            format='%(asctime)s - %(levelname)s - %(message)s',
            stream=sys.stdout
        )

    MODE = getenv("MODE", "local").lower()

    if MODE == "redis":
        REDIS_HOST = getenv("REDIS_HOST", "127.0.0.1")
        REDIS_PORT = int(getenv("REDIS_PORT", "6379"))
        rdb = redis.Redis(REDIS_HOST, REDIS_PORT, db=0)
        worker = Worker(rdb=rdb, is_local=False)
    else:
        worker = Worker(is_local=True)

    logging.info(
        f"🚀 Starting worker in {'LOCAL' if MODE != "redis" else 'REDIS'} mode...")
    worker.run()
