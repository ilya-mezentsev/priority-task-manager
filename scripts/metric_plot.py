"""
Скрипт для рисования графиков.
В частности, для удобного сравнения графиков
"""
import os
import sys

import matplotlib.pyplot as plt


assert len(sys.argv) == 3, f'Script can process exactly two arguments; first - is metric before, second - is after'

before_title = os.environ.get('BEFORE_TITLE') or 'Before'
after_title = os.environ.get('AFTER_TITLE') or 'After'

labels = [
    before_title,
    after_title,
]
time_points_list: list[list[float]] = []
for file in sys.argv[1:]:
    with open(file, 'r') as f:
        time_points_list.append(list(map(float, f.read().split(','))))

min_points_count = len(min(*time_points_list, key=lambda points: len(points)))
x_points = [i for i in range(min_points_count)]
for time_points in time_points_list:
    plt.plot(
        x_points,
        list(map(float, time_points[:min_points_count])),
        label=labels.pop(0),
    )

plt.ylabel('Время ожидания задачи взятия в работу, [с]')
plt.xlabel('Время, [с]')

plt.grid()
plt.title(os.environ.get('TITLE') or f'{before_title} vs {after_title}')
plt.legend()

plt.show()
