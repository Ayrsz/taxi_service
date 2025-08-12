package utils

import "github.com/gofiber/fiber/v2"

import (
	"fmt"
	"math"
	
)

// CalculateDistance usa a fórmula de Haversine para calcular a distância entre dois pontos (em KM)
func CalculateDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371 // Raio médio da Terra em quilômetros

	var degToRad = func(deg float64) float64 {
		return deg * (math.Pi / 180)
	}

	dLat := degToRad(lat2 - lat1)
	dLon := degToRad(lon2 - lon1)

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(degToRad(lat1))*math.Cos(degToRad(lat2))*
			math.Sin(dLon/2)*math.Sin(dLon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return R * c
}

// CalculateETA (Estimated Time of Arrival) simula o tempo de chegada em minutos.
// Para um sistema real, você integraria com uma API de mapas (Google Maps, OpenStreetMap).
// A `speedKmPerMin` é uma velocidade média simulada em km/min.
func CalculateETA(targetLat, targetLon, speedKmPerMin float64) string {
	// Para este exemplo, vamos considerar uma origem fixa (posição atual do motorista ao aceitar a corrida)
	// e o targetLat/targetLon como o destino de embarque da corrida.
	// Em um cenário real, você passaria a latitude/longitude atual do motorista.
	// Por simplicidade, assumimos que esta função é chamada com a posição atual do motorista e o ponto de embarque.
	// Vamos usar um ponto de partida fictício para o cálculo da distância, se não tiver a localização do motorista.
	// A função deveria receber a localização ATUAL do motorista e a ORIGEM da corrida.
	// Como a pergunta é "Tempo estimado de chegada até o local de embarque", o `targetLat/targetLon` é a `OriginLatitude/OriginLongitude` da `Ride`.
	// E a origem do cálculo da distância seria a localização atual do motorista.
	// Para o teste, vamos simular que a "posição inicial do cálculo" é fixa para simplificar.

	// Distância fixa ou calculada do motorista para o ponto de embarque
	// Aqui, estamos simplificando, você precisaria da localização real do motorista.
	// Apenas para que o cálculo não seja 0, vamos usar um ponto de partida diferente do alvo para ter uma distância.
	// Para o Cenário 5, o ETA é para o local de EMBARQUE (OriginLatitude/OriginLongitude da corrida)
	simulatedDriverLat := -8.061738
	simulatedDriverLon := -34.880579

	distanceToPickup := CalculateDistance(simulatedDriverLat, simulatedDriverLon, targetLat, targetLon)

	if speedKmPerMin <= 0 {
		speedKmPerMin = 0.5 // Default: 30 km/h (0.5 km/min)
	}

	etaMinutes := distanceToPickup / speedKmPerMin
	if etaMinutes < 1.0 {
		return "1 minuto" // Mínimo de 1 minuto
	}
	// Arredonda para o minuto inteiro mais próximo
	return fmt.Sprintf("%.0f minutos", math.Round(etaMinutes))
}

// GetAuthenticatedUserID simula a obtenção do ID do usuário autenticado do contexto do Fiber
// Em uma aplicação real, você usaria o middleware JWT para extrair isso do token.
func GetAuthenticatedUserID(c *fiber.Ctx) uint {
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		// Logar ou retornar um erro se o user_id não estiver no contexto
		return 0 // Ou algum valor de erro apropriado
	}
	return userID
}

// GetAuthenticatedUserRole simula a obtenção da role do usuário autenticado do contexto do Fiber
func GetAuthenticatedUserRole(c *fiber.Ctx) string {
	userRole, ok := c.Locals("user_role").(string)
	if !ok {
		return ""
	}
	return userRole
}