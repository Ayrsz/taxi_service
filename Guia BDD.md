# **Guia Detalhado de Boas Práticas para Desenvolvimento BDD (Behavior-Driven Development)**

---

## **1. Introdução ao BDD**
O BDD (Behavior-Driven Development) é uma abordagem de desenvolvimento de software que combina técnicas de **testes automatizados** e **especificação de requisitos** por meio de cenários escritos em linguagem natural. O objetivo é alinhar as expectativas dos stakeholders com a implementação técnica, garantindo que o sistema seja construído corretamente e atenda às necessidades reais.

---

## **2. Fundamentos do BDD**
### **2.1. O que é um Teste?**
- **Objetivo principal**: Garantir qualidade e fornecer evidências de que o sistema se comporta conforme o esperado.
- **Tipos de comportamentos verificados**:
  - Funcionalidades (corretude)
  - Robustez
  - Desempenho e escalabilidade
  - Usabilidade (GUI)
  - Segurança

### **2.2. Testes de Aceitação vs. Outros Tipos de Testes**
- **Testes de Aceitação**: Validam se o sistema atende aos requisitos dos stakeholders.
- **Testes de Regressão**: Garantem que alterações não quebrem funcionalidades existentes.
- **Testes de Fumaça (Smoke Tests)**: Verificam o funcionamento básico do sistema.

---

## **3. Escrevendo Cenários BDD**
### **3.1. Estrutura de um Cenário (Gherkin)**
Um cenário BDD segue o formato **Given-When-Then**:
- **Given (Dado)**: Pré-condições (estado inicial do sistema).
- **When (Quando)**: Ação realizada pelo usuário ou sistema.
- **Then (Então)**: Resultados esperados (pós-condições).

Exemplo:
```gherkin
Scenario: Adicionar nova nota
  Given Eu estou na página "Notas"
  And O aluno "João" não tem nota para "Matemática"
  When Eu adiciono a nota "8.5" para "João" em "Matemática"
  Then A nota "8.5" é exibida para "João" em "Matemática"
```

### **3.2. Boas Práticas para Cenários**
#### ✅ **Recomendações**
1. **Seja Declarativo, Não Imperativo**  
   - ❌ Evite: "Eu clico no botão X, depois no botão Y..."  
   - ✅ Prefira: "Eu adiciono a nota '8.5' para 'João' em 'Matemática'."

2. **Use Dados Concretos**  
   - Evite generalizações como "um aluno" ou "uma nota".  
   - ✅ Exemplo: `Given O aluno "Maria Silva" está cadastrado com CPF "123456789"`.

3. **Mantenha os Cenários Curtos e Objetivos**  
   - Um cenário deve testar **um único comportamento**.

4. **Reutilize Passos**  
   - Evite duplicação usando steps compartilhados (ex.: `Given Eu estou logado como "professor"`).

5. **Especifique Resultados Esperados com Clareza**  
   - Não basta dizer "o sistema mostra uma mensagem".  
   - ✅ Exemplo: `Then Eu vejo a mensagem "Aluno cadastrado com sucesso"`.

6. **Parametrize Cenários com `Scenario Outline`**  
   - Útil para testar múltiplos casos com a mesma estrutura.  
   ```gherkin
   Scenario Outline: Login com credenciais válidas
     Given Eu estou na página de login
     When Eu insiro "<email>" e "<senha>"
     Then Eu sou redirecionado para a página inicial
     Examples:
       | email           | senha   |
       | user1@test.com  | 123456  |
       | admin@test.com  | admin   |
   ```

7. **Inclua Cenários de Falha**  
   - Teste casos onde o sistema deve **rejeitar entradas inválidas**.  
   ```gherkin
   Scenario: Tentativa de login com senha incorreta
     Given Eu estou na página de login
     When Eu insiro "user1@test.com" e "senha_errada"
     Then Eu vejo a mensagem "Credenciais inválidas"
   ```

#### ❌ **Críticas e Problemas Comuns**
1. **Cenários Muito Longos ou Complexos**  
   - Se um cenário tem muitos `And`, pode estar testando múltiplos comportamentos (violando o princípio da responsabilidade única).

