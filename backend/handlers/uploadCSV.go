package handlers

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"home_money/models"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type CSVRow map[string]string

type CategoryGPT struct {
	CategoryName string
	Description  string
}

type Response struct {
	Message string `json:"message"`
}

type ImportCSVHandler struct {
	db              *gorm.DB
	categoryHandler *CategoryHandler
}

// CategoriesStruct representa a estrutura do JSON
type CategoriesStruct struct {
	Categories []string `json:"categorias"`
}

func NewImportCSVHandler(db *gorm.DB, categoryHandler *CategoryHandler) *ImportCSVHandler {
	return &ImportCSVHandler{
		db:              db,
		categoryHandler: categoryHandler,
	}
}

func (h *ImportCSVHandler) UploadCSV(c echo.Context) error {
	// Extrair o user_id do banco de dados usando o email do token JWT
	userID, err := getUserIDFromJWT(c, h.db)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	// Obtenha o arquivo do formulário
	file, err := c.FormFile("file")

	csvStructureBytes := []byte(c.FormValue("structure"))
	csvStructure := string(csvStructureBytes)
	columns := strings.Split(csvStructure, ";")

	// if err != nil {
	// 	fmt.Println("Erro ao obter o arquivo do formulário:", err)
	// 	return c.JSON(http.StatusBadRequest, Response{Message: "Erro ao obter o arquivo do formulário"})
	// }

	// Abra o arquivo CSV
	src, err := file.Open()
	if err != nil {
		fmt.Println("Erro ao abrir o arquivo CSV:", err)
		return c.JSON(http.StatusInternalServerError, Response{Message: "Erro ao abrir o arquivo CSV"})
	}
	defer src.Close()

	reader := csv.NewReader(src)
	// Defina o delimitador (substitua '\t' pelo delimitador que você está usando)
	reader.Comma = ';'

	var header []string

	//Todas as linhas do CSV
	var rows []CSVRow

	// Leia todas as linhas do arquivo CSV (ignorando a primeira linha)
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Erro ao ler o arquivo CSV:", err)
		return c.JSON(http.StatusInternalServerError, Response{Message: "Erro ao ler o arquivo CSV"})
	}

	var index = 0

	for _, record := range records {
		headerIndex := 0

		if index == 0 {
			// Ordenar o slice
			sort.Strings(columns)

			for _, colName := range record {

				// Verificar se a string está presente no slice
				indice := sort.SearchStrings(columns, colName)

				// Verificar se a string foi encontrada ou não
				if indice < len(columns) && columns[indice] == colName {
					header = append(header, colName)
				}
			}

			if len(header) == 0 {
				response := Response{
					Message: "Não foi definido o mapeamento do cabeçalho!",
				}

				return c.JSON(http.StatusOK, response)
			}

		}

		//Linha do CSV
		row := make(CSVRow)

		for _, field := range record {
			row[header[headerIndex]] = field

			headerIndex++
		}
		// Adicione a linha processada ao slice de rows
		if index != 0 {
			rows = append(rows, row)
		}

		index++
	}

	// Tamanho do pedaço
	chunkSize := 60
	var categories []*CategoryGPT

	// Iterar sobre o array em pedaços de 10
	for i := 0; i < len(rows); i += chunkSize {
		end := i + chunkSize
		if end > len(rows) {
			end = len(rows)
		}

		// Criar o pedaço
		chunk := rows[i:end]

		h.createCategoriesByDescription(chunk, userID, &categories)
	}

	//REGISTRA NO BANCO

	if len(categories) > 0 {
		for _, category := range categories {

			newCategory := new(models.Category)

			// Define o campo user_id no objeto Category
			newCategory.CategoryName = category.CategoryName
			newCategory.UserID = userID

			// Verificar se a categoria com o mesmo nome já existe
			if !h.categoryHandler.categoryExists(category.CategoryName) {

				result := h.db.Create(newCategory)

				if result.Error != nil {
					fmt.Println("ERRO AO CRIAR CATEGORIA:", result.Error)
				}

				fmt.Println(newCategory.CategoryID)

				descriptionsByCategories := new(models.DescriptionsByCategories)

				description := removeExtraSpaces(category.Description)
				descriptionFormated := strings.TrimSuffix(description, string(description[len(description)-1]))

				descriptionsByCategories.CategoryID = newCategory.CategoryID
				descriptionsByCategories.Description = descriptionFormated

				result2 := h.db.Create(descriptionsByCategories)

				if result2.Error != nil {
					fmt.Println("ERRO AO CRIAR ASSOCIAÇÃO DE CATEGORIA COM DESCRIÇÃO:", result2.Error)
				}
			} else {
				categoryID := h.categoryHandler.GetCategoryByName(category.CategoryName, userID)

				description := removeExtraSpaces(category.Description)
				descriptionFormated := strings.TrimSuffix(description, string(description[len(description)-1]))

				descriptionsByCategories := new(models.DescriptionsByCategories)
				descriptionsByCategories.CategoryID = categoryID
				descriptionsByCategories.Description = descriptionFormated

				result2 := h.db.Create(descriptionsByCategories)

				if result2.Error != nil {
					fmt.Println("ERRO AO CRIAR ASSOCIAÇÃO DE CATEGORIA COM DESCRIÇÃO 2:", result2.Error)
				} else {
					fmt.Println("ASSOCIOU DESCRIÇÃO A CATEGORIA")
				}
			}
		}
	}

	h.saveTransactions(rows, userID)

	// Retorne uma resposta indicando que a importação foi bem-sucedida
	// response := Response{
	// 	Message: "Arquivo CSV importado com sucesso!",
	// }

	return c.JSON(http.StatusOK, rows)
}

