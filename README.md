# Telegram-бот для домовых чатов

Данный проект реализует серверную часть Telegram-бота, способного управлять домовым чатом. Основные особенности:

- Данный проект является бесплатным и распространяется свободно, вы можете самостоятельно создать и настроить бота,
  после чего быть уверенным, что ваши данные и данные ваших соседей не доступны владельцам сторонних ботов.
- Ваш бот в Telegram будет работать только в вашем домовом чате, его не смогут использовать другие Telegram-чаты.
- Вы можете использовать как современные облачные решения для хранения кодовой базы бота, так и обычные.
- Данный проект рассчитан на людей с техническим образованием. Если вы просто владелец домового чата, попробуйте найти специалистов, которые помогут вам разобраться с настройкой и установкой бота, или [обратитесь](mailto:stanislav.belichenko@gmail.com) к автору данного проекта.

## Содержание

- [Функции бота](#функции-бота)
- [Установка](#установка)
    - [Через GitLab Actions](#через-gitlab-actions)
    - [Вручную](#вручную)
- [Разработка](#разработка)
- [Конфигурация](#конфигурация)
    - [Для доставки в Yandex Cloud](#для-доставки-в-yandex-cloud)
    - [Настройки бота](#настройки-бота)
- [Лицензия](#лицензия)

## Функции бота

Команды бота:

- в личной переписке (доступны любому пользователю Telegram):
    - `/start` – начало общения с ботом;
    - `/help` – справка по командам в личной переписке и в домовом чате;
- в домовом чате (доступны всем его участникам):
    - `/report` – жалоба на нарушение правил другим участником;
    - `/keys` – шуточная команда, отвечающая пользователю когда выдадут ключи (можно ограничить определенной темой чата), актуально для чатов, где собрались пока еще дольщики;
- в административном чате:
    - для администрирования домового чата (доступны всем его участникам):
        - `/help_admin` – справка по командам в административном чате;
        - `/mute` – ограничить общение пользователя в домовом чате;
        - `/unmute` – снять ограничения на общение пользователя в домовом чате;
        - `/ban` – заблокировать (удалить) пользователя из домового чата;
        - `/unban` – разблокировать пользователя в домовом чате (без автоматического добавления в чат, но он сможет заново подать заявку на вступление);
    - сервисные команды (доступны администраторам):
        - `/set_commands` – установить команды бота (_только для отображения их в меню, это не влияет на их доступность!_);
        - `/delete_commands` – удалить команды бота.

Прочие функции:

- при отправке пользователем заявки на вступление в домовой чат (заявки должны быть включены в настройках чата) бот уведомит об этом в административном чате;
- при отправке жалобы на нарушение правил бот уведомит об этом в административном чате.

## Установка

[//]: # (TODO: Написать)

### Через GitLab Actions

[//]: # (TODO: Написать)

### Вручную

[//]: # (TODO: Написать)

## Разработка

## Конфигурация

### Для доставки в Yandex Cloud

Секреты:

- `YC_SA_JSON_CREDENTIALS` – JSON-файл с
  ключом ([инструкция для получения](https://yandex.cloud/ru/docs/iam/operations/iam-token/create-for-sa)).
- `YC_FOLDER_ID` – идентификатор заранее созданного каталога в YC.
- `YC_FUNCTION_NAME` – имя функции для бота в YC (может быть создана при первой доставке бота в YC).
- `YC_SERVICE_ACCOUT_ID` – идентификатор заранее созданного сервисного аккаунта в YC.
- `LOG_GROUP_ID` – идентификатор заранее созданной группы логов в YC.

Секреты для установки секретов бота из Lockbox:

- `YC_LOCKBOX_ID` – идентификатор секретов в Lockbox.
- `YC_LOCKBOX_VERSION` – версия секретов в Lockbox.

В указываемой версии должны быть созданы все [секреты бота](#настройки-бота).

Переменные:

- `FUNCTION_RUNTIME` – среда запуска бота. Используйте `golang121`, если не уверены.
- `FUNCTION_MEMORY` – лимит памяти для среды запуска бота. Начните с `128Mb`.
- `EXECUTION_TIMEOUT` – таймаут обработки запроса в секундах. Используйте `10`.
- `LOG_LEVEL` – минимальный уровень журналирования бота. Оставьте пустую строку, если не уверены.

### Настройки бота

Секреты (при работе бота в YC будут установлены автоматически из Lockbox):

- `TELEGRAM_BOT_TOKEN` – токен Telegram-бота, полученный от [@BotFather](https://t.me/BotFather).
- `LLM_API_TOKEN` – токен для работы с YandexGPT
  API ([инструкция для получения](https://yandex.cloud/ru/docs/iam/operations/authentication/manage-api-keys)).
- `LLM_FOLDER_ID` – идентификатор заранее созданного каталога в YC.

Переменные:

- `HOUSE_CHAT_ID` – идентификатор домового чата, управляемого ботом. Должен быть супергруппой и форумом (с темами). Бот должен быть добавлен в данный чат с полными правами.
- `HOME_THREAD_BOT` – идентификатор темы чата, в котором бот может использовать шуточные команды (например, "Оффтоп").
- `ADMINISTRATION_CHAT_ID` – идентификатор чата, где собраны администраторы домового чата. В этот чат будут приходить уведомления о заявках на вступление в группу, репорты на нарушение правил и прочее. Баны и удаление пользователей производятся именно в нем. Бот должен быть добавлен в данный чат с полными правами.

## Лицензия

[GNU AGPLv3](LICENSE).