"""
Скрипт для мерджа json-файлов.
Использовать так:

python3 json_merge.py out-dir file1.json file2.json ... fileN.json

результат:
Создаем/перезаписываем файл out-dir/main.json. Содержание такое:
{
  file1: { ... },
  file2: { ... },
  ...
  fileN: { ... }
}

P.S. Нюанс, конечно, в том, что файлы с одинаковым названием адекватно никак не обрабатываются
"""

import json
import os.path
import sys
from typing import (
    Sequence,
    Mapping,
    Any,
    Dict,
)


MainContent = Dict[str, Mapping[str, Any]]


def _process(
        main_content: MainContent,
        file_path: str,
) -> None:

    assert os.path.exists(file_path), f'File {file_path} does not exist!'

    with open(file_path, 'r') as f:
        base_name = os.path.basename(f.name)
        section_name = os.path.splitext(base_name)[0]

        assert section_name not in main_content, f'Filename {base_name} met twice'
        main_content[section_name] = json.loads(f.read())


def main(
        out_dir: str,
        file_paths: Sequence[str],
) -> None:

    main_content: MainContent = {}
    main_file = os.path.join(out_dir, 'main.json')

    for file_path in file_paths:
        try:
            _process(
                main_content=main_content,
                file_path=file_path,
            )
        except Exception as e:
            raise RuntimeError(f'Unable to process file {file_path}') from e

    with open(main_file, 'w') as f:
        f.write(json.dumps(main_content, indent=4))


if __name__ == '__main__':
    assert len(sys.argv) > 3, 'Invalid number of arguments (expected more than 3)'

    main(
        out_dir=sys.argv[1],
        file_paths=sys.argv[2:],
    )