const (
	TRANSACTION_STATUS_PENDING = 1
)

func (h *ImportCSVHandler) saveTransactions(rows []CSVRow, userID int) {

	for _, row := range rows {
		newTransaction := new(models.Transaction)
		newTransaction.AccountID = 1 //ESSE DEVE VIR POR RESQUEST

		// Converte a string para time.Time
		dateConverted, err := time.Parse("02/01/2006", row["Data Lançamento"])
		if err != nil {
			panic(err)
		}

		newTransaction.Date = dateConverted

		var transactionType string

		// Remover espaços em branco
		strNumero := strings.TrimSpace(row["Valor"])

		// Remover todos os pontos
		strNumero = strings.ReplaceAll(strNumero, ".", "")

		// Substituir a última vírgula por ponto
		strNumero = strings.Replace(strNumero, ",", ".", 1)

		// Convertendo a string para um número de ponto flutuante
		value, err := strconv.ParseFloat(strNumero, 64)

		// Verificando se houve algum erro na conversão
		if err != nil {
			fmt.Println("Erro na conversão:", err)
			return
		}

		if value > 0 {
			transactionType = "I" //I - ENTRADA
		} else {
			transactionType = "O" //O - SAÍDA

		}
		newTransaction.TransactionType = transactionType //AQUI FAZER PODER RECEBER POR RESQUEST
		newTransaction.AccountID = 4                     //ESSE ID DA CONTA TAMBÉM TEM QUE VIR POR REQUEST
		newTransaction.Amount = value
		newTransaction.UserID = userID
		newTransaction.Description = removeExtraSpaces(row["Descrição"])
		newTransaction.Status = TRANSACTION_STATUS_PENDING

		categoryID, err := h.categoryHandler.GetCategoryIDByDescription(removeExtraSpaces(row["Descrição"]), userID)
		if err != nil || categoryID == 0 {
			// Tratar o erro
			fmt.Println("Erro ao buscar categoria com descrição:", removeExtraSpaces(row["Descrição"]))
			newTransaction.CategoryID = 1 //Esse aqui é o Outros
		} else {
			newTransaction.CategoryID = categoryID
		}

		fmt.Println("ID DA CATEGORIA:", newTransaction.CategoryID)

		result := h.db.Create(newTransaction)

		if result.Error != nil {
			fmt.Println("ERRO AO CRIAR TRANSAÇÃO:", result.Error)
		}

		fmt.Println("Transação criada com ID:", newTransaction.TransactionID)
	}

}

func removeExtraSpaces(frase string) string {
	palavras := strings.Fields(frase)
	resultado := strings.Join(palavras, " ")
	return resultado
}

func (h *ImportCSVHandler) createCategoriesByDescription(rows []CSVRow, userID int, categories *[]*CategoryGPT) {
	var descriptionList []string
	indexList := 0

	for _, row := range rows {
		descriptionFormatted := removeExtraSpaces(row["Descrição"])

		if !h.categoryByDescriptionExists(descriptionFormatted) {
			descriptionListString := strings.Join(descriptionList, " ")

			exists := strings.Contains(descriptionListString, descriptionFormatted)

			if !exists {
				descriptionList = append(descriptionList, descriptionFormatted+",")
				indexList++
			}
		}
	}

	if len(descriptionList) > 0 {
		descriptionListString := strings.Join(descriptionList, " ")
		fmt.Println("Lista de descrições:", descriptionListString)

		categoriesGPT := h.defineCategoriesNameWithChatGPT(descriptionListString, len(descriptionList))

		// Converte a string JSON para a estrutura CategoriesStruct
		var categoriesStr CategoriesStruct
		err := json.Unmarshal([]byte(categoriesGPT), &categoriesStr)
		if err != nil {
			fmt.Println("RESPOSTA DO CHATGPT:", categoriesGPT)
			fmt.Println("Erro ao converter JSON:", err)
		}

		//Extrai as categorias da estrutura
		extractedCategories := categoriesStr.Categories

		for index, categoryName := range extractedCategories {
			newCategory := &CategoryGPT{
				CategoryName: categoryName,
				Description:  descriptionList[index],
			}

			*categories = append(*categories, newCategory)
		}
	}
}

