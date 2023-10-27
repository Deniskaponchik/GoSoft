# [Main]([https://github.com/Deniskaponchik/](https://github.com/Deniskaponchik/GoSoft/tree/main/Unifi))
| | |
| - | - |
| Code Structure | github.com/golang-standards/project-layout |
| ORM | нет |
| Aps | тикеты создаются со второй попытки (через 24 минуты) |
| Clients | При обработке каждого клиента подключение к мапе точек для получения имени точки |
| Anomalies | Ежедневная обработка. При обработке каждой аномалии подключение к мапе клиентов для получения имени точки |
| Poly | с перезагрузками |


# [Dev1]([https://github.com/Deniskaponchik/](https://github.com/Deniskaponchik/GoSoft/tree/dev1/Unifi))
| | |
| - | - |
| Code Structure | github.com/golang-standards/project-layout |
| ORM | github.com/jinzhu/gorm |
| Aps | тикеты создаются со второй попытки (через 24 минуты) |
| Clients | При обработке каждого клиента подключение к мапе точек для получения имени точки НЕ производится. По умолчанию доступен только мак точки |
| Anomalies | Ежедневная обработка. При обработке каждой аномалии подключение к мапе клиентов для получения имени точки НЕ производится. Получение имени точки только на финальной стадии при создании заявки |
| Poly | с перезагрузками |

# [Dev2]([https://github.com/Deniskaponchik/](https://github.com/Deniskaponchik/GoSoft/tree/dev2/Unifi))
- AP. тикеты со второй попытки
- Anomalies. Ежедневная обработка
- Poly. с перезагрузками
- 


