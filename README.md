# Effective-Mobile-TestCase

0. Рабочая ветка - 0

1. Перед началом работы с сервисом необходимо создать файл .env, пример файла лежит в репозитории, можно просто переименовать его с .env.example в .env. Сервисы Endpoint и Controller поднимаются на разных портах, в связи с чем общение происходит через gRPC (некое подобие микросервисов);
   
2. Для запуска сервиса можно воспользоваться:
   1.  стандартной командой docker: docker compose up -d --build
   2.  раннер task: task start

3. Документация сервиса:
   1. GET /service/subscription - возврат списка всех сервисов;
   2. GET /service/subscription/totalprice/start/{date_start}/stop/{date_stop}?user_uuid={some_string}&service_name={some_string} - возврат суммы подписок за выбранный период с возможностью фильтрации по user_uuid и service_name. Параметры {date_start} и {date_stop} необходимо вводить в формате ММ.ГГГГ;
   3. POST /service/subscription - запрос по добавлению новой подписки в сервис. Необходимо передать параметры в формате json в тело запроса:  {
    "service_name": "second service",
    "price": 4321,
    "user_uuid": "unique_uuid",
    "start_date": "06.2025",
    "stop_date": "10.2025"
   }
   4. PUT /service/{service_name}/user/{user_uuid}/subscription - запрос изменения параметров подписки, новые параметры необходимо передать параметры в формате json в тело запроса: {
    "price": 4321,
    "start_date": "06.2025",
    "stop_date": "10.2025"
   }
   5. DELETE /service/{service_name}/user/{user_uuid}/subscription - запрос на удаление выбранной подписке по user_uuid и service_nameю.