func (h *ImportCSVHandler) categoryByDescriptionExists(description string) bool {
	var count int64
	h.db.Model(&models.DescriptionsByCategories{}).Where("description LIKE ?", "%"+description+"%").Count(&count)
	return count > 0
}

func (h *ImportCSVHandler) defineCategoriesNameWithChatGPT(description string, qtdCategories int) string {
	// Sua chave de API do OpenAI
	apiKey := ""

	// URL da API do OpenAI
	apiURL := "https://api.openai.com/v1/chat/completions"

	//1. Retorne apenas um texto no formato de JSON assim: , não retorne nenhuma frase ou caracter fora do JSON \n 2. Você deve me retornar apenas um json com as categoria. \n 3. Crie uma lista de tipos de transações para organização financeira mais completa que você tiver \n 5. A mesma quantidade de categorias deve ser igual a quantidade de descrição de entrada \n 5. Usando ela categorize as descrições a seguir: " + description
	// Mensagens para o modelo GPT-3.5-turbo
	messages := []map[string]string{
		{
			"role":    "user",
			"content": "Com base em uma lista de categorias de transações financeiras que podem ajudar a organizar finanças pessoais. Categorize as descrições de transações com essas categorias. Você deve retornar apenas um JSON, e não pode ter nenhum comentário fora do JSON. Vai existir um array que a mesma quantidade de categorias deve ser igual a de descrições de entrada. Você tem a obrigação de retornar apenas esse JSON: { categorias: [<CATEGORIAS QUE A IA DEFINIU SEPARADAS POR VIRGULA>]}, retorne apenas" + strconv.Itoa(qtdCategories) + "categorias. As descrições para categorizar são essas a seguir:" + description,
			// "content": "Com base em uma lista de categorias de transações financeiras que podem ajudar a organizar finanças pessoais."+
			// "Categorize as transações com essas categorias."+
			// "Você sempre deve retornar apenas um JSON, e jamais/nunca retornar um texto fora do JSON."+
			// "Vai existir um array que a mesma quantidade de categorias deve ser igual a de descrições de entrada."+
			// "Você tem a obrigação de retornar apenas esse JSON: { categorias: [<CATEGORIAS QUE A IA DEFINIU SEPARADAS POR VIRGULA CASO TENHA MAIS DE UMA>]}."+
			// "Retorne a seguinte quantidade de categorias:" + strconv.Itoa(qtdCategories) +
			// "Por fim, você deve categorizar as transações a seguir:" + description,
		},
		{
			"role":    "system",
			"content": "Você é capaz de categorizar transações financeiras para organização financeira pessoal, e você só retorna um JSON no formato de texto",
		},
	}
	// Configurando a solicitação HTTP POST
	requestBody := map[string]interface{}{
		"messages": messages,
		"model":    "gpt-3.5-turbo",
	}

	req, err := http.NewRequest("POST", apiURL, encodeRequestBody(requestBody))
	if err != nil {
		fmt.Println("Erro ao criar a solicitação HTTP:", err)

	}

	// Configurando o cabeçalho de autorização com a chave de API
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	// Realizando a solicitação HTTP
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Erro ao fazer a solicitação HTTP:", err)

	}
	defer resp.Body.Close()

	// Lendo a resposta
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Erro ao ler a resposta HTTP:", err)

	}

	// Decodificando a resposta JSON
	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		fmt.Println("Erro ao decodificar resposta JSON:", err)

	}
	fmt.Println(response)
	// Extraindo a resposta como uma string
	choices, ok := response["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		fmt.Println("Erro: Resposta inválida.")

	}

	message, ok := choices[0].(map[string]interface{})["message"].(map[string]interface{})
	if !ok {
		fmt.Println("Erro ao extrair a resposta.")

	}

	content, ok := message["content"].(string)
	if !ok {
		fmt.Println("Erro ao extrair a resposta.")

	}

	return content
}

// Pegue uma lista de tipos de transações  mais completa que você tiver com 30 opções, e usando ela categorize a descrição: Pix Maycom Hoepers, você deve me retornar apenas um json com a categoria, e não deve ter mais nenhum caracter na sua resposta, apenas o nome da categoria que você definiu, igual a { category=<categoria definida pelo chatgpt> }

func encodeRequestBody(data map[string]interface{}) *bytes.Buffer {
	buffer := new(bytes.Buffer)
	_ = json.NewEncoder(buffer).Encode(data)
	return buffer
}
