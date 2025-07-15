package database

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"taxi_service/models"
)

type CorridaRepository interface {
	BuscarPorID(id int) (*models.Corrida, error)
	Salvar(corrida *models.Corrida) error
	ListarTodas() ([]models.Corrida, error)
}

type jsonCorridaRepository struct {
	filePath string
	mutex    sync.RWMutex
}

func NewJSONCorridaRepository(filePath string) CorridaRepository {
	repo := &jsonCorridaRepository{
		filePath: filePath,
	}
	repo.garantirArquivoExistente()
	return repo
}

func (r *jsonCorridaRepository) garantirArquivoExistente() {
	r.mutex.Lock()
	defer r.mutex.Unlock()
    if err := os.MkdirAll("./data", 0755); err != nil {
        panic(fmt.Sprintf("não foi possível criar o diretório data: %v", err))
    }
	if _, err := os.Stat(r.filePath); os.IsNotExist(err) {
		os.WriteFile(r.filePath, []byte("[]"), 0644)
	}
}

func (r *jsonCorridaRepository) lerCorridas() ([]models.Corrida, error) {
	data, err := os.ReadFile(r.filePath)
	if err != nil {
		return nil, fmt.Errorf("falha ao ler arquivo de corridas: %w", err)
	}

	var corridas []models.Corrida
	if err := json.Unmarshal(data, &corridas); err != nil {
		return nil, fmt.Errorf("falha ao decodificar JSON de corridas: %w", err)
	}
	return corridas, nil
}

func (r *jsonCorridaRepository) salvarCorridas(corridas []models.Corrida) error {
	data, err := json.MarshalIndent(corridas, "", "  ")
	if err != nil {
		return fmt.Errorf("falha ao codificar corridas para JSON: %w", err)
	}
	return os.WriteFile(r.filePath, data, 0644)
}

func (r *jsonCorridaRepository) BuscarPorID(id int) (*models.Corrida, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	corridas, err := r.lerCorridas()
	if err != nil {
		return nil, err
	}

	for _, c := range corridas {
		if c.Id == id {
			return &c, nil
		}
	}
	return nil, fmt.Errorf("corrida com ID %d não encontrada", id)
}

func (r *jsonCorridaRepository) Salvar(corrida *models.Corrida) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	corridas, err := r.lerCorridas()
	if err != nil {
		return err
	}

	encontrado := false
	for i, c := range corridas {
		if c.Id == corrida.Id {
			corridas[i] = *corrida // Atualiza a corrida existente
			encontrado = true
			break
		}
	}

	if !encontrado {
		//Gera um novo ID se for uma nova corrida
		maxID := 0
		for _, c := range corridas {
			if c.Id > maxID {
				maxID = c.Id
			}
		}
		corrida.Id = maxID + 1
		corridas = append(corridas, *corrida)
	}

	return r.salvarCorridas(corridas)
}

func (r *jsonCorridaRepository) ListarTodas() ([]models.Corrida, error) {
    r.mutex.RLock()
    defer r.mutex.RUnlock()
    return r.lerCorridas()
}
