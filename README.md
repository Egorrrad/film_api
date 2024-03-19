# REST API для создания фильмотеки на Go


### Запуск приложения:

```
make build && make run
```

Если запускается впервые, необходимо применить миграции к базе данных:

```
make migrate
```

Также нужно подключиться к базе данных:
```
psql -h localhost -p 5436 -d film_api -U api_tester -W
```

Ввести пароль:
```
testing
```

И создать API ключ "root":
```
INSERT INTO public.users (id, role, api_key) VALUES (1, 'admin', 'root');
```

Даллее при каждом запросе к API нужно указывать ключ "root"