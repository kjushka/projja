projja is a telegram task tracker system

### Сборка и запуск
1. Установить необходимое ПО
    1. Docker, Docker-Compose ([Docker](https://docs.docker.com/engine/install/ubuntu/), [Docker-Compose](https://docs.docker.com/compose/install/))
    2. Добавить Docker в группу ([Docker](https://itsecforu.ru/2018/04/12/как-использовать-docker-без-sudo-на-ubuntu/) 1-й вариант)
    3. Установить make командой sudo apt-get install make
2. Перейти в директорию проекта
3. Для сборки сервисов используется команда make install   
4. Для поднятия сервера используется команда make run
5. Для остановки работы сервера - make stop
6. Для удаления контейнеров - make down   
7. Логи можно посмотреть с помощью команды make logs
8. Сервер api запускается по адресу localhost:8080
9. Сервер exec запускается по адресу localhost:8090
