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

from flask import Flask, request
from flask_httpauth import HTTPBasicAuth


queue_increase_alert_names = frozenset({
    'QueueSmallIncrease',
    'QueueMediumIncrease',
    'QueueDamnIncrease',
})


min_workers_count = 30
max_workers_count = 240
root_dir = os.environ.get('ROOT_DIR')
merged_configs_file = os.environ.get('MERGED_CONFIGS_FILE')
assert merged_configs_file and root_dir, 'Unable to work without MERGED_CONFIGS_FILE and ROOT_DIR in env'

with open(merged_configs_file, 'r') as f:
    content = json.loads(f.read())
    basic_auth_settings: Mapping[str, str] = content['basic_auth']
    current_workers_count: int = content['workers_pool']['max_workers_count']


app = Flask(__name__)
auth = HTTPBasicAuth()


def start_task_executor(workers_count: int) -> None:
    os.system(f'WORKERS_COUNT={workers_count} make -C {root_dir} executor-run > {root_dir}/scripts/logs 2>&1 &')


@auth.verify_password
def verify_password(username, password):
    if (
            username == basic_auth_settings.get('user') and
            password == basic_auth_settings.get('password')
    ):
        return username


@app.route('/alert', methods=['POST'])
@auth.login_required
def hello_world():
    global current_workers_count

    alert = request.json['alerts'][0]
    alert_name = alert["labels"]["alertname"]
    if alert_name == 'QueueIncrease':
        should_increase = True
        coefficient = alert['annotations']['queue_increase_coefficient']
        coefficient = float(coefficient) if coefficient else 0
    elif alert_name == 'QueueDecrease':
        should_increase = False
        coefficient = alert['annotations']['queue_decrease_coefficient']
        coefficient = float(coefficient) if coefficient else 0
    else:
        raise RuntimeError(f'Unknown alert name: {alert_name}')

    # если получили "значимый" коэффициент
    if round(coefficient, 1):
        # todo тут по-хорошему надо использовать значение коэффициента
        current_workers_count = (
            current_workers_count * 2
            if should_increase
            else int(current_workers_count / 2)
        )

        current_workers_count = max(min_workers_count, current_workers_count)
        current_workers_count = min(max_workers_count, current_workers_count)

        app.logger.info(
            f'Got alert {alert["labels"]["alertname"]} with coefficient: {coefficient}; '
            f'trying to start {current_workers_count} workers'
        )

        start_task_executor(current_workers_count)

    return 'ok'


if __name__ == '__main__':
    app.run(
        host='0.0.0.0',
        port=5000,
        debug=True,
    )
