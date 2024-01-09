# Request For Proposal (RFP) - FLUTTER AWESOME

Идея социально-значимого мега-восстребованного пет-проекта для вашего портфолио

Распарсить flutterawesome.com на go-colly (см. пример (https://github.com/comerc/try-colly)) и сделать свой реестр open source проектов на Flutter для поиска и аналитики с различными срезами. Ещё можно сделать подписку на ежемесячную рассылку дайджестов. Ещё можно сделать бот в Телеге. Ещё можно сделать мониторинг: какие проекты в активной стадии, а какие уже дохлые.

Можно применить: ClickHouse, NATS/RabbitMQ/Kafka, Hasura/GraphQL, GoHugo/NextJS/Flutter!

Для вдохновения: https://fluttergems.dev/

---

Звёздочки можно подключить:

- https://img.shields.io/github/stars/user/repo
- https://api.github.com/repos/user/repo

Но во втором случае данные будут на моей стороне (что важно для поиска)

---

Посты:

`https://flutterawesome.com/sitemap-posts.xml`

Название:

`<meta property="og:title" content="">`

Картинка обложки:

`<meta property="og:image" content="">`

- `<meta property="og:image:width" content="719">`
- `<meta property="og:image:height" content="522">`

Теги:

`<meta property="article:tag" content="">`...

Ссылка на GitHub:

`<a class="github-view" href="https://github.com/user/repo?ref=flutterawesome.com" target="_blank" rel="noopener">View Github</a>`

---

UI - SPA на vercel или telegram-бот:

- поле поиска

  - сортировка по названию
  - сортировка по звёздочкам
  - сортировка по "search relevance"

- результат поиска

- облако тегов
- пост:
  - сколько звёздочек
  - сколько форков
  - отношение открытые/закрытые issues
  - отношение открытые/закрытые pull requests
  - когда стартовал
  - когда последний commit
  - состояние (активная разработка, на поддержке, брошен)
  - сколько контрибьютеров и их вклад
  - какая версия внутри pubspec.yaml > environment > sdk

---

Crowler:

- cli-утилита на go
- на входе:
  - /sitemap-tags.xml по Last Modified
  - /sitemap-posts.xml по Last Modified
  - https://api.github.com/repos/user/repo
- на выходе
  - сводный json
  - графики Contributors
- перезапускать регулярно для актуализации доступности репок и их параметров
