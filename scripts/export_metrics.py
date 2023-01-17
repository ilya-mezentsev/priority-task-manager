"""
Скрипт для скачивания метрик из Prometheus по API,
чтобы склеить значения в файлы для разных ролей для дальнейшей обработки.
"""
import os
import sys
from collections import defaultdict

import requests


assert len(sys.argv) > 1, f'Missed arguments. Please pass out dir'

url = 'http://localhost:9090/api/v1/query?query=extracted_tasks_waiting_time[12h]'
try:
    r = requests.get(
        url=url,
        timeout=5,
    )
    response = r.json()
except Exception as e:  # noqa
    print(f'Got exception: {e!r}')
    raise RuntimeError('Unable to make request to prometheus') from e

assert response.get('status') == 'success', f'Got unsuccessful status: {response.get("status")}'
assert response.get('data') and response['data'].get('result'), f'No result in response: {response}'

role_to_metric_values: defaultdict[str, list[str]] = defaultdict(list)
for res in response['data']['result']:
    assert res.get('metric') and res['metric'].get('role'), f'Invalid metric structure: {res.get("metric")}'
    assert res.get('values'), f'Invalid metric values structure: {res.get("values")}'

    for val in res['values']:
        role_to_metric_values[res['metric']['role']].append(val[1])

for key, val in role_to_metric_values.items():
    file_path = os.path.join(sys.argv[1], key)
    with open(file_path, 'w') as f:
        f.write(','.join(val))
