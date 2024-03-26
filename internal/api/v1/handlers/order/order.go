package order

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"text/template"

	"github.com/plusik10/cmd/order-info-service/internal/service"
)

type responseErr struct {
	Err        string `json:"err,omitempty"`
	StatusCode int    `json:"statusCode,omitempty"`
}

func sendResponseError(w http.ResponseWriter, err error, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	resp := responseErr{Err: err.Error(), StatusCode: statusCode}

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
	}
}

func Info(ctx context.Context, orderService service.OrderService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		order, err := orderService.GetOrderByUID(ctx, id)
		if err != nil {
			sendResponseError(w, err, 400)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(order)
		if err != nil {
			sendResponseError(w, err, 500)
		}
	}
}

func GetOrderUIDs(ctx context.Context, service service.OrderService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := struct {
			OrderUIDs []string
		}{}

		arrStr, err := service.GetOrderUIDs(ctx)
		if err != nil {
			fmt.Println("error getting order")
			return
		}

		data.OrderUIDs = arrStr

		tmpl, err := template.ParseFiles("templates/index.html")
		if err != nil {
			fmt.Println("error parsing templates")
			return
		}
		err = tmpl.Execute(w, data)
		if err != nil {
			fmt.Println("error executing order", err)
			return
		}
	}
}
