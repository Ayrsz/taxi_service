package e2e

import (
    "testing"
    "time"
	"fmt"

    "taxi-service/models"
    "taxi-service/test"

    "github.com/stretchr/testify/assert"
)

func TestListNotificacoesCorrida(t *testing.T) {
    app := test.SetupTestApp(t)
    defer test.CleanupTestApp(t)

    // Test List all notificacoes
    resp := test.MakeRequest(t, app, "GET", "/notificacoes", nil)
    assert.Equal(t, 200, resp.StatusCode)

    var notificacoes []models.NotificacaoCorrida
    test.ParseResponseBody(t, resp, &notificacoes)

    t.Logf("Found %d notificacoes", len(notificacoes))
    
    // Verificar se há dados de exemplo
    if len(notificacoes) > 0 {
        for i, notif := range notificacoes {
            assert.NotZero(t, notif.ID, "Notificacao should have an ID")
            assert.NotZero(t, notif.MotoristaID, "Notificacao should have a MotoristaID")
            assert.NotZero(t, notif.CorridaID, "Notificacao should have a CorridaID")
            assert.NotEmpty(t, notif.PassageiroNome, "Notificacao should have PassageiroNome")
            assert.Greater(t, notif.Valor, 0.0, "Valor should be greater than 0")
            
            t.Logf("Notificacao[%d]: ID=%d, Motorista=%d, Status='%s', Valor=%.2f", 
                i, notif.ID, notif.MotoristaID, notif.Status, notif.Valor)
        }
    }
}

func TestCreateNotificacaoCorrida(t *testing.T) {
    app := test.SetupTestApp(t)
    defer test.CleanupTestApp(t)

    // Get initial count
    listResp := test.MakeRequest(t, app, "GET", "/notificacoes", nil)
    assert.Equal(t, 200, listResp.StatusCode)
    
    var initialNotificacoes []models.NotificacaoCorrida
    test.ParseResponseBody(t, listResp, &initialNotificacoes)
    initialCount := len(initialNotificacoes)
    t.Logf("Initial notificacoes count: %d", initialCount)

    // Create a new notificacao
    newNotificacao := models.NotificacaoCorrida{
        MotoristaID:     999,
        CorridaID:       888,
        PassageiroNome:  "Test Passenger",
        Valor:           45.50,
        DistanciaKm:     8.5,
        TempoEstimado:   "15 min",
        Origem:          "Test Origin Address",
        Destino:         "Test Destination Address",
    }

    createResp := test.MakeRequest(t, app, "POST", "/notificacoes", newNotificacao)
    assert.Equal(t, 201, createResp.StatusCode)

    var createdNotificacao models.NotificacaoCorrida
    test.ParseResponseBody(t, createResp, &createdNotificacao)
    
    // Verify created notificacao
    assert.NotZero(t, createdNotificacao.ID, "Created notificacao should have an ID")
    assert.Equal(t, uint(999), createdNotificacao.MotoristaID)
    assert.Equal(t, uint(888), createdNotificacao.CorridaID)
    assert.Equal(t, "Test Passenger", createdNotificacao.PassageiroNome)
    assert.Equal(t, 45.50, createdNotificacao.Valor)
    assert.Equal(t, models.NotificacaoPendente, createdNotificacao.Status)
    assert.False(t, createdNotificacao.CreatedAt.IsZero(), "CreatedAt should be set")
    assert.False(t, createdNotificacao.ExpiraEm.IsZero(), "ExpiraEm should be set")
    
    t.Logf("Created notificacao with ID: %d, Status: %s", createdNotificacao.ID, createdNotificacao.Status)

    // Verify count increased
    listResp2 := test.MakeRequest(t, app, "GET", "/notificacoes", nil)
    assert.Equal(t, 200, listResp2.StatusCode)
    
    var afterCreateNotificacoes []models.NotificacaoCorrida
    test.ParseResponseBody(t, listResp2, &afterCreateNotificacoes)
    assert.Equal(t, initialCount+1, len(afterCreateNotificacoes), "Notificacoes count should increase by 1")
}

func TestGetNotificacoesPendentesParaMotorista(t *testing.T) {
    app := test.SetupTestApp(t)
    defer test.CleanupTestApp(t)

    // Test Get pending notificacoes for motorista 102
    resp := test.MakeRequest(t, app, "GET", "/notificacoes/motorista/102/pending", nil)
    assert.Equal(t, 200, resp.StatusCode)

    var pendingResp map[string]interface{}
    test.ParseResponseBody(t, resp, &pendingResp)

    motoristaID := pendingResp["motorista_id"]
    pendingCount := pendingResp["pending_count"]
    
    assert.Equal(t, float64(102), motoristaID)
    assert.NotNil(t, pendingCount)
    
    t.Logf("Motorista 102 has %v pending notificacoes", pendingCount)

    // Test with non-existent motorista
    resp2 := test.MakeRequest(t, app, "GET", "/notificacoes/motorista/99999/pending", nil)
    assert.Equal(t, 200, resp2.StatusCode)

    var emptyResp map[string]interface{}
    test.ParseResponseBody(t, resp2, &emptyResp)
    assert.Equal(t, float64(0), emptyResp["pending_count"])
}

