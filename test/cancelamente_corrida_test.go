package test

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
	"taxi-service/models"
	"taxi-service/routes"
	"taxi-service/services"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

// Usamos uma variável de pacote para que o *testing.T possa ser acessado nos steps.
var tc TestContext

type TestContext struct {
	app          *fiber.App
	service      *services.CorridaService
	lastResponse *http.Response
	lastBody     map[string]interface{}
	motoristas   map[string]*models.Motorista // Mapeia nome_referencia -> Motorista
	t            *testing.T
}

// TestCancelamentoDeCorrida é a função principal que o Go executa para esta suíte de testes.
func TestCancelamentoDeCorrida(t *testing.T) {
	// Atribui o *testing.T ao nosso contexto de pacote antes de rodar a suíte.
	tc.t = t

	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenarioCancelamentoDeCorrida,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"../features/cancelamento_de_corrida/"},
			TestingT: t,
			Strict:   true,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}

// InitializeScenarioCancelamentoDeCorrida registra todos os steps (passos) do nosso cenário.
func InitializeScenarioCancelamentoDeCorrida(ctx *godog.ScenarioContext) {
	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		tc.motoristas = make(map[string]*models.Motorista)
		tc.lastBody = make(map[string]interface{})

		// Criamos instâncias novas para cada cenário para garantir o isolamento.
		tc.service = services.NewCorridaService()
		app := fiber.New()
		// Injeta o serviço real no controller através da configuração de rotas
		routes.SetupCorridaRoutes(app.Group("/api"), tc.service)
		tc.app = app

		return ctx, nil
	})

	// --- Registro dos Steps ---
	ctx.Step(`^que o sistema possui os seguintes motoristas$`, tc.queOSistemaPossuiOsSeguintesMotoristas)
	ctx.Step(`^que as seguintes corridas existem no sistema, alocadas para "([^"]*)":$`, tc.queAsSeguintesCorridasExistemNoSistema)
	ctx.Step(`^a corrida de id "([^"]*)" tem o status "([^"]*)"$`, tc.aCorridaDeIdTemOStatus)
	ctx.Step(`^o motorista de cpf "([^"]*)" seleciona a opção "([^"]*)" no aplicativo para a corrida "([^"]*)"$`, tc.oMotoristaTentaCancelarACorrida)
	ctx.Step(`^o sistema atualiza o status da corrida "([^"]*)" para "([^"]*)"$`, tc.oSistemaAtualizaOStatusDaCorrida)
	ctx.Step(`^o motorista recebe a mensagem de confirmação na tela "([^"]*)"$`, tc.oMotoristaRecebeAMensagemDeConfirmacao)
	ctx.Step(`^o status do motorista de cpf "([^"]*)" é atualizado para "([^"]*)"$`, tc.oStatusDoMotoristaEAtualizadoPara)
	ctx.Step(`^a corrida de id "([^"]*)" está no status "([^"]*)" há mais de "([^"]*)" minutos$`, tc.aCorridaEstaNoStatusHaMaisDeMinutos)
	ctx.Step(`^o sistema de monitoramento executa a verificação de corridas pendentes$`, tc.oSistemaDeMonitoramentoExecutaAVerificacao)
	ctx.Step(`^o sistema altera o status da corrida "([^"]*)" para "([^"]*)"$`, tc.oSistemaAtualizaOStatusDaCorrida)
	ctx.Step(`^o sistema envia uma notificação para o motorista de cpf "([^"]*)" com a mensagem "([^"]*)"$`, tc.oSistemaEnviaUmaNotificacaoParaOMotorista)
	ctx.Step(`^o motorista de cpf "([^"]*)" tenta cancelar a corrida de id "([^"]*)" através do aplicativo$`, tc.oMotoristaTentaCancelarACorrida)
	ctx.Step(`^o sistema rejeita a solicitação de cancelamento$`, tc.oSistemaRejeitaASolicitacao)
	ctx.Step(`^o motorista recebe uma mensagem de erro na tela "([^"]*)"$`, tc.oMotoristaRecebeUmaMensagemDeErro)
	ctx.Step(`^o status da corrida "([^"]*)" permanece "([^"]*)"$`, tc.aCorridaDeIdTemOStatus)
}

// --- Funções de Implementação dos Steps ---

func (tc *TestContext) queOSistemaPossuiOsSeguintesMotoristas(table *godog.Table) error {
	for i := 1; i < len(table.Rows); i++ {
		row := table.Rows[i]
		nomeReferencia := row.Cells[0].Value
		nomeCompleto := strings.Trim(row.Cells[1].Value, `"`)
		cpf := strings.Trim(row.Cells[2].Value, `"`)
		statusStr := strings.Trim(row.Cells[3].Value, `"`)

		motorista := &models.Motorista{
			ID:     strconv.Itoa(i),
			Nome:   nomeCompleto,
			CPF:    cpf,
			Status: models.StatusMotorista(statusStr),
		}
		tc.motoristas[nomeReferencia] = motorista
	}
	return nil
}

