import random
import uuid
from abc import abstractmethod

from locust import (
    HttpUser,
    task,
    between,
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
    wait_time = between(0.1, 0.2)

    @task
    def do_task(self) -> None:
        self.do_request()

    def get_account_id(self) -> str:
        return ROLE_1_ID


class Role2User(HttpUser, UserMixin):
    wait_time = between(0.1, 0.2)

    @task
    def do_task(self) -> None:
        self.do_request()

    def get_account_id(self) -> str:
        return ROLE_2_ID


class Role3User(HttpUser, UserMixin):
    wait_time = between(0.1, 0.2)

    @task
    def do_task(self) -> None:
        self.do_request()

    def get_account_id(self) -> str:
        return ROLE_3_ID
