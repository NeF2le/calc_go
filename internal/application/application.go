package application

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"github.com/NeF2le/calc_go/pkg/calculation"
	"github.com/NeF2le/calc_go/pkg/logging"
)

type Config struct {
	Addr string
}

func ConfigFromEnv() *Config {
	config := new(Config)
	config.Addr = os.Getenv("PORT")
	if config.Addr == "" {
		config.Addr = "8080"
	}
	return config
}

type Application struct {
	config *Config
}

func NewApplication() *Application {
	return &Application{
		config: ConfigFromEnv(),
	}
}

func (a *Application) Run() error {
	for {
		log.Println("Input expression:")
		reader := bufio.NewReader(os.Stdin)
		text, err := reader.ReadString('\n')
		if err != nil {
			log.Println("Failed to read expression from console")
			continue
		}

		text = strings.TrimSpace(text)

		if text == "exit" {
			log.Println("Application was successfully closed")
			return nil
		}

		result, err := calculation.Calc(text)
		if err != nil {
			log.Println(text, "calculation failed with error:", err)
		} else {
			log.Println(text, "=", result)
		}
	}
}

type Request struct {
	Expression string `json:"expression"`
}

type Response struct {
	Result *float64 `json:"result,omitempty"`
	Error  string   `json:"error,omitempty"`
}

func CalcHanlder(w http.ResponseWriter, r *http.Request) {
	request := new(Request)
	defer r.Body.Close()

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{Error: "Internal server error"})
		return
	}

	if strings.TrimSpace(request.Expression) == "" {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(Response{Error: "Expression is not valid"})
		return
	}

	result, err := calculation.Calc(request.Expression)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(Response{Error: "Expression is not valid"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Response{Result: &result})
}

func (a *Application) RunServer() error {
	r := mux.NewRouter()

	logger := logging.SetupLogger()

	r.Use(logging.LoggingMiddleware(logger))

	r.HandleFunc("/api/v1/calculate", CalcHanlder)
	fmt.Println("Server is starting on port:", a.config.Addr)
	err := http.ListenAndServe(":"+a.config.Addr, r)
	if err != nil {
		fmt.Println("Error starting server:", err)
		return err
	}
	return nil
}