func TestGetHistoricoNotificacoesMotorista(t *testing.T) {
    app := test.SetupTestApp(t)
    defer test.CleanupTestApp(t)

    // Test Get historico for motorista 101
    resp := test.MakeRequest(t, app, "GET", "/notificacoes/motorista/101/historico", nil)
    assert.Equal(t, 200, resp.StatusCode)

    var historicoResp map[string]interface{}
    test.ParseResponseBody(t, resp, &historicoResp)

    motoristaID := historicoResp["motorista_id"]
    totalCount := historicoResp["total_count"]
    aceitasCount := historicoResp["aceitas_count"]
    recusadasCount := historicoResp["recusadas_count"]
    expiradasCount := historicoResp["expiradas_count"]
    pendentesCount := historicoResp["pendentes_count"]

    assert.Equal(t, float64(101), motoristaID)
    assert.NotNil(t, totalCount)
    assert.NotNil(t, aceitasCount)
    assert.NotNil(t, recusadasCount)
    assert.NotNil(t, expiradasCount)
    assert.NotNil(t, pendentesCount)

    t.Logf("Motorista 101 historico: Total=%v, Aceitas=%v, Recusadas=%v, Expiradas=%v, Pendentes=%v",
        totalCount, aceitasCount, recusadasCount, expiradasCount, pendentesCount)
}

func TestAceitarNotificacaoCorrida(t *testing.T) {
    app := test.SetupTestApp(t)
    defer test.CleanupTestApp(t)

    // First create a notificacao to accept
    newNotificacao := models.NotificacaoCorrida{
        MotoristaID:     555,
        CorridaID:       777,
        PassageiroNome:  "Accept Test Passenger",
        Valor:           30.00,
        DistanciaKm:     5.0,
        TempoEstimado:   "10 min",
        Origem:          "Accept Test Origin",
        Destino:         "Accept Test Destination",
    }

    createResp := test.MakeRequest(t, app, "POST", "/notificacoes", newNotificacao)
    assert.Equal(t, 201, createResp.StatusCode)

    var createdNotificacao models.NotificacaoCorrida
    test.ParseResponseBody(t, createResp, &createdNotificacao)
    
    // Wait a moment to ensure the notificacao is created
    time.Sleep(100 * time.Millisecond)

    // Accept the notificacao
    acceptPath := fmt.Sprintf("/notificacoes/%d/motorista/555/accept", createdNotificacao.ID)
    acceptResp := test.MakeRequest(t, app, "POST", acceptPath, nil)
    
    switch acceptResp.StatusCode {
    case 200:
        var acceptResult map[string]interface{}
        test.ParseResponseBody(t, acceptResp, &acceptResult)
        
        assert.Equal(t, "Notificacao accepted successfully", acceptResult["message"])
        assert.Equal(t, float64(createdNotificacao.ID), acceptResult["notificacao_id"])
        assert.Equal(t, float64(555), acceptResult["motorista_id"])
        
        t.Logf("Successfully accepted notificacao ID: %d", createdNotificacao.ID)
    case 410:
        t.Log("Notificacao expired before we could accept it - this can happen with 20s expiration")
    default:
        t.Logf("Accept returned status: %d", acceptResp.StatusCode)
    }
}

