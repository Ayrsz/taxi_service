package models

import "gorm.io/gorm"
import "time"

type Corrida struct {
	gorm.Model
	Data  string `json:"data"` //dia da corrida
	Horario time.Time `json:"horario"` // horário de inicio 
	Tempo int `json:"tempo"` // tempo para chegar ao destino
	Valor int `json:"valor"` // valor da corrida
	Avaliacao *int `json:"avaliacao"` // avaliacao 1, 2, 3, 4, 5 ou nil, * permite nil
	Status string  `json:"status"`// andamento, finalizada ou cancelada
}
