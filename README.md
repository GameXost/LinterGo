# Linter for Golang

## Базовые правила
1. Лог-сообщения должны начинаться со строчной буквы
2. Лог-сообщения должны быть только на английском языке 
3. Лог-сообщения не должны содержать спецсимволы или эмодзи 
4. Лог-сообщения не должны содержать потенциально чувствительные данные

## Работает с:
- log
- log/slog
- go.uber.org/zap

## Установка и запуск
- Необходимо установить golangci-lint-custom
```
choco install golangci-lint
или же
go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.10.1
```
- Создать .yml файлы конфигурации в корне проекта

##### .custom-gcl.yml

```yml
version: v2.9.0

plugins:
  - module: github.com/GameXost/LinterGo
    version: v1.0.0
```

##### .golangci.yml
```yml
version: "2"
linters:
  enable:
    - loglinter


  settings:
    custom:
      loglinter:
        type: "module"
        description: "Checks logs to be proper"
        settings:
          extra-ban-words: ["private", "IP"]
          disable-flags: [] #"low_first_letter", "english-only", "special-symbols", "sensitive-words"
```

- Собрать бинарник
```
golangci-lint custom
```
- Запустить линтер на проекте или отдельной папке
```
./custom-gcl run ./...
./custom-gcl run путь до папки
```
## Прогон тестов
` go test ./...`

### Сделал бонуски 3 и 4
- Можно указывать свои паттерны - бан ворды для 3 правила
```
Для добавления бан вордов необходимо прописать в файле golangci.yml в поле settings
По примеру: `extra-ban-words: ["private", "IP"]`
```
- Использовал github/workflows/go.yml для автоматической сборки и тестирования
```
Для отлкючения какого-либо правила необходимо записать его в файле golangci.yml в поле settings
По примеру: `disable-flags: ["low_first_letter", "english-only"]`
названия строго в соответсвии со следующими:
1 правило - "low_first_letter"
2 правило - "english-only"
3 правило - "special-symbols"
4 правило - "sensitive-words"
```

 

## Примеры
```Go
// logger, _ := zap.NewProduction()  для работы запа на тестах сделал пакет, чтоб не было проблем с импортом(у меня были :)))

log.Printf("Start smth")           //"the log message must start with lowercase letter"
slog.Error("Database fell asleep") //"the log message must start with lowercase letter"
log.Print("Error occured")         //"the log message must start with lowercase letter"
logger.Info("Staaart")             // "the log message must start with lowercase letter"
zap.L().Info("Staaart")            // "the log message must start with lowercase letter"


log.Printf("полёт успешен")    // "the log message must be in english"
slog.Error("ошибка полёта гг") // "the log message must be in english"
log.Println("улёт")            // "the log message must be in english"
logger.Info("ну уж нет")        // "the log message must be in english"


log.Printf("started server!!!") // "the log message must not contain any special symbols"
slog.Error("smth failed...")    // "the log message must not contain any special symbols"
log.Println("what rocket 🚀")    // "the log message must not contain any special symbols"
logger.Error("wtf!!!")          // "the log message must not contain any special symbols"


log.Println("password" + password) // "the log message must not contain any sensitive data: password"
log.Println("api key " + apiKey)   // "the log message must not contain any sensitive data: apiKey"
slog.Error("token " + token)       // "the log message must not contain any sensitive data: token"
slog.Error("authh" + authToken)    // "the log message must not contain any sensitive data: authToken"
slog.Error("password " + password + "apikey " + apiKey) // "the log message must not contain any sensitive data: password, apiKey"
logger.Error("my password is hehehe" + password)        // "the log message must not contain any sensitive data: password"

log.Println("user id" + id) // fine
slog.Error("user name" + name) // fine
log.Printf("")  // fine
log.Println("   ") // fine 
log.Println("1111") // fine

slog.Info("user authenticated", "user", "john", "role", "admin") // fine
logger.Error("be happy dont worry") // fine
```

### Проверил на своем проектике
```
cmd\main.go:62:3: the log message must start with lowercase letter (loglinter)
                log.Println("Cache loaded")
internal\service\service_bruh.go:45:3: the log message must not contain any special symbols (loglinter)
                log.Printf("inalid order data...")
                ^                                        
```

## ПЫ СЫ:
- Решил, что цифры не спец символ, поэтому они проходят проверки
- Любые знаки препинания решил считать спец символами
- log тоже поддерживается, при желании можно добавить в пути других логгеров в availableLoggers из run.go и они будут работать
- Проверка на логер осуществляется через проверку пакета метода(через его путь)
- Проверка на логгер начинается с проверки самого метода, я вынес основные в условия и первичная проверка проходит оп ним, далее уже следует проверка по пакету метода
