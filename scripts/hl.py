import random
import uuid
from abc import abstractmethod
from typing import (
    Sequence,
    Mapping,
    Optional,
    Tuple,
)

from locust import (
    HttpUser,
    task,
    constant,
    LoadTestShape,
)
from locust.clients import HttpSession

ROLE_1_ID = '86d3f3a95c324c9479bd8986968f4327'
ROLE_2_ID = '11c9beec53034beb3a6687891c9e248a'
ROLE_3_ID = 'eafae37058a254f5dfdfe22ede8cca1f'
TASK_TYPES = (
    'foo',
    'bar',
    'baz',
)


class UserMixin:
    client: HttpSession

    @abstractmethod
    def get_account_id(self) -> str:
        raise NotImplementedError()

    def do_request(self) -> None:
        self.client.post(
            '/external/task',
            json=dict(
                type=random.choice(TASK_TYPES),
                data=dict(
                    exec_time=random.randint(1, 3),
                    some_id=str(uuid.uuid4()),
                ),
            ),
            headers={
                'X-Account-Hash': self.get_account_id(),
            },
        )


class Role1User(HttpUser, UserMixin):
    wait_time = constant(1)

    @task
    def do_task(self) -> None:
        self.do_request()

    def get_account_id(self) -> str:
        return ROLE_1_ID


class Role2User(HttpUser, UserMixin):
    wait_time = constant(1)

    @task
    def do_task(self) -> None:
        self.do_request()

    def get_account_id(self) -> str:
        return ROLE_2_ID


class Role3User(HttpUser, UserMixin):
    wait_time = constant(1)

    @task
    def do_task(self) -> None:
        self.do_request()

    def get_account_id(self) -> str:
        return ROLE_3_ID


class StagesShape(LoadTestShape):
    """
    A simple load test shape class that has different user and spawn_rate at
    different stages.
    Keyword arguments:
        stages -- A list of dicts, each representing a stage with the following keys:
            duration -- When this many seconds pass the test is advanced to the next stage
            users -- Total user count
            spawn_rate -- Number of users to start/stop per second
        stop_at_end -- Can be set to stop once all stages have run.
    """

    stages: Sequence[Mapping[str, int]] = [
        {'duration': 30, 'users': 30, 'spawn_rate': 10},
        {'duration': 210, 'users': 90, 'spawn_rate': 10},
        {'duration': 450, 'users': 180, 'spawn_rate': 10},
        {'duration': 630, 'users': 90, 'spawn_rate': 10},
        {'duration': 660, 'users': 30, 'spawn_rate': 5},
    ]
    stop_at_end = True

    def tick(self) -> Optional[Tuple[int, int]]:
        run_time = self.get_run_time()

        for stage in self.stages:
            if run_time < stage['duration']:
                tick_data = (stage['users'], stage['spawn_rate'])
                return tick_data

        return None
