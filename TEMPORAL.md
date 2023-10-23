# [Determinism, Wokrflow + Temporal](https://www.youtube.com/watch?v=YfWu5swj-Gg)

## Виды окрестрации

- BPM
  - Camunda
  - Zeebe
- DSL
  - Conductor
  - Yandex.Procaas
- Code
  - Cadence
  - Temporal

**Temporal - это масштабируемая распределённая платформа оркестрации рабочих процессов**

## Workflow As Code

Определение рабочего процесса:

- Отказоустойчивая программа
- Выполняющая задачи
- Реагирующая на внешние события, включая
- Таймеры и таймауты

![](./assets/temporal/workflow-as-code.png)
![](./assets/temporal/workflow-cut.png)

## Stateful Execution Model

- **Workflow Function** - код, написанный на Go, который описывает работу Workflow с использованием Temporal SDK
- **Workflow Execution** - выполнение Workflow на стороне Temporal, который обеспечивает надёжное выполнение Workflow Function

![](./assets/temporal/stateful-execution-model.png)

**Replay** - это метод, с помощью которого выполение рабочего процесса возобновляет выполение. Во время повтора сгенерированные команды проверяются на соответсвие существующей истории событий.

![](./assets/temporal/event-log.png)

## Inversion Of Execution

![](./assets/temporal/inversion-of-execution.png)

## Order Management System

![](./assets/temporal/order-management-system.png)

"Elastic - альтернатив нет; позволяет индексировать workflow и делать выборки из него, что очень полезно для нашей админки, для саппорта".

## Use the State, Luke

![](./assets/temporal/two-db.png)

![](./assets/temporal/use-state.png)

![](./assets/temporal/true-way.png)

## Polling Activity

![](./assets/temporal/polling-activity.png)

Разработчики Temporal не рекомендуют поллить чаще, чем с периодом в 1 минуту.

![](./assets/temporal/high-freq-activity.png)

![](./assets/temporal/low-freq-activity.png)

## Signal is not Promise

![](./assets/temporal/signal-is-not-promise.png)

Temporal гарантирует отправку сигнала. Но какие изменения произойдут в state - синхронно узнать невозможно.

![](./assets/temporal/signal-is-not-promise-code.png)

## Workflow Not-Determinism

![](./assets/temporal/workflow-not-determinism.png)

### Как недопустить?

- Версионировать
- Использовать [workflowcheck](https://pkg.go.dev/go.temporal.io/sdk/contrib/tools/workflowcheck)
- Использовать replay-тесты
- Использовать Side Effects для недетерминированной логики
- Использовать _правильные конструкции_ GoLang

![](./assets/temporal/workflow-not-determinism-version.png)

![](./assets/temporal/workflow-not-determinism-replay-test.png)

![](./assets/temporal/workflow-not-determinism-true-go.png)

### Как локализовать?

- Выгрузить проблемный workfow, и прогнать его через replay-тест с помощью дебагера
- Изучать стек-трейсы в UI Temporal

![](./assets/temporal/workflow-not-determinism-fuckup.png)

## Fail State

Если workflow нельзя завершить ошибкой, а activity не может выполниться успешно, то workflow переходит в Fail State.

### Как такое возможно?

- Проблема во внешней системе
- Ошибка в исходных данных
- Невалидные данные из внешней системы

### А как решать?

- Бесконечные попытки отправки
- Возможность отредактировать исходные данные
- Ручное повторение активити
- Эскалация в саппорт
- Не сразу завершать workflow

DataDog помогает мониторить

## Синхронный workflow

![](./assets/temporal/sync-workflow.png)

компромисс между скоростью и возможностью компенсаций

## Что же такое Temporal?

- Гарантия исполнения
- Огромная нагрузка на DB
- Парадигма программирования
- Фреймворк/SDK
- Транспорт (сигналы вместо gRPC/HTTP ручек)
- Не про realtime
- Комьюнити
