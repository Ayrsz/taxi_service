package services

import (
	"errors"
	"your-app/database"
	"your-app/models"
)

// CancelarCorrida contém a regra de negócio para cancelar uma corrida.
func CancelarCorrida(id int) (*models.Corrida, error) {
	corrida, err := database.GetCorridaByID(id) // Busca a corrida no banco de dados.
	if err != nil {
		return nil, err //Corrida não encontrada
	}

	//Só pode cancelar se estiver "em andamento" ou "cancelada".
	if corrida.Status != models.StatusAndamento && corrida.Status != models.StatusPendente{
		return nil, errors.New("operação inválida: só é possível cancelar corridas que estão em andamento ou Pendentes")
	}

	//Altera o status da corrida para "cancelada".
	corrida.Status = models.StatusCancelada

	//Atualiza a corrida no banco de dados.
	err = database.UpdateCorrida(*corrida)
	if err != nil {
		return nil, err // Caso ocorra um erro ao salvar.
	}

	//Retorna a corrida com o status atualizado.
	return corrida, nil
}

// GetCorrida busca uma corrida pelo ID (função auxiliar para teste).
func GetCorrida(id int) (*models.Corrida, error) {
    return database.GetCorridaByID(id)
}
