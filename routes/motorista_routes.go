package routes

import (
	"taxi-service/controllers"
	"taxi-service/repositories"
	"taxi-service/services"

	"github.com/gofiber/fiber/v2"
)

func SetupMotoristaRoutes(api fiber.Router) {
	// Inicializar dependências
	motoristaRepo := repositories.NewJSONMotoristaRepository()
	emailService := services.NewSMTPEmailServiceFromEnv()
	motoristaService := services.NewMotoristaService(motoristaRepo, emailService)
	motoristaController := controllers.NewMotoristaController(motoristaService)

	// Grupo de rotas da API
	apiGroup := api.Group("/api")

	// Rotas de motoristas
	motoristas := apiGroup.Group("/motoristas")
	motoristas.Post("/", motoristaController.CadastrarMotorista)                      // Cadastro de motorista
	motoristas.Get("/:id", motoristaController.BuscarMotorista)                       // Buscar motorista
	motoristas.Post("/:id/documentos", motoristaController.UploadDocumento)           // Upload de documentos
	motoristas.Post("/:id/validar-documentos", motoristaController.ValidarDocumentos) // Validar documentos
	motoristas.Put("/:id/aprovar", motoristaController.AprovarMotorista)              // Aprovar motorista
	motoristas.Put("/:id/rejeitar", motoristaController.RejeitarMotorista)            // Rejeitar motorista

	// Utilitários
	motoristas.Post("/verificar-senha", motoristaController.VerificarForcaSenha)
	motoristas.Post("/validar-documento", motoristaController.ValidarDocumentoUpload)
}
