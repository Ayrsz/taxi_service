# **Relatório Detalhado do Projeto de Serviço de Táxi**  
*(Foco em Motoristas e Corridas)*  

## **1. Visão Geral**  
Este documento detalha o escopo do projeto de desenvolvimento de um sistema de serviço de táxi, semelhante ao Uber ou 99pop, com foco nas funcionalidades de **Motorista** e **Corridas**, conforme solicitado pelo stakeholder.  

---

## **2. Escopo do Projeto**  
### **Funcionalidades Prioritárias**  
| Módulo          | Funcionalidades Implementadas |  
|----------------|-------------------------------|  
| **Motorista**  | Cadastro, autenticação, histórico de corridas, avaliação, notificações |  
| **Corridas**   | Geração automática, cancelamento, notificações, tempo estimado, histórico |  

*Obs:* Funcionalidades de **Pagamento** e **Clientes** não serão implementadas nesta fase.  

---

## **3. Detalhamento das Funcionalidades**  

### **3.1. Módulo Motorista**  
#### **Cadastro e Autenticação**  
- **Autocadastro**: O motorista pode se cadastrar fornecendo:  
  - Nome, idade (>18 anos), CPF (único), CNH, CRLV, e-mail e senha.  
- **Login**: Acesso seguro com e-mail/senha.  
- **Recuperação de conta**: Fluxo de "Esqueci a senha" via e-mail.  
- **Restrições**:  
  - Um motorista só pode editar/excluir o próprio cadastro.  
  - Um carro só pode ser vinculado a um único motorista.  

#### **Notificações**  
- **Novas corridas próximas**: Alertas em tempo real quando uma corrida estiver disponível na região.  
- **Configuração**: O motorista pode habilitar/desabilitar notificações.  

#### **Histórico de Corridas**  
- Lista todas as corridas realizadas, incluindo:  
  - Data/hora, valor, distância, status (concluída/cancelada).  
- **Avaliação**:  
  - Opção de avaliar corridas (1 a 5 estrelas + comentário opcional).  
  - Disponível no momento da conclusão ou posteriormente no histórico.  

---

### **3.2. Módulo Corridas**  
#### **Geração Automática de Corridas**  
- **Mock de dados**: JSON pré-definido com 20-30 corridas (origem, destino, valor, distância).  
- **Opcional**: Integração com API do Google Maps para cálculo de rotas (não obrigatório).  

#### **Fluxo da Corrida**  
1. **Aceitação**:  
   - Motorista recebe notificação e pode aceitar a corrida.  
2. **Tempo Estimado**:  
   - Exibição do tempo de chegada ao cliente (baseado em mock ou API).  
3. **Cancelamento**:  
   - Confirmação com mensagem de incentivo ("Tem certeza que deseja cancelar?").  
   - Registro no histórico como "cancelada".  
4. **Conclusão**:  
   - Notificação automática ao chegar ao destino.  

#### **Registro no Histórico**  
- Todas as corridas (concluídas/canceladas) são armazenadas com:  
  - Data/hora, valor, distância, status e avaliação (se aplicável).  

---

## **4. Requisitos Não Funcionais**  
- **Segurança**:  
  - Autenticação obrigatória para acessar funcionalidades.  
  - Validação de CPF e CNH para evitar cadastros duplicados.  
- **Usabilidade**:  
  - Interface intuitiva para aceitar/recusar/cancelar corridas.  
  - Feedback claro em notificações.  

---

## **5. Entregáveis**  
1. **Backend**:  
   - API para cadastro/login de motoristas.  
   - Gerenciamento de corridas (aceitar, cancelar, concluir).  
   - Mock de corridas em JSON.  
2. **Frontend (Mobile/Web)**:  
   - Telas de cadastro, login, histórico e notificações.  
   - Fluxo de aceitação/cancelamento de corridas.  

---

## **6. Próximos Passos**  
- Validação do mock de corridas com o stakeholder.  
- Definição de UI/UX para telas de notificação e histórico.  
- Testes de integração (notificações, cancelamento, avaliação).  

**Observação**: As funcionalidades de pagamento e clientes ficarão para uma futura fase do projeto.  

