package testCases

import (
	"go.uber.org/zap"
	"log"
	"log/slog"
)

func testLowerCaseStart() {
	log.Printf("Start smth")           // want "the log message must start with lowercase letter"
	slog.Error("Database fell asleep") // want "the log message must start with lowercase letter"
	log.Print("Error occured")         // want "the log message must start with lowercase letter"

	log.Println("everything is fine")
	slog.Error("nothing happened")
}

func testEnglishOnly() {
	log.Printf("полёт успешен")    // want "the log message must be in english"
	slog.Error("ошибка полёта гг") // want "the log message must be in english"
	log.Println("улёт")            // want "the log message must be in english"

	log.Printf("everything is fine")
	slog.Error("nothing happend")
}

func testSpecialSymbols() {
	log.Printf("started server!!!") // want "the log message must not contain any special symbols"
	slog.Error("smth failed...")    // want "the log message must not contain any special symbols"
	log.Println("what rocket 🚀")    // want "the log message must not contain any special symbols"

	log.Printf("everything is fine")
	slog.Error("nothing happend")
}

func testSensitiveData() {
	password := "1234"
	apiKey := "ke"
	token := "tok"
	authToken := "auth token"

	log.Println("password" + password) // want "the log message must not contain any sensitive data: password"
	log.Println("api key " + apiKey)   // want "the log message must not contain any sensitive data: apiKey"
	slog.Error("token " + token)       // want "the log message must not contain any sensitive data: token"
	slog.Error("authh" + authToken)    // want "the log message must not contain any sensitive data: authToken"

	name := "Ian"
	id := "1234"

	log.Println("user id" + id)
	slog.Error("user name" + name)
}

func testExtra() {
	log.Printf("")
	log.Println("   ")
	log.Println("1111")

	password := "1234"
	apiKey := "1111"
	slog.Error("password " + password + "apikey " + apiKey) // want "the log message must not contain any sensitive data: password, apiKey"
	slog.Info("user authenticated", "user", "john", "role", "admin")

	logger, _ := zap.NewProduction()
	logger.Info("Staaart")  // want "the log message must start with lowercase letter"
	zap.L().Info("Staaart") // want "the log message must start with lowercase letter"

	logger.Info("ну уж нет") // want "the log message must be in english"
	logger.Error("wtf!!!")   // want "the log message must not contain any special symbols"

	logger.Error("my password is hehehe" + password) // want "the log message must not contain any sensitive data: password"

	logger.Error("be happy dont worry")
}
