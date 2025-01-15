## Установка Ansible через venv на macOS

0. Версиями Python лучше управлять через pyenv:
```
pyenv --help
```
1. Python обычно уже установлен, но если нет, установите его через Homebrew:
```
brew install python3
```
2. Создайте виртуальное окружение (команда такая же, как и на Linux):
```
python3 -m venv .venv
```
3. Активируйте окружение:
```
source .venv/bin/activate
```
4. Обновите pip:
```
python3 -m pip install --upgrade pip
```
5. Установите Ansible:
```
pip install ansible
```
6. Проверьте установку:
```
ansible --version
```
7. Для деактивации виртуального окружения:
```
deactivate
```
