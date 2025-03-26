package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

// Estruturas para armazenar os dados
type PostData struct {
	ImageURL string `json:"image_url"`
	Message  string `json:"message"`
}

type TextData struct {
	Message string `json:"message"`
}

type LogEntry struct {
	Timestamp string `json:"timestamp"`
	Source    string `json:"source"`
	Message   string `json:"message"`
	ImageURL  string `json:"image_url,omitempty"`
}

// Caminhos dos arquivos
const dataFile = "data.json"
const logFile = "log.txt"

// Salvar a última mensagem no arquivo JSON
func saveData(data PostData) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(dataFile, jsonData, 0644)
}

// Salvar a mensagem no log
func appendToLog(source, message, imageURL string) error {
	// Criar a string de log
	logEntry := fmt.Sprintf("[%s] (%s) Mensagem: %s",
		time.Now().Format("2006-01-02 15:04:05"), source, message)

	if imageURL != "" {
		logEntry += fmt.Sprintf(" | Imagem: %s", imageURL)
	}

	logEntry += "\n"

	// Abrir o arquivo em modo append (adicionar ao final)
	f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	// Escrever no arquivo
	_, err = f.WriteString(logEntry)
	return err
}

// Handler para GET "/"
func helloWorldHandler(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadFile(dataFile)
	if err != nil {
		http.Error(w, "Erro ao carregar dados", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

// Handler para POST "/post" (mensagem + imagem)
func postHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	var data PostData
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Erro ao processar JSON", http.StatusBadRequest)
		return
	}

	// Salvar a última mensagem
	if err := saveData(data); err != nil {
		http.Error(w, "Erro ao salvar dados", http.StatusInternalServerError)
		return
	}

	// Registrar no log como "Assistente"
	if err := appendToLog("Assistente", data.Message, data.ImageURL); err != nil {
		http.Error(w, "Erro ao salvar no log", http.StatusInternalServerError)
		return
	}

	// Responder ao cliente
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message":  fmt.Sprintf("Mensagem recebida: %s", data.Message),
		"imageURL": data.ImageURL,
	})
}

// Handler para POST "/text" (apenas texto)
func textHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	var data TextData
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Erro ao processar JSON", http.StatusBadRequest)
		return
	}

	// Registrar no log como "Usuário"
	if err := appendToLog("Usuário", data.Message, ""); err != nil {
		http.Error(w, "Erro ao salvar no log", http.StatusInternalServerError)
		return
	}

	// Responder ao cliente
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": fmt.Sprintf("Ele me respondeu com a frase %s , estou satisfeita com a interacao! ", data.Message),
	})
}

// Handler para GET "/logs" (exibe o conteúdo do log)
func logsHandler(w http.ResponseWriter, r *http.Request) {
	// Ler o arquivo de log
	data, err := ioutil.ReadFile(logFile)
	if err != nil {
		http.Error(w, "Erro ao carregar logs", http.StatusInternalServerError)
		return
	}

	// Converter os logs para um formato mais legível em JSON
	var logEntries []LogEntry
	logLines := string(data)

	// Dividir o arquivo de log em linhas
	lines := splitLines(logLines)
	for _, line := range lines {
		// Ignorar linhas vazias
		if line == "" {
			continue
		}

		// Processar cada linha de log e extrair os dados usando expressões regulares
		logEntry := parseLogLine(line)
		logEntries = append(logEntries, logEntry)
	}

	// Enviar os logs em formato JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logEntries)
}

// Função para dividir as linhas do log
func splitLines(logData string) []string {
	return strings.Split(logData, "\n")
}

// Função para analisar uma linha de log e extrair as informações
func parseLogLine(line string) LogEntry {
	// Usar uma expressão regular para capturar os dados da linha
	// O formato da expressão regular deve corresponder ao formato do seu log
	// Agora, a captura do campo de imagem será opcional
	re := regexp.MustCompile(`\[(.*?)\] \((.*?)\) Mensagem: (.*?) \|? Imagem: (.*)`)
	matches := re.FindStringSubmatch(line)

	if len(matches) == 0 {
		// Verificar se a linha é do tipo "Usuário" sem imagem
		reUser := regexp.MustCompile(`\[(.*?)\] \((.*?)\) Mensagem: (.*?)$`) // Alterado para garantir captura até o final
		matchesUser := reUser.FindStringSubmatch(line)

		if len(matchesUser) > 0 {
			// Linha do usuário sem imagem
			return LogEntry{
				Timestamp: matchesUser[1],
				Source:    matchesUser[2],
				Message:   matchesUser[3],
				ImageURL:  "", // Garantir que o campo de imagem fique vazio
			}
		}

		// Caso não consiga extrair os dados, retorna um LogEntry com valores padrão
		return LogEntry{}
	}

	// Caso encontre os dados com a imagem
	return LogEntry{
		Timestamp: matches[1],
		Source:    matches[2],
		Message:   matches[3],
		ImageURL:  matches[4],
	}
}

func main() {
	http.HandleFunc("/", helloWorldHandler)
	http.HandleFunc("/post", postHandler)
	http.HandleFunc("/text", textHandler)
	http.HandleFunc("/logs", logsHandler) // Novo endpoint para obter os logs

	log.Println("Servidor rodando na porta 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