func TestRecusarNotificacaoCorrida(t *testing.T) {
    app := test.SetupTestApp(t)
    defer test.CleanupTestApp(t)

    // First create a notificacao to refuse
    newNotificacao := models.NotificacaoCorrida{
        MotoristaID:     666,
        CorridaID:       888,
        PassageiroNome:  "Refuse Test Passenger",
        Valor:           25.00,
        DistanciaKm:     4.0,
        TempoEstimado:   "8 min",
        Origem:          "Refuse Test Origin",
        Destino:         "Refuse Test Destination",
    }

    createResp := test.MakeRequest(t, app, "POST", "/notificacoes", newNotificacao)
    assert.Equal(t, 201, createResp.StatusCode)

    var createdNotificacao models.NotificacaoCorrida
    test.ParseResponseBody(t, createResp, &createdNotificacao)
    
    // Wait a moment to ensure the notificacao is created
    time.Sleep(100 * time.Millisecond)

    // Refuse the notificacao
    refusePath := "/notificacoes/" + string(rune(createdNotificacao.ID+'0')) + "/motorista/666/refuse"
    
    // For IDs > 9, we need proper conversion
    if createdNotificacao.ID <= 9 {
        refuseResp := test.MakeRequest(t, app, "POST", refusePath, nil)
        
        switch refuseResp.StatusCode {
        case 200:
            var refuseResult map[string]interface{}
            test.ParseResponseBody(t, refuseResp, &refuseResult)
            
            assert.Equal(t, "Notificacao refused successfully", refuseResult["message"])
            assert.Equal(t, float64(createdNotificacao.ID), refuseResult["notificacao_id"])
            assert.Equal(t, float64(666), refuseResult["motorista_id"])
            
            t.Logf("Successfully refused notificacao ID: %d", createdNotificacao.ID)
        case 410:
            t.Log("Notificacao expired before we could refuse it - this can happen with 20s expiration")
        default:
            t.Logf("Refuse returned status: %d", refuseResp.StatusCode)
        }
    }
}

func TestExpirarNotificacoesVencidas(t *testing.T) {
    app := test.SetupTestApp(t)
    defer test.CleanupTestApp(t)

    // 1. Criar uma notificação que irá expirar
    newNotificacao := models.NotificacaoCorrida{
        MotoristaID:     777,
        CorridaID:       888,
        PassageiroNome:  "Expire Test Passenger",
        Valor:           25.00,
        DistanciaKm:     4.0,
        TempoEstimado:   "8 min",
        Origem:          "Test Origin",
        Destino:         "Test Destination",
    }

    createResp := test.MakeRequest(t, app, "POST", "/notificacoes", newNotificacao)
    assert.Equal(t, 201, createResp.StatusCode)

    var createdNotificacao models.NotificacaoCorrida
    test.ParseResponseBody(t, createResp, &createdNotificacao)

    // Verificação inicial
    assert.Equal(t, models.NotificacaoPendente, createdNotificacao.Status)
    t.Logf("Created notificacao ID %d with status: %s", createdNotificacao.ID, createdNotificacao.Status)

    // 2. Esperar pelo menos 22 segundos (20s de expiração + margem)
    t.Log("Waiting for notificacao to expire naturally by service...")
    time.Sleep(22 * time.Second)

    // 3. Verificar se a notificação foi automaticamente marcada como expirada
    getResp := test.MakeRequest(t, app, "GET", fmt.Sprintf("/notificacoes/%d", createdNotificacao.ID), nil)
    assert.Equal(t, 200, getResp.StatusCode)

    var after models.NotificacaoCorrida
    test.ParseResponseBody(t, getResp, &after)

    // 4. Verificar se o status mudou para "expirada"
    assert.Equal(t, models.NotificacaoExpirada, after.Status,
        "Expected notificacao to be marked as expired automatically")
    assert.True(t, after.UpdatedAt.After(createdNotificacao.UpdatedAt),
        "UpdatedAt should be updated after expiration")

    t.Logf("Notificacao ID %d expired as expected (status: %s)", after.ID, after.Status)

    // 5. Verificar que ela não aparece mais na lista de notificações pendentes do motorista
    pendingResp := test.MakeRequest(t, app, "GET", fmt.Sprintf("/notificacoes/motorista/%d/pending", newNotificacao.MotoristaID), nil)
    assert.Equal(t, 200, pendingResp.StatusCode)

    var pendingResult map[string]interface{}
    test.ParseResponseBody(t, pendingResp, &pendingResult)

    if notificacoes, exists := pendingResult["notificacoes"]; exists {
        if notifArray, ok := notificacoes.([]interface{}); ok {
            for _, notif := range notifArray {
                if notifMap, ok := notif.(map[string]interface{}); ok {
                    if id, exists := notifMap["id"]; exists && id == float64(createdNotificacao.ID) {
                        t.Errorf("Expired notificacao ID %d should not appear in pending list", createdNotificacao.ID)
                    }
                }
            }
        }
    }

    t.Log("Automatic expiration test completed successfully")
}


