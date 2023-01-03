"""
Файл-обертка для запуска сервиса, выполняющего задачи.
Основное назначение - запуск нового экземпляра сервиса в ответ на алерт;
также выставляет кол-во доступных сервису воркеров в зависимости от данных в алерте.

Для работы требуются переменные окружения:
    * ROOT_DIR - корневая директория, в которой лежит Makefile, с помощью которого запускается сервис
    * MERGED_CONFIGS_FILE - абсолютный путь до файла с настройками
"""

import json
import os
from typing import Mapping

from flask import Flask
from flask_httpauth import HTTPBasicAuth


root_dir = os.environ.get('ROOT_DIR')
merged_configs_file = os.environ.get('MERGED_CONFIGS_FILE')
assert (
    merged_configs_file and root_dir,
    'Unable to work without MERGED_CONFIGS_FILE and ROOT_DIR in env'
)

with open(merged_configs_file, 'r') as f:
    content = json.loads(f.read())
    basic_auth_settings: Mapping[str, str] = content['basic_auth']
    current_workers_count: int = content['workers_pool']['max_workers_count']


app = Flask(__name__)
auth = HTTPBasicAuth()


def start_task_executor(workers_count: int) -> None:
    os.system(f'WORKERS_COUNT={workers_count} make -C {root_dir} executor-run > ./logs 2>&1 &')


@auth.verify_password
def verify_password(username, password):
    if (
            username == basic_auth_settings.get('user') and
            password == basic_auth_settings.get('password')
    ):
        return username


@app.route('/alert')
@auth.login_required
def hello_world():
    start_task_executor(42)
    return 'ok'