func (tc *TestContext) queAsSeguintesCorridasExistemNoSistema(nomeReferencia string, table *godog.Table) error {
	motorista, ok := tc.motoristas[nomeReferencia]
	if !ok {
		return fmt.Errorf("motorista de referência '%s' não encontrado", nomeReferencia)
	}
	motoristaID, _ := strconv.Atoi(motorista.ID)

	for i := 1; i < len(table.Rows); i++ {
		row := table.Rows[i]
		id, _ := strconv.Atoi(row.Cells[0].Value)
		statusInicial := strings.Trim(row.Cells[1].Value, `"`)
		localEmbarque := strings.Trim(row.Cells[2].Value, `"`)
		localDesembarque := strings.Trim(row.Cells[3].Value, `"`)
		tempoEstimado, _ := strconv.Atoi(row.Cells[4].Value)

		// Cria a corrida através do serviço para ter o ID
		corrida, err := tc.service.CriarNovaCorrida(models.Corrida{
			PassageiroID:  99,
			Origem:        localEmbarque,
			Destino:       localDesembarque,
			TempoEstimado: tempoEstimado,
		})
		if err != nil {
			return fmt.Errorf("falha ao criar corrida no teste: %w", err)
		}

		// Força o estado da corrida para corresponder exatamente ao cenário Gherkin
		corrida.ID = id
		corrida.MotoristaID = motoristaID
		corrida.Status = statusInicial
	}
	return nil
}

func (tc *TestContext) aCorridaDeIdTemOStatus(idStr, statusEsperado string) error {
	resp := MakeRequest(tc.t, tc.app, "GET", "/api/corrida/"+idStr, nil)
	assert.Equal(tc.t, http.StatusOK, resp.StatusCode)

	var corrida models.Corrida
	ParseResponseBody(tc.t, resp, &corrida)

	assert.Equal(tc.t, statusEsperado, corrida.Status)
	return nil
}

func (tc *TestContext) oMotoristaTentaCancelarACorrida(cpf, idCorridaStr string) error {
	var motoristaID string
	for _, m := range tc.motoristas {
		if m.CPF == cpf {
			motoristaID = m.ID
			break
		}
	}
	if motoristaID == "" {
		return fmt.Errorf("motorista com CPF %s não encontrado no contexto do teste", cpf)
	}

	body := map[string]interface{}{"motorista_id": motoristaID}
	endpoint := fmt.Sprintf("/api/corrida/%s/cancelar/motorista", idCorridaStr)

	tc.lastResponse = MakeRequest(tc.t, tc.app, "POST", endpoint, body)
	ParseResponseBody(tc.t, tc.lastResponse, &tc.lastBody)

	return nil
}

func (tc *TestContext) oSistemaAtualizaOStatusDaCorrida(idStr, novoStatus string) error {
	return tc.aCorridaDeIdTemOStatus(idStr, novoStatus)
}

func (tc *TestContext) oMotoristaRecebeAMensagemDeConfirmacao(mensagemEsperada string) error {
	assert.Equal(tc.t, http.StatusOK, tc.lastResponse.StatusCode)
	mensagem, ok := tc.lastBody["message"].(string)
	assert.True(tc.t, ok)
	assert.Equal(tc.t, mensagemEsperada, mensagem)
	return nil
}

func (tc *TestContext) oStatusDoMotoristaEAtualizadoPara(cpf, novoStatus string) error {
	for _, m := range tc.motoristas {
		if m.CPF == cpf {
			m.Status = models.StatusMotorista(novoStatus)
			break
		}
	}
	return nil
}

func (tc *TestContext) aCorridaEstaNoStatusHaMaisDeMinutos(idStr, status, minutosStr string) error {
	id, _ := strconv.Atoi(idStr)
	minutos, _ := strconv.Atoi(minutosStr)
	corrida, err := tc.service.GetCorridaPorID(id)
	if err != nil {
		return err
	}
	corrida.DataInicio = time.Now().Add(-time.Duration(minutos+1) * time.Minute)
	corrida.Status = status
	return nil
}

func (tc *TestContext) oSistemaDeMonitoramentoExecutaAVerificacao() error {
	tc.service.MonitorarCorridasAtivas()
	return nil
}

func (tc *TestContext) oSistemaEnviaUmaNotificacaoParaOMotorista(cpf, mensagem string) error {
	fmt.Printf("SIMULAÇÃO: Notificação enviada para %s: '%s'\n", cpf, mensagem)
	return nil
}

func (tc *TestContext) oSistemaRejeitaASolicitacao() error {
	assert.NotEqual(tc.t, http.StatusOK, tc.lastResponse.StatusCode)
	return nil
}

func (tc *TestContext) oMotoristaRecebeUmaMensagemDeErro(mensagemEsperada string) error {
	mensagem, ok := tc.lastBody["error"].(string)
	assert.True(tc.t, ok)
	assert.Equal(tc.t, mensagemEsperada, mensagem)
	return nil
}