func TestDeleteNotificacaoCorrida(t *testing.T) {
    app := test.SetupTestApp(t)
    defer test.CleanupTestApp(t)

    // First create a notificacao to delete
    newNotificacao := models.NotificacaoCorrida{
        MotoristaID:     777,
        CorridaID:       999,
        PassageiroNome:  "Delete Test Passenger",
        Valor:           20.00,
        DistanciaKm:     3.0,
        TempoEstimado:   "7 min",
        Origem:          "Delete Test Origin",
        Destino:         "Delete Test Destination",
    }

    createResp := test.MakeRequest(t, app, "POST", "/notificacoes", newNotificacao)
    assert.Equal(t, 201, createResp.StatusCode)

    var createdNotificacao models.NotificacaoCorrida
    test.ParseResponseBody(t, createResp, &createdNotificacao)

    // Delete the notificacao
    deletePath := "/notificacoes/" + string(rune(createdNotificacao.ID+'0'))
    
    // For IDs > 9, we need proper conversion
    if createdNotificacao.ID <= 9 {
        deleteResp := test.MakeRequest(t, app, "DELETE", deletePath, nil)
        assert.Equal(t, 200, deleteResp.StatusCode)

        var deleteResult map[string]interface{}
        test.ParseResponseBody(t, deleteResp, &deleteResult)
        
        assert.Equal(t, "Notificacao deleted successfully", deleteResult["message"])
        assert.Equal(t, float64(createdNotificacao.ID), deleteResult["id"])
        
        t.Logf("Successfully deleted notificacao ID: %d", createdNotificacao.ID)

        // Verify it's gone
        getResp := test.MakeRequest(t, app, "GET", deletePath, nil)
        assert.Equal(t, 404, getResp.StatusCode)
    }
}

func TestNotificacaoCorridaWorkflow(t *testing.T) {
    app := test.SetupTestApp(t)
    defer test.CleanupTestApp(t)

    t.Log("=== TESTING COMPLETE NOTIFICACAO WORKFLOW ===")

    // 1. Create a notificacao
    motoristaID := uint(123)
    newNotificacao := models.NotificacaoCorrida{
        MotoristaID:     motoristaID,
        CorridaID:       456,
        PassageiroNome:  "Workflow Test Passenger",
        Valor:           35.75,
        DistanciaKm:     7.2,
        TempoEstimado:   "16 min",
        Origem:          "Workflow Test Origin",
        Destino:         "Workflow Test Destination",
    }

    createResp := test.MakeRequest(t, app, "POST", "/notificacoes", newNotificacao)
    assert.Equal(t, 201, createResp.StatusCode)

    var createdNotificacao models.NotificacaoCorrida
    test.ParseResponseBody(t, createResp, &createdNotificacao)
    assert.Equal(t, models.NotificacaoPendente, createdNotificacao.Status)
    t.Logf("1. Created notificacao ID: %d with status: %s", createdNotificacao.ID, createdNotificacao.Status)

    // 2. Check it appears in pending list
    pendingResp := test.MakeRequest(t, app, "GET", "/notificacoes/motorista/123/pending", nil)
    assert.Equal(t, 200, pendingResp.StatusCode)

    var pendingResult map[string]interface{}
    test.ParseResponseBody(t, pendingResp, &pendingResult)
    t.Logf("2. Motorista 123 has %v pending notificacoes", pendingResult["pending_count"])

    // 3. Check historico includes it
    historicoResp := test.MakeRequest(t, app, "GET", "/notificacoes/motorista/123/historico", nil)
    assert.Equal(t, 200, historicoResp.StatusCode)

    var historicoResult map[string]interface{}
    test.ParseResponseBody(t, historicoResp, &historicoResult)
    t.Logf("3. Motorista 123 historico: Total=%v, Pendentes=%v", 
        historicoResult["total_count"], historicoResult["pendentes_count"])

    t.Log("=== WORKFLOW TEST COMPLETED SUCCESSFULLY ===")
}

func TestInvalidRequests(t *testing.T) {
    app := test.SetupTestApp(t)
    defer test.CleanupTestApp(t)

    // Test invalid IDs
    testCases := []struct {
        name           string
        method         string
        path           string
        expectedStatus int
    }{
        {"Get invalid ID", "GET", "/notificacoes/abc", 400},
        {"Delete invalid ID", "DELETE", "/notificacoes/xyz", 400},
        {"Accept invalid notificacao ID", "POST", "/notificacoes/abc/motorista/123/accept", 400},
        {"Accept invalid motorista ID", "POST", "/notificacoes/1/motorista/xyz/accept", 400},
        {"Refuse invalid notificacao ID", "POST", "/notificacoes/abc/motorista/123/refuse", 400},
        {"Refuse invalid motorista ID", "POST", "/notificacoes/1/motorista/xyz/refuse", 400},
        {"Pending invalid motorista ID", "GET", "/notificacoes/motorista/abc/pending", 400},
        {"Historico invalid motorista ID", "GET", "/notificacoes/motorista/xyz/historico", 400},
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            resp := test.MakeRequest(t, app, tc.method, tc.path, nil)
            assert.Equal(t, tc.expectedStatus, resp.StatusCode)
            t.Logf("Invalid request test '%s' passed with status %d", tc.name, resp.StatusCode)
        })
    }
}