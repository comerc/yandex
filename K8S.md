https://youtu.be/MHn-taXfQ8o?t=3748

Определения:

- CFS (Completely Fair Scheduler)
- GMP
- CPU MEM Limits
- Containers
- Pods
- Nodes
- Requests

Рекомендации:

- Requests: всегда определяйте CPU MEM requests
- CPU Limits: использовать в связке с GOMAXPROCS (иначе тротлятся)
- MEM Limits: выставлять равным requests (иначе убиваются)
- Мониторить CPU Throttling ваших контейнеров
- Всегда контролируйте значение GOMAXPROCS
- Оптимальный GOMAXPROCS всегда зависит от вашего контекста
- Выставляйте целые числа в CPU Limit
- Пробуйте разные значения, пока не найдёте оптимальную комбинацию:
  - Server Latency (50, 90, 99, 99.9 персентили)
  - CPU Usage
  - GC Duration
