package services

import (
	"taxi-service/models"
	"time"
	"strconv"
	"errors"
	"encoding/json"
	"os"
	"log"
	"fmt"
	"sync"
)

var corridas []models.Corrida 



// CorridaService gerencia a lógica de negócio das corridas.
type CorridaService struct {
	corridas map[int]*models.Corrida
	mutex    sync.RWMutex
	nextID   int
}

// NewCorridaService cria uma nova instância de CorridaService.
func NewCorridaService() *CorridaService {
	service := &CorridaService{
		corridas: make(map[int]*models.Corrida),
		nextID:   1,
	}
	// Inicia o monitoramento em background
	go service.MonitorarCorridasAtivas()
	return service
}

// CriarNovaCorrida cria uma nova corrida e a prepara para ser aceita.
func (s *CorridaService) CriarNovaCorrida(corridaInput models.Corrida) (*models.Corrida, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	corrida := &corridaInput
	corrida.ID = s.nextID
	s.nextID++
	corrida.Status = models.StatusProcurandoMotorista
	corrida.DataInicio = time.Now()
	// Em um sistema real, o tempo estimado seria calculado com base na distância, trânsito, etc.
	// Para este exemplo, vamos fixar em 1 minuto para facilitar os testes.
	corrida.TempoEstimado = 1 // minutos

	s.corridas[corrida.ID] = corrida

	return corrida, nil
}

// GetCorridaPorID busca uma corrida pelo seu ID.
func (s *CorridaService) GetCorridaPorID(id int) (*models.Corrida, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	corrida, exists := s.corridas[id]
	if !exists {
		return nil, fmt.Errorf("corrida com ID %d não encontrada, na procura", id)
	}
	return corrida, nil
}

// AceitarCorrida permite que um motorista aceite uma corrida.
func (s *CorridaService) AceitarCorrida(corridaID int, motoristaID int) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	corrida, exists := s.corridas[corridaID]
	if !exists {
		return fmt.Errorf("corrida com ID %d não encontrada", corridaID)
	}

	if corrida.Status != models.StatusProcurandoMotorista {
		return fmt.Errorf("corrida %d não está mais procurando por motorista", corridaID)
	}

	corrida.Status = models.StatusMotoristaEncontrado
	corrida.MotoristaID = motoristaID
	fmt.Printf("Corrida %d: Motorista %d aceitou a corrida.\n", corrida.ID, corrida.MotoristaID)

	return nil
}

// AtualizarPosicao atualiza a localização do motorista para uma corrida específica.
func (s *CorridaService) AtualizarPosicao(corridaID int, lat, lng float64) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	corrida, exists := s.corridas[corridaID]
	if !exists {
		return fmt.Errorf("corrida com ID %d não encontrada", corridaID)
	}

	corrida.MotoristaLat = lat
	corrida.MotoristaLng = lng
	return nil
}

// FinalizarCorrida finaliza uma corrida, aplicando a lógica de tempo.
func (s *CorridaService) FinalizarCorrida(corridaID int) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	corrida, exists := s.corridas[corridaID]
	if !exists {
		return fmt.Errorf("corrida com ID %d não encontrada", corridaID)
	}

	duracaoReal := time.Since(corrida.DataInicio)
	duracaoEstimada := time.Duration(corrida.TempoEstimado) * time.Minute

	if duracaoReal < duracaoEstimada {
		corrida.Status = models.StatusConcluidaAntecedencia
		corrida.BonusAplicado = true
		fmt.Printf("Corrida %d: Finalizada com antecedência! Bônus aplicado.\n", corrida.ID)
	} else if duracaoReal > duracaoEstimada+time.Duration(15)*time.Minute { // Limite de tolerância para cancelamento
		corrida.Status = models.StatusCanceladaPorExcessoTempo
		fmt.Printf("Corrida %d: Cancelada por excesso de tempo.\n", corrida.ID)
	} else if duracaoReal > duracaoEstimada {
		corrida.Status = models.StatusAtrasado
		fmt.Printf("Corrida %d: Finalizada com atraso.\n", corrida.ID)
	} else {
		corrida.Status = models.StatusConcluidaNoTempo
		fmt.Printf("Corrida %d: Finalizada no tempo previsto.\n", corrida.ID)
	}

	now := time.Now()
    corrida.DataFim = &now

	return nil
}

