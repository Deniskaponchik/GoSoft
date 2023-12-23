# Описание проекта
Приложение подключается по API к Wi-Fi контроллеру на базу Unifi. <br/>
Обрабатывает полученную информацию и создаёт обращения в корпоративной тикет-системе через SOAP, 
упорядочивает данные в БД MySQL 
или по запросу отображает клиентов и точки доступа через web-интерфейс благодаря GIN-фреймворку. 

# Внешний вид
<div align="left">
  <img src="https://github.com/Deniskaponchik/GoSoft/blob/main/web/png/client.png" width="600"/>
</div>
<div align="left">
  <img src="https://github.com/Deniskaponchik/GoSoft/blob/main/web/png/swager.PNG" width=600"/>
</div>

# [Main]([https://github.com/Deniskaponchik/](https://github.com/Deniskaponchik/GoSoft/tree/main/Unifi))
|                |                                                                                                                                                                                                           |
|----------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| version        | 3.5                                                                                                                                                                                                       |
| Code Structure | https://www.github.com/golang-standards/project-layout <br/> https://www.youtube.com/watch?v=V6lQG6d5LgU                                                                                                  |
| ORM            | https://www.github.com/jinzhu/gorm                                                                                                                                                                        |
| Web server     | https://www.github.com/gin-gonic/gin                                                                                                                                                                      |
| Arguments      | -mode : TEST or PROD <br/> -time : server time zone <br/> -httpUrl : server url for created ticket <br/> -db : db name                                                                                    |
| Unifi          | Оба контроллера обрабатываются в одном приложении                                                                                                                                                         |
| Aps            | тикеты создаются со второй попытки (через 24 минуты)                                                                                                                                                      |
| Clients        | При обработке каждого клиента информция заносится в 2 мапы: по маку и по имени машины.                                                                                                                    |
| Anomalies      | ЧАС. Информация заносится клиенту в МАССИВ аномалий. <br/> СУТКИ. Пробегаемся по МАПЕ клиентов, а не делаем запрос к БД. <br/> МЕСЯЦ. Дропаем массив аномалий у клиентов и загружаем снова из БД за 30 дн. |
| Poly           | с перезагрузками, без отдельного web-интерфейса                                                                                                                                                           |


# [Dev3]([https://github.com/Deniskaponchik/](https://github.com/Deniskaponchik/GoSoft/tree/dev3/Unifi))
|               |                                                                                                                                                                                                           |
|---------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| version       | 3.5                                                                                                                                                                                                       |
| Code Structure | https://www.github.com/golang-standards/project-layout <br/> https://www.youtube.com/watch?v=V6lQG6d5LgU                                                                                                  |
| ORM           | https://www.github.com/jinzhu/gorm                                                                                                                                                                        |
| Web server    | https://www.github.com/gin-gonic/gin                                                                                                                                                                      |
| Log           |                                                                                                                                                                                                           |
| Unifi         | Оба контроллера обрабатываются в одном приложении                                                                                                                                                         |
| Aps           | тикеты создаются со второй попытки (через 24 минуты)                                                                                                                                                      |
| Clients       | При обработке каждого клиента информция заносится в 2 мапы: по маку и по имени машины.                                                                                                                    |
| Anomalies     | ЧАС. Информация заносится клиенту в МАССИВ аномалий. <br/> СУТКИ. Пробегаемся по МАПЕ клиентов, а не делаем запрос к БД. <br/> МЕСЯЦ. Дропаем массив аномалий у клиентов и загружаем снова из БД за 30 дн. |
| Poly          | с перезагрузками, без отдельного web-интерфейса                                                                                                                                                           |


# [Dev2]([https://github.com/Deniskaponchik/](https://github.com/Deniskaponchik/GoSoft/tree/dev2/Unifi))
|               |                                                                                                                                                                                                |
|---------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| version       | 2.0                                                                                                                                                                                            |
| Code Structure | github.com/golang-standards/project-layout                                                                                                                                                     |
| ORM           | Нет                                                                                                                                                                                            |
| Web server    | github.com/gin-gonic/gin                                                                                                                                                                       |
| Unifi         | Каждый контроллер запускается в отдельном терминале                                                                                                                                            |
| Aps           | тикеты создаются со второй попытки (через 24 минуты)                                                                                                                                           |
| Clients       | При обработке каждого клиента информция заносится в 2 мапы: по маку и по имени машины.                                                                                                         |
| Anomalies     | ЧАС. Информация заносится клиенту в МАССИВ аномалий. <br/>СУТКИ. Пробегаемся по МАПЕ клиентов, а не делаем запрос к БД. <br/>МЕСЯЦ. Дропаем массив аномалий у клиентов и загружаем снова из БД за 30 дн. |
| Poly          | с перезагрузками                                                                                                                                                                               |

# [Dev1]([https://github.com/Deniskaponchik/](https://github.com/Deniskaponchik/GoSoft/tree/dev1/Unifi))
| | |
| - | - |
| version | 1.3 |
| Code Structure | github.com/golang-standards/project-layout |
| ORM | нет |
| Aps | тикеты создаются со второй попытки (через 24 минуты) |
| Clients | При обработке каждого клиента подключение к мапе точек для получения ИМЕНИ точки НЕ производится. <br/>По умолчанию доступен только мак точки |
| Anomalies | При обработке каждой аномалии раз в час производится подключение к мапе клиентов для получения МАКА точки, потом подключение к мапе точек для получения и занесение имени точки в БД. |
| Poly | с перезагрузками |

# [Dev]([https://github.com/Deniskaponchik/](https://github.com/Deniskaponchik/GoSoft/tree/dev/Unifi))
| | |
| - | - |
| version | 0.5 |
| Code Structure | старый лапша-код, всё в одном main файле практически |
| Aps | тикеты создаются со второй попытки (через 24 минуты) |
| Clients |  |
| Anomalies | Ежедневная обработка |
| Poly | с перезагрузками |


