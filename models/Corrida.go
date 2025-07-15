package models

import (
	"time"
)

type Corrida struct {
	Id int `json:"id"`
	Data  string `json:"data"` //dia da corrida
	Horario time.Time `json:"horario"` // horário de inicio 
	Tempo int `json:"tempo"` // tempo para chegar ao destino
	Valor int `json:"valor"` // valor da corrida
	Avaliacao *int `json:"avaliacao"` // avaliacao 1, 2, 3, 4, 5 ou nil, * permite nil
	Status string  `json:"status"`// andamento, finalizada ou cancelada
	CPFMotorista *int `json:"cpfMotorista"` // chave estrangeira pro motorista responsavel
	LocalEmbarque string `json:"localEmbarque"` 
	LocalDesembarque string `json:"localDesembarque"`
}

const (
	StatusAndamento = "andamento"
	StatusCancelada = "cancelada"
	StatusFinalizada = "finalizada"
	StatusPendente = "pendente" //Ainda não aceita
)