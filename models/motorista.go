package models

import "gorm.io/gorm"

type Motorista struct {
	gorm.Model
	Nome  string `json:"nome"` 
	CPF int `json:"cpf"` //apenas um por cpf, chave primaria
	Idade int `json:"idade"` //precisa ser maior que 18
	CNH int `json:"cnh"` 
	CRLV int `json:"crlv"`
}
