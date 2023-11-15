# [Main]([https://github.com/Deniskaponchik/](https://github.com/Deniskaponchik/GoSoft/tree/main/Unifi))
| | |
| - | - |
| version | не доведён до работающего приложения|
| Code Structure | github.com/golang-standards/project-layout |
| ORM | нет |
| Aps | тикеты создаются со второй попытки (через 24 минуты) |
| Clients | При обработке каждого клиента подключение к мапе точек для получения имени точки |
| Anomalies | Ежедневная обработка. При обработке каждой аномалии подключение к мапе клиентов для получения имени точки |
| Poly | с перезагрузками |

# [Dev1]([https://github.com/Deniskaponchik/](https://github.com/Deniskaponchik/GoSoft/tree/dev1/Unifi))
| | |
| - | - |
| version | 2.3 |
| Code Structure | github.com/golang-standards/project-layout |
| ORM | github.com/jinzhu/gorm |
| Web server | github.com/gin-gonic/gin |
| Aps | тикеты создаются со второй попытки (через 24 минуты) |
| Clients | При обработке каждого клиента информция заносится в 2 мапы: по маку и по имени машины.  |
| Anomalies | ЧАС. Информация заносится клиенту в МАССИВ аномалий. СУТКИ. Пробегаемся по МАПЕ клиентов, а не делаем запрос к БД. МЕСЯЦ. Дропаем массив аномалий у клиентов и загружаем снова из БД за 30 дн.|
| Poly | с перезагрузками |

# [Dev1]([https://github.com/Deniskaponchik/](https://github.com/Deniskaponchik/GoSoft/tree/dev1/Unifi))
| | |
| - | - |
| version | 1.3 |
| Code Structure | github.com/golang-standards/project-layout |
| ORM | нет |
| Aps | тикеты создаются со второй попытки (через 24 минуты) |
| Clients | При обработке каждого клиента подключение к мапе точек для получения ИМЕНИ точки НЕ производится. По умолчанию доступен только мак точки |
| Anomalies | При обработке каждой аномалии раз в час производится подключение к мапе клиентов для получения МАКА точки, потом подключение к мапе точек для получения и занесение имени точки в БД. |
| Poly | с перезагрузками |

# [Dev]([https://github.com/Deniskaponchik/](https://github.com/Deniskaponchik/GoSoft/tree/dev2/Unifi))
| | |
| - | - |
| version | 0.5 |
| Code Structure | старый лапша-код, всё в одном main файле практически |
| Aps | тикеты создаются со второй попытки (через 24 минуты) |
| Clients |  |
| Anomalies | Ежедневная обработка |
| Poly | с перезагрузками |