2. **Passos Vagos ou Incompletos**  
   - ❌ Exemplo: `When Eu faço alguma coisa` → **Sem valor para testes**.

3. **Dependência Entre Cenários**  
   - Cada cenário deve ser **independente** (não depender de estado deixado por outro).

4. **Falta de Dados Concretos**  
   - ❌ Exemplo: `Given Um aluno existe` → **Qual aluno? Quais dados?**

5. **Cenários Imperativos (Foco em GUI)**  
   - ❌ Exemplo: `When Eu clico no botão "Salvar"` → **Muito sensível a mudanças na interface**.

6. **Redundância**  
   - Evite cenários que só mudam parâmetros sem acrescentar novos comportamentos.

---

## **4. Tipos de Testes no BDD**
### **4.1. Testes GUI vs. Testes de Serviço**
| **Critério**       | **Testes GUI**                          | **Testes de Serviço**                     |
|--------------------|----------------------------------------|-------------------------------------------|
| **Foco**           | Interface do usuário (navegador)       | API/Backend (requisições HTTP)            |
| **Vantagens**      | Testa fluxo completo do usuário        | Mais rápidos e estáveis                   |
| **Desvantagens**   | Frágeis (mudanças na UI quebram testes | Não cobrem interação visual               |

### **4.2. Pirâmide de Testes**
```
        GUI
      /      \
  Service   Service
 /             \
Unit           Unit
```
- **Base (Unit)**: Testes de unidade (métodos/classes).  
- **Meio (Service)**: Testes de API/integração.  
- **Topo (GUI)**: Testes de interface (menos frequentes).

---

## **5. Implementação e Manutenção de Testes**
### **5.1. Automatização**
- Use ferramentas como:
  - **Cucumber** (para cenários Gherkin).
  - **Selenium** (para testes GUI).
  - **RestAssured** (para testes de API).

### **5.2. Boas Práticas de Código de Teste**
1. **Mantenha os Testes Limpos**  
   - Evite código duplicado (use funções auxiliares).
2. **Testes Independentes**  
   - Cada teste deve rodar em um ambiente limpo.
3. **Lidere com Dados Falsos (Mocks/Stubs)**  
   - Para testes de unidade, simule dependências externas.
4. **Verifique Resultados com Precisão**  
   - Não apenas "o sistema não quebrou", mas "o sistema fez **exatamente** o esperado".

### **5.3. Evite Armadilhas Comuns**
- **Problema do Oráculo**: Testes que passam sem verificar o resultado real.
- **Testes Frágeis**: Dependência excessiva de detalhes de implementação.
- **Falta de Cobertura**: Cenários importantes não testados.

---

## **6. Integração com Desenvolvimento**
### **6.1. BDD no Fluxo de Desenvolvimento**
1. **Escreva Cenários Antes do Código** (BDD puro).
2. **Use Cenários como Especificação Viva**.
3. **Execute Testes Automaticamente no CI/CD**.

### **6.2. TDD vs. BDD**
| **Critério**       | **TDD**                          | **BDD**                          |
|--------------------|----------------------------------|----------------------------------|
| **Foco**           | Testes unitários (código)        | Comportamento (requisitos)       |
| **Linguagem**      | Código (ex.: JUnit)              | Gherkin (Given-When-Then)        |
| **Público**        | Desenvolvedores                  | Stakeholders + Devs + QA         |

---

## **7. Checklist de Revisão de Cenários**
1. [ ] O cenário é **SMART** (Específico, Mensurável, etc.)?
2. [ ] Os passos são **declarativos** (não imperativos)?
3. [ ] Os dados são **concretos** (não genéricos)?
4. [ ] O cenário testa **apenas um comportamento**?
5. [ ] Há **cenários de sucesso e falha**?
6. [ ] Os resultados são **verificáveis**?
7. [ ] O cenário é **independente** de outros?

---

## **8. Conclusão**
O BDD é uma **ferramenta poderosa** para alinhar requisitos, testes e implementação. Cenários bem escritos:
- **Melhoram a comunicação** entre stakeholders e devs.
- **Reduzem ambiguidades** nos requisitos.
- **Facilitam a manutenção** dos testes.

Seguindo essas práticas, sua equipe pode desenvolver software **mais confiável, testável e alinhado com as necessidades reais**.