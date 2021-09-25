# Задача

Реализовать микросервис на Golang, который позволяет понять во сколько потоков можно парсить определенный сайт без ошибок (нагрузочное тестирование).

Сервису на вход приходит поисковая строка, например “playstation купить”. Из поисковой выдачи Яндекса (готовый парсер тут) получаем список урлов. Далее для каждого урла нужно провести небольшой бенчмарк -- сколько параллельных запросов с одного IP этот урл выдерживает без ошибок. Максимальное время ответа - до 3 секунд. Ответом на исходный GET должна быть мапа “хост” => “рекомендуемое количество одновременных потоков”. 

Обязательные моменты, которые должны быть реализованы:

* сервис должен быть обернут в docker
* взаимодействие через один эндпоинт GET /sites?search=foobar
* с непрогретым кэшом сервис должен отвечать не дольше 30 сек.
* настройка параметров через конфиг

# Решение

Два http сервиса

Основной, service, ловит /sites?search= запросы, дергает внутренний сервис для выкачивания yandex-страницы, парсинга и дергает по каждому URL checker. Содержит внутренний кэш. Если после парсинга yandex-страницы по каким-то URL нет статистики, то дергает по ним checker. По той статистике что есть checker не дергает. Работает с ограничением ответа в 3 секунды, если за 3 секунды от checker не пришел ответ, то отправляет клиенту все что смог собрать.

Checker это http сервис, который держит несколько открытых http-клиентов, получает url для тестирования и запускает запросы через своих клиентов одновременно. Содержит внутренний кэш. Получив url от service выполняет нагрузочное тестирование и отдает результат в service. Если service уже отвалился по таймауту все равно продолжает работать, т.е. перед нагрузочным тестированием стоит очередь из задач.

Оба сервиса содержат внутренний кэш, возможно надо поменять на внешний, например на Redis чтобы checker клал результаты для service асинхронно.

Конфигурация:

Service:

* SERVICE_ADDR - на каком порту запускать
* CHECKER_URL - адрес checker
* YANDEX_FETCHERS - количество одновременно работающих Yandex-fetchers

Checker:

* CHECKER_ADDR - на каком порту пускать
* CHECKER_FETCHERS - сколько одновременных клиентов нагружать на страницы
    

![в процессе прогрева, timeout 3 sec](images/s1.png)


