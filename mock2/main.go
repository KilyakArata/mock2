package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Структура для ответа нашего сервиса
type Student struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// Функция, которая делает запрос к внешнему API и возвращает данные о пользователе
func getUserFromAPI(userID string) (*Student, error) {
	apiURL := fmt.Sprintf("https://jsonplaceholder.typicode.com/users/%s", userID)

	// Выполняем запрос к API
	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("ошибка при выполнении запроса: %v", err)
	}
	defer resp.Body.Close()

	// Читаем тело ответа
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ошибка при чтении ответа: %v", err)
	}

	// Декодируем JSON из API в структуру Student
	var student Student
	err = json.Unmarshal(body, &student)
	if err != nil {
		return nil, fmt.Errorf("ошибка при декодировании JSON: %v", err)
	}

	return &student, nil
}

// Handler для маршрута /getUser/{id}
func getUserHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем ID из параметров URL
	vars := mux.Vars(r)
	userID := vars["id"]

	// Получаем данные пользователя из внешнего API
	student, err := getUserFromAPI(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Устанавливаем заголовки ответа и возвращаем JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(student)
}

func main() {
	// Создаем новый роутер
	router := mux.NewRouter()

	// Определяем маршрут /getUser/{id}
	router.HandleFunc("/getUser/{id}", getUserHandler).Methods("GET")

	// Запуск HTTP-сервера на порту 8080
	fmt.Println("Сервер запущен на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
