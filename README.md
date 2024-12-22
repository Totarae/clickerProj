Задача 1. Счетчик кликов.
Есть набор баннеров (от 10 до 100). У каждого есть ИД и название (id, name)

Нужно сделать сервис, который будет считать клики и собирать их в поминутную статистику (timestamp, bannerID, count)



Нужно сделать АПИ с двумя методами:

1. /counter/<bannerID> (GET)

Должен посчитать +1 клик по баннеру с заданным ИД



2. /stats/<bannerID> (POST)

Должен выдать статистику показов по баннеру за указанный промежуток времени (tsFrom, tsTo)



Язык: golang

СУБД:  pg

Примеры запросов:
````
curl --location --request POST "http://localhost:8080/stats/17" \
--header "Content-Type: application/json" \
--data-raw "{
\"tsFrom\": \"2024-12-22T18:36:00Z\",
\"tsTo\": \"2024-12-22T19:50:00Z\"
}"
````
````
curl --location --request GET "http://localhost:8080/counter/10"
````