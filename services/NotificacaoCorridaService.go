package services

import (
    "encoding/json"
    "errors"
    "os"
    "path/filepath"
    "taxi-service/models"
    "time"
)

const notificacaoCorridaFile = "./data/notificacao_corrida.json"

// ============= FUNÇÕES AUXILIARES DE ARQUIVO =============

func readNotificacoesCorrida() ([]models.NotificacaoCorrida, error) {
    
    // Criar diretório se não existir
    dataDir := filepath.Dir(notificacaoCorridaFile)
    if err := os.MkdirAll(dataDir, 0755); err != nil {
        return nil, err
    }

    file, err := os.Open(notificacaoCorridaFile)
    if err != nil {
        if os.IsNotExist(err) {
            return []models.NotificacaoCorrida{}, nil
        }
        return nil, err
    }
    defer file.Close()

    var notificacoes []models.NotificacaoCorrida
    err = json.NewDecoder(file).Decode(&notificacoes)
    if err != nil {
        return nil, err
    }
    
    return notificacoes, nil
}

func writeNotificacoesCorrida(notificacoes []models.NotificacaoCorrida) error {
    
    // Garantir que o diretório existe antes de escrever
    dataDir := filepath.Dir(notificacaoCorridaFile)
    if err := os.MkdirAll(dataDir, 0755); err != nil {
        return err
    }
    
    data, err := json.MarshalIndent(notificacoes, "", "  ")
    if err != nil {
        return err
    }
    
    err = os.WriteFile(notificacaoCorridaFile, data, 0644)
    if err != nil {
        return err
    }
    
    return nil
}

// ============= FUNÇÕES PRINCIPAIS DE SERVIÇO =============

// ListNotificacoesCorrida - Lista todas as notificações
func ListNotificacoesCorrida() ([]models.NotificacaoCorrida, error) {
    return readNotificacoesCorrida()
}

// GetNotificacaoCorrida - Busca notificação por ID
func GetNotificacaoCorrida(id uint) (models.NotificacaoCorrida, error) {
    
    notificacoes, err := readNotificacoesCorrida()
    if err != nil {
        return models.NotificacaoCorrida{}, err
    }
    
    for _, notificacao := range notificacoes {
        if notificacao.ID == id {
            return notificacao, nil
        }
    }
    
    return models.NotificacaoCorrida{}, errors.New("notificacao not found")
}

// CreateNotificacaoCorrida - Cria nova notificação para motorista
func CreateNotificacaoCorrida(notificacao *models.NotificacaoCorrida) error {
    notificacoes, err := readNotificacoesCorrida()
    if err != nil {
        return err
    }

    // Atribuir novo ID (máximo + 1)
    var maxID uint = 0
    for _, n := range notificacoes {
        if n.ID > maxID {
            maxID = n.ID
        }
    }
    notificacao.ID = maxID + 1

    // Definir valores padrão
    now := time.Now()
    notificacao.Status = models.NotificacaoPendente
    notificacao.CreatedAt = now
    notificacao.UpdatedAt = now
    notificacao.ExpiraEm = now.Add(20 * time.Second) // Expira em 20 segundos

    // Adicionar nova notificação à lista
    notificacoes = append(notificacoes, *notificacao)

    if err := writeNotificacoesCorrida(notificacoes); err != nil {
        return err
    }

    // Iniciar rotina para expiração automática após 21 segundos
    go func(id uint) {
        time.Sleep(21 * time.Second) // 20s + 1s de margem

        // Ler novamente as notificações
        notificacoes, err := readNotificacoesCorrida()
        if err != nil {
            // Idealmente logar erro, mas estamos em goroutine
            return
        }

        updated := false
        now := time.Now()
        for i, n := range notificacoes {
            if n.ID == id && n.Status == models.NotificacaoPendente && now.After(n.ExpiraEm) {
                notificacoes[i].Status = models.NotificacaoExpirada
                notificacoes[i].UpdatedAt = now
                updated = true
                break
            }
        }

        if updated {
            _ = writeNotificacoesCorrida(notificacoes) // Ignorar erro por simplicidade
        }
    }(notificacao.ID)

    return nil
}