// MonitorarCorridasAtivas é um processo em background para verificar status.
func (s *CorridaService) MonitorarCorridasAtivas() {
	// Ticker para verificar a cada 30 segundos
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		s.mutex.Lock()
		for _, corrida := range s.corridas {
			// Apenas verifica corridas que estão em andamento
			if corrida.Status == models.StatusMotoristaEncontrado || corrida.Status == models.StatusCorridaIniciada {
				duracaoReal := time.Since(corrida.DataInicio)
				duracaoEstimada := time.Duration(corrida.TempoEstimado) * time.Minute

				// Lógica para cancelamento automático
				if duracaoReal > duracaoEstimada+time.Duration(15)*time.Minute {
					corrida.Status = models.StatusCanceladaPorExcessoTempo
					now := time.Now()
    				corrida.DataFim = &now
					fmt.Printf("Corrida %d: Cancelada automaticamente por excesso de tempo.\n", corrida.ID)
				} else if duracaoReal > duracaoEstimada && corrida.Status != models.StatusAtrasado {
					// Lógica para marcar como atrasado
					corrida.Status = models.StatusAtrasado
					fmt.Printf("Corrida %d: Marcada como atrasada.\n", corrida.ID)
				}
			}
		}
		s.mutex.Unlock()
	}
}

func AvaliarCorrida(id int, nota int) error {
    for i := range corridas {
        if corridas[i].ID == id {
            corridas[i].Avaliacao = &nota
            return nil
        }
    }
    return errors.New("corrida não encontrada")
}

func (s *CorridaService) AdicionarCorrida(corrida models.Corrida) {
	corridas = append(corridas, corrida)
}

func CarregarCorridasDoArquivo() {
	file, err := os.Open("data/corridas.json")
	if err != nil {
		log.Println("Erro ao abrir arquivo JSON de corridas:", err)
		return
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&corridas)
	if err != nil {
		log.Println("Erro ao fazer parse do JSON:", err)
	}
}

func GetCorridas() []models.Corrida {
	return corridas
}

func (s *CorridaService) IniciarCorrida(corridaID int) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	corrida, exists := s.corridas[corridaID]
	if !exists {
		return fmt.Errorf("corrida com ID %d não encontrada", corridaID)
	}

	// Apenas permite iniciar se o motorista já foi encontrado
	if corrida.Status != models.StatusMotoristaEncontrado {
		return fmt.Errorf("corrida %d não pode ser iniciada, pois seu status é '%s'", corridaID, corrida.Status)
	}

	corrida.Status = models.StatusEmAndamento
	fmt.Printf("Corrida %d: Iniciada e agora está em andamento.\n", corrida.ID)

	return nil
}


// ALTERADO: A função agora aceita o ID do motorista como string para alinhar com o modelo e o controller.
func (s *CorridaService) CancelarCorridaPeloMotorista(corridaID int, motoristaIDStr string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Converte o ID string do motorista para int para comparação,
	// já que o campo na struct Corrida é int.
	motoristaID, err := strconv.Atoi(motoristaIDStr)
	if err != nil {
		return fmt.Errorf("ID do motorista inválido: '%s'", motoristaIDStr)
	}

	corrida, exists := s.corridas[corridaID]
	if !exists {
		return fmt.Errorf("corrida com ID %d não encontrada", corridaID)
	}

	if corrida.MotoristaID != motoristaID {
		return fmt.Errorf("motorista %d não pode cancelar a corrida %d, status %s", motoristaID, corridaID, corrida.Status)
	}

	switch corrida.Status {
	case models.StatusConcluidaAntecedencia,
		models.StatusConcluidaNoTempo,
		models.StatusCanceladaPeloUsuario,
		models.StatusCanceladaPeloMotorista,
		models.StatusEmAndamento,
		models.StatusCanceladaPorExcessoTempo:
		return fmt.Errorf("corrida %d não pode ser cancelada pois seu status é '%s'", corridaID, corrida.Status)
	}

	corrida.Status = models.StatusCanceladaPeloMotorista
	now := time.Now()
	corrida.DataFim = &now
	fmt.Printf("Corrida %d: Cancelada pelo motorista %d.\n", corrida.ID, motoristaID)

	return nil
}