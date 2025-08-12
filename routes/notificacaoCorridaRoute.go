package routes

import (
    "github.com/gofiber/fiber/v2"
    "taxi-service/controllers"
)

func NotificacaoCorridaRoutes(api fiber.Router) {
    notificacoes := api.Group("/notificacoes")

    // ============= ROTAS CRUD BÁSICAS =============
    // GET /notificacoes - Lista todas as notificações
    notificacoes.Get("/", controllers.ListNotificacoesCorrida)
    
    // GET /notificacoes/:id - Busca notificação por ID
    notificacoes.Get("/:id", controllers.GetNotificacaoCorrida)
    
    // POST /notificacoes - Cria nova notificação
    notificacoes.Post("/", controllers.CreateNotificacaoCorrida)
    
    // DELETE /notificacoes/:id - Remove notificação
    notificacoes.Delete("/:id", controllers.DeleteNotificacaoCorrida)

    // ============= ROTAS ESPECÍFICAS DE MOTORISTA =============
    // GET /notificacoes/motorista/:motoristaID/pending - Notificações pendentes do motorista
    notificacoes.Get("/motorista/:motoristaID/pending", controllers.GetNotificacoesPendentesParaMotorista)
    
    // GET /notificacoes/motorista/:motoristaID/historico - Histórico do motorista
    notificacoes.Get("/motorista/:motoristaID/historico", controllers.GetHistoricoNotificacoesMotorista)

    // ============= ROTAS DE AÇÕES =============
    // POST /notificacoes/:id/motorista/:motoristaID/accept - Aceitar notificação
    notificacoes.Post("/:id/motorista/:motoristaID/accept", controllers.AceitarNotificacaoCorrida)
    
    // POST /notificacoes/:id/motorista/:motoristaID/refuse - Recusar notificação
    notificacoes.Post("/:id/motorista/:motoristaID/refuse", controllers.RecusarNotificacaoCorrida)

    // ============= ROTAS ADMINISTRATIVAS =============
    // POST /notificacoes/expire - Expirar notificações vencidas
    notificacoes.Post("/expire", controllers.ExpirarNotificacoesVencidas)
    
    // PUT /notificacoes/:id/status - Atualizar status (auxiliar)
    notificacoes.Put("/:id/status", controllers.UpdateNotificacaoStatus)
}