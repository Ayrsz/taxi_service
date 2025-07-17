package database

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"taxi_service/models"
)

// Implementação do repositório JSON para dummy users
type JSONDummyUserRepository struct {
	filePath string
	mutex    sync.RWMutex
}

func NewJSONDummyUserRepository() *JSONDummyUserRepository {
	return &JSONDummyUserRepository{
		filePath: "./data/dummy_users.json",
	}
}

func (r *JSONDummyUserRepository) criarArquivoSeNaoExistir() error {
	if _, err := os.Stat(r.filePath); os.IsNotExist(err) {
		// Criar diretório se não existir
		if err := os.MkdirAll("./data", 0755); err != nil {
			return err
		}

		// Criar arquivo vazio
		file, err := os.Create(r.filePath)
		if err != nil {
			return err
		}
		defer file.Close()

		// Escrever array vazio
		_, err = file.WriteString("[]")
		return err
	}
	return nil
}

func (r *JSONDummyUserRepository) lerTodos() ([]models.DummyUser, error) {
	if err := r.criarArquivoSeNaoExistir(); err != nil {
		return nil, err
	}

	data, err := os.ReadFile(r.filePath)
	if err != nil {
		return nil, err
	}

	var users []models.DummyUser
	if err := json.Unmarshal(data, &users); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *JSONDummyUserRepository) salvarTodos(users []models.DummyUser) error {
	data, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(r.filePath, data, 0644)
}

func (r *JSONDummyUserRepository) BuscarPorEmail(email string) (*models.DummyUser, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	users, err := r.lerTodos()
	if err != nil {
		return nil, err
	}

	for _, user := range users {
		if user.Email == email {
			return &user, nil
		}
	}

	return nil, fmt.Errorf("user não encontrado")
}

func (r *JSONDummyUserRepository) BuscarPorID(id int) (*models.DummyUser, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	users, err := r.lerTodos()
	if err != nil {
		return nil, err
	}

	for _, user := range users {
		if user.ID == id {
			return &user, nil
		}
	}

	return nil, fmt.Errorf("user não encontrado")
}

func (r *JSONDummyUserRepository) Criar(user *models.DummyUser) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	users, err := r.lerTodos()
	if err != nil {
		return err
	}

	// Gerar novo ID
	maxID := 0
	for _, u := range users {
		if u.ID > maxID {
			maxID = u.ID
		}
	}
	user.ID = maxID + 1

	users = append(users, *user)
	return r.salvarTodos(users)
}

func (r *JSONDummyUserRepository) Atualizar(id int, user *models.DummyUser) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	users, err := r.lerTodos()
	if err != nil {
		return err
	}

	for i, u := range users {
		if u.ID == id {
			user.ID = id
			users[i] = *user
			return r.salvarTodos(users)
		}
	}

	return fmt.Errorf("user não encontrado")
}

func (r *JSONDummyUserRepository) Excluir(id int) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	users, err := r.lerTodos()
	if err != nil {
		return err
	}

	for i, u := range users {
		if u.ID == id {
			users = append(users[:i], users[i+1:]...)
			return r.salvarTodos(users)
		}
	}

	return fmt.Errorf("user não encontrado")
}
