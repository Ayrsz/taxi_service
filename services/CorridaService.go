package services

import (
	"fmt"
	"taxi_service/database"
	"taxi_service/models"
	"time" // Import necessário para lidar com o tempo
)

// CorridaService encapsula a lógica de negócio para corridas.
type CorridaService struct {
	repo database.CorridaRepository
}

// NewCorridaService cria uma nova instância de CorridaService.
func NewCorridaService(repo database.CorridaRepository) *CorridaService {
	return &CorridaService{repo: repo}
}

// CancelarCorrida aplica a regra de negócio para o cancelamento.
func (s *CorridaService) CancelarCorrida(id int) (*models.Corrida, error) {
	corrida, err := s.repo.BuscarPorID(id)
	if err != nil {
		return nil, err
	}

	// Regra de negócio: Apenas corridas "pendentes" podem ser canceladas.
	if corrida.Status != models.StatusPendente {
		return nil, fmt.Errorf("não é possível cancelar uma corrida que não está em andamento. Status atual: %s", corrida.Status)
	}

	corrida.Status = models.StatusCancelada

	if err := s.repo.Salvar(corrida); err != nil {
		return nil, fmt.Errorf("falha ao salvar o cancelamento da corrida: %w", err)
	}

	return corrida, nil
}

// NOVA FUNÇÃO: Verifica corridas pendentes que excederam o tempo limite.
func (s *CorridaService) VerificarCorridasPendentesPorTimeout() error {
	corridas, err := s.repo.ListarTodas()
	if err != nil {
		return fmt.Errorf("falha ao listar corridas: %w", err)
	}

	const limiteMinutos = 10.0 // Limite de 10 minutos para o motorista partir

	for i := range corridas {
		corrida := &corridas[i] // Usamos um ponteiro para modificar a corrida original no slice

		// Só processa corridas pendentes
		if corrida.Status != models.StatusPendente {
			continue
		}

		// Calcula o tempo decorrido desde a criação/alocação da corrida
		tempoDecorrido := time.Since(corrida.Horario).Minutes()

		if tempoDecorrido > limiteMinutos {
			corrida.Status = models.StatusCancelada
			// Salva a alteração no repositório
			if err := s.repo.Salvar(corrida); err != nil {
				// Apenas loga o erro, mas continua o processo para outras corridas
				fmt.Printf("Erro ao salvar corrida %d cancelada: %v\n", corrida.Id, err)
			}
		}
	}
	return nil
}
