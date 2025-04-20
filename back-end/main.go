package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/google/uuid"
)

const (
	MercadoPagoBaseURL = "https://api.mercadopago.com/v1"
	AccessToken        = "" // CHAVE PRIVADA SOPHIA
)

type PaymentRequest struct {
	Token         string  `json:"token"`
	Amount        float64 `json:"amount"`
	Description   string  `json:"description"`
	Email         string  `json:"email"`
	PaymentMethod string  `json:"payment_method"`
}

type MercadoPagoRequest struct {
	TransactionAmount float64 `json:"transaction_amount"`
	Token             string  `json:"token"`
	Description       string  `json:"description"`
	Payer             struct {
		Email string `json:"email"`
	} `json:"payer"`
	PaymentMethodId string `json:"payment_method_id"`
	Installments    int    `json:"installments"`
}

type ApiResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func handlePayment(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Erro ao ler corpo da requisição", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var paymentReq PaymentRequest
	if err := json.Unmarshal(body, &paymentReq); err != nil {
		http.Error(w, "Erro ao decodificar JSON", http.StatusBadRequest)
		return
	}
	if paymentReq.Amount != 1 {
		paymentReq.Amount = 1
	}

	mpRequest := MercadoPagoRequest{
		TransactionAmount: paymentReq.Amount,
		Token:             paymentReq.Token,
		Description:       paymentReq.Description,
		PaymentMethodId:   paymentReq.PaymentMethod,
		Installments:      1,
	}
	mpRequest.Payer.Email = paymentReq.Email

	mpPayload, err := json.Marshal(mpRequest)
	if err != nil {
		http.Error(w, "Erro ao preparar requisição", http.StatusInternalServerError)
		return
	}

	req, err := http.NewRequest("POST", MercadoPagoBaseURL+"/payments", bytes.NewBuffer(mpPayload))
	if err != nil {
		http.Error(w, "Erro ao criar requisição", http.StatusInternalServerError)
		return
	}

	idempotencyKey := uuid.New().String()

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+AccessToken)
	req.Header.Set("X-Idempotency-Key", idempotencyKey)

	log.Printf("Enviando pagamento com X-Idempotency-Key: %s", idempotencyKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Erro ao processar pagamento", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	mpResponseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Erro ao ler resposta", http.StatusInternalServerError)
		return
	}

	var mpResponse interface{}
	if err := json.Unmarshal(mpResponseBody, &mpResponse); err != nil {
		http.Error(w, "Erro ao decodificar resposta", http.StatusInternalServerError)
		return
	}

	response := ApiResponse{
		Success: resp.StatusCode >= 200 && resp.StatusCode < 300,
		Message: "Processamento de pagamento concluído",
		Data:    mpResponse,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	fmt.Fprintf(w, "API funcionando!")
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/api/payment", handlePayment)

	corsHandler := enableCORS(mux)

	fmt.Println("Servidor rodando na porta 8080...")
	log.Fatal(http.ListenAndServe(":8080", corsHandler))
}
