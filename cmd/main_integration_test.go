package main

import (
	"CRUD_API/pkg/handler"
	"CRUD_API/pkg/repository"
	"CRUD_API/pkg/service"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"testing"
	"time"
)

// docker run --name=testdb -e POSTGRES_PASSWORD=qwerty -e POSTGRES_USER=postgres -e POSTGRES_DB=postgres -e POSTGRES_HOST=localhost -p 5432:5432 -d postgres
func setupTestDatabase() (*pgxpool.Pool, error) {
	//exec.Command-создает команду в терминале
	cmd := exec.Command("docker", "run", "--name=testdb", "-e", "POSTGRES_PASSWORD=qwerty",
		"-e", "POSTGRES_USER=postgres", "-e", "POSTGRES_DB=postgres", "-p", "5432:5432", "-d", "postgres")
	if err := cmd.Run(); err != nil { //запускает команду в терминале
		log.Fatalf("Error building Docker image: %v", err)
	}

	migrateCmd := exec.Command("migrate", "-path", "../schema", "-database", "postgres://postgres:qwerty@localhost:5432/postgres?sslmode=disable", "up")
	/*../schema-используем две точки указывают на то,что папка находится на 1 уровень выше текущей директории*/
	//вывод о том,что миграции применились:
	migrateCmd.Stdout = os.Stdout
	migrateCmd.Stderr = os.Stderr
	time.Sleep(2 * time.Second) //даем серверу время запуститься перед применением миграций
	if err := migrateCmd.Run(); err != nil {
		log.Fatalf("Error applying migrations: %v", err)
	}

	time.Sleep(2 * time.Second) //даем докеру время создать контейнер,для дальнейшего подкл к нему
	/*pgConfig, err := pgxpool.ParseConfig("postgres://postgres:qwerty@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		log.Fatalf("Error parsing Postgres config: %v", err)
	}*/
	data := "postgres://postgres:qwerty@localhost:5432/postgres?sslmode=disable"
	db, err := pgxpool.Connect(context.Background(), data)
	if err != nil {
		log.Fatalf("Error connecting to Postgres: %v", err)
	}

	return db, nil
}

func tearDownTestDatabase() error {
	cmd := exec.Command("docker", "stop", "testdb")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error stopping Docker container: %w", err)
	}

	cmd = exec.Command("docker", "rm", "testdb")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error removing Docker container: %v", err)
	}
	return nil
}

func TestIntegration(t *testing.T) {
	db, err := setupTestDatabase()
	if err != nil {
		fmt.Printf("err in DBtest:%s", err)
	}
	defer db.Close()
	handler.InitRedis()
	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)
	r := handlers.InitRoutes()
	serv := httptest.NewServer(r) //тестируемый http serv
	defer serv.Close()
	signUpReq := map[string]interface{}{
		"name":     "TestName",
		"username": "TestDescription",
		"password": "TestPass",
	}
	signUpJson, _ := json.Marshal(signUpReq)
	signUp, err := http.Post(serv.URL+"/auth/sign-up", "application/json", bytes.NewBuffer(signUpJson))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, signUp.StatusCode)

	signInReq := map[string]interface{}{
		"id":       1,
		"username": "TestDescription",
		"password": "TestPass",
	}
	signInReqJson, _ := json.Marshal(signInReq)
	signInResp, err := http.Post(serv.URL+"/auth/sign-in", "application/json", bytes.NewBuffer(signInReqJson))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, signInResp.StatusCode)
	var signUpResponseData map[string]interface{}
	_ = json.NewDecoder(signInResp.Body).Decode(&signUpResponseData)
	accessToken := signUpResponseData["access_token"].(string)

	createProductRequest := map[string]interface{}{
		"name":        "TestProduct",
		"description": "TestDescription",
	}
	createProductJSON, _ := json.Marshal(createProductRequest) //преобразуем в json
	req, err := http.NewRequest("POST", serv.URL+"/api/product", bytes.NewBuffer(createProductJSON))
	/*serv.URL-url тестового сервера,bytes.NewBuffer преобразует наше тело запроса к нужному io.reader*/
	req.Header.Set("Authorization", "Bearer "+accessToken)
	createProduct, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, createProduct.StatusCode)

	req, err = http.NewRequest("GET", serv.URL+"/api/product", nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)
	readAllProducts, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, readAllProducts.StatusCode)

	req, err = http.NewRequest("GET", serv.URL+"/api/product/1", nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)
	readProductById, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, readProductById.StatusCode)

	updateProductRequest := map[string]interface{}{ //наш обновленный запрос
		"name":        "UpdatedProduct",
		"description": "UpdatedDescription",
	}
	updateProductJSON, _ := json.Marshal(updateProductRequest)
	req, err = http.NewRequest("PUT", serv.URL+"/api/product/1", bytes.NewBuffer(updateProductJSON)) //create new http req
	req.Header.Set("Authorization", "Bearer "+accessToken)
	updateProduct, err := http.DefaultClient.Do(req) //perform our http req
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, updateProduct.StatusCode)

	req, err = http.NewRequest("DELETE", serv.URL+"/api/product/1", nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)
	deleteProduct, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, deleteProduct.StatusCode)

	if err := tearDownTestDatabase(); err != nil {
		t.Errorf("err in DeleteOurDb:%s", err)
	}
}