// GetNotificacoesPendentesParaMotorista - Busca notificações pendentes para um motorista específico
func GetNotificacoesPendentesParaMotorista(motoristaID uint) ([]models.NotificacaoCorrida, error) {
    
    notificacoes, err := readNotificacoesCorrida()
    if err != nil {
        return nil, err
    }
    
    var notificacoesPendentes []models.NotificacaoCorrida
    agora := time.Now()
    
    for _, notificacao := range notificacoes {
        // Verificar se é para o motorista e está pendente
        if notificacao.MotoristaID == motoristaID && notificacao.Status == models.NotificacaoPendente {
            // Verificar se não expirou
            if agora.Before(notificacao.ExpiraEm) {
                notificacoesPendentes = append(notificacoesPendentes, notificacao)
            } else {
            }
        }
    }
    
    return notificacoesPendentes, nil
}

// AceitarNotificacaoCorrida - Aceita uma notificação de corrida
func AceitarNotificacaoCorrida(notificacaoID uint, motoristaID uint) error {
    
    notificacoes, err := readNotificacoesCorrida()
    if err != nil {
        return err
    }
    
    for i, notificacao := range notificacoes {
        if notificacao.ID == notificacaoID && notificacao.MotoristaID == motoristaID {
            // Verificar se ainda está pendente e não expirou
            if notificacao.Status != models.NotificacaoPendente {
                return errors.New("notificacao already processed")
            }
            
            if time.Now().After(notificacao.ExpiraEm) {
                // Marcar como expirada
                notificacoes[i].Status = models.NotificacaoExpirada
                notificacoes[i].UpdatedAt = time.Now()
                writeNotificacoesCorrida(notificacoes)
                return errors.New("notificacao expired")
            }
            
            // Aceitar a notificação
            notificacoes[i].Status = models.NotificacaoAceita
            notificacoes[i].UpdatedAt = time.Now()
            
            if err := writeNotificacoesCorrida(notificacoes); err != nil {
                return err
            }
            
            return nil
        }
    }
    
    return errors.New("notificacao not found")
}

// RecusarNotificacaoCorrida - Recusa uma notificação de corrida
func RecusarNotificacaoCorrida(notificacaoID uint, motoristaID uint) error {
    
    notificacoes, err := readNotificacoesCorrida()
    if err != nil {
        return err
    }
    
    for i, notificacao := range notificacoes {
        if notificacao.ID == notificacaoID && notificacao.MotoristaID == motoristaID {
            // Verificar se ainda está pendente
            if notificacao.Status != models.NotificacaoPendente {
                return errors.New("notificacao already processed")
            }
            
            // Recusar a notificação
            notificacoes[i].Status = models.NotificacaoRecusada
            notificacoes[i].UpdatedAt = time.Now()
            
            if err := writeNotificacoesCorrida(notificacoes); err != nil {
                return err
            }
            
            return nil
        }
    }
    
    return errors.New("notificacao not found")
}

// ExpirarNotificacoesVencidas - Marca como expiradas as notificações que passaram do tempo limite
func ExpirarNotificacoesVencidas() error {
    
    notificacoes, err := readNotificacoesCorrida()
    if err != nil {
        return err
    }
    
    agora := time.Now()
    expiradas := 0
    
    for i, notificacao := range notificacoes {
        if notificacao.Status == models.NotificacaoPendente && agora.After(notificacao.ExpiraEm) {
            notificacoes[i].Status = models.NotificacaoExpirada
            notificacoes[i].UpdatedAt = agora
            expiradas++
        }
    }
    
    if expiradas > 0 {
        if err := writeNotificacoesCorrida(notificacoes); err != nil {
            return err
        }
    } else {
    }
    
    return nil
}

// GetHistoricoNotificacoesMotorista - Busca histórico de notificações de um motorista
func GetHistoricoNotificacoesMotorista(motoristaID uint) ([]models.NotificacaoCorrida, error) {
    
    notificacoes, err := readNotificacoesCorrida()
    if err != nil {
        return nil, err
    }
    
    var historico []models.NotificacaoCorrida
    for _, notificacao := range notificacoes {
        if notificacao.MotoristaID == motoristaID {
            historico = append(historico, notificacao)
        }
    }
    
    return historico, nil
}

// DeleteNotificacaoCorrida - Remove uma notificação (para limpeza de dados antigos)
func DeleteNotificacaoCorrida(id uint) error {
    
    notificacoes, err := readNotificacoesCorrida()
    if err != nil {
        return err
    }
    
    newNotificacoes := []models.NotificacaoCorrida{}
    found := false
    
    for _, notificacao := range notificacoes {
        if notificacao.ID != id {
            newNotificacoes = append(newNotificacoes, notificacao)
        } else {
            found = true
        }
    }
    
    if !found {
        return errors.New("notificacao not found")
    }
    
    if err := writeNotificacoesCorrida(newNotificacoes); err != nil {
        return err
    }
    
    return nil
}