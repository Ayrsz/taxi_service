# language: pt
Funcionalidade: Cadastro de Motorista
  Como um futuro motorista
  Eu quero me cadastrar no sistema
  Para poder aceitar corridas e trabalhar como motorista

  Contexto:
    Dado que o sistema de cadastro está funcionando
    E o serviço de validação de documentos está disponível

  Cenário: Cadastro bem-sucedido com todos os dados válidos
    Quando realizo o cadastro com dados válidos:
      | nome           | João Silva                    |
      | data_nascimento| 15/03/1990                   |
      | cpf            | 123.456.789-00               |
      | cnh            | 12345678901                  |
      | categoria_cnh  | B                            |
      | validade_cnh   | 15/03/2030                   |
      | placa_veiculo  | ABC1234                      |
      | modelo_veiculo | Honda Civic 2020             |
      | telefone       | (11) 99999-9999              |
      | email          | joao.silva@email.com         |
      | senha          | MinhaSenh@123                |
    Então vejo a mensagem "Cadastro realizado com sucesso"
    E recebo um email de confirmação
    E meu status inicial é "aguardando_aprovacao"

  Esquema do Cenário: Validação de campos obrigatórios
    Quando tento realizar cadastro sem preencher o campo "<campo_faltante>"
    Então vejo a mensagem de erro "<mensagem_erro>"
    E não consigo concluir o cadastro

    Exemplos:
      | campo_faltante  | mensagem_erro                       |
      | nome            | Nome é obrigatório                  |
      | cpf             | CPF é obrigatório                   |
      | cnh             | CNH é obrigatória                   |
      | email           | Email é obrigatório                 |
      | senha           | Senha é obrigatória                 |
      | telefone        | Telefone é obrigatório              |
      | placa_veiculo   | Placa do veículo é obrigatória      |

  Esquema do Cenário: Validação de formatos de documentos
    Quando preencho o campo "<campo>" com o valor "<valor_invalido>"
    Então vejo a mensagem de erro "<mensagem_erro>"
    E não consigo concluir o cadastro

    Exemplos:
      | campo    | valor_invalido  | mensagem_erro                |
      | cpf      | 123.456.789-99  | CPF inválido                 |
      | cpf      | 111.111.111-11  | CPF inválido                 |
      | cnh      | 1234567890      | CNH deve ter 11 dígitos      |
      | email    | email_invalido  | Formato de email inválido    |
      | telefone | 999999999       | Formato de telefone inválido |
      | placa    | ABC12345        | Formato de placa inválido    |

  Cenário: Rejeição por idade mínima
    Quando tento realizar cadastro com data de nascimento "15/03/2010"
    Então vejo a mensagem de erro "Motorista deve ter pelo menos 18 anos"
    E não consigo concluir o cadastro

  Cenário: Rejeição por CNH vencida
    Quando tento realizar cadastro com CNH válida até "15/03/2020"
    Então vejo a mensagem de erro "CNH vencida. Renove sua CNH para prosseguir"
    E não consigo concluir o cadastro

  Esquema do Cenário: Tentativa de cadastro com dados já existentes
    Dado que existe um motorista cadastrado com <campo> "<valor_existente>"
    Quando tento realizar cadastro com <campo> "<valor_existente>"
    Então vejo a mensagem de erro "<mensagem_erro>"
    E não consigo concluir o cadastro

    Exemplos:
      | campo | valor_existente      | mensagem_erro              |
      | cpf   | 123.456.789-00       | CPF já cadastrado          |
      | cnh   | 12345678901          | CNH já cadastrada          |
      | email | joao.silva@email.com | Email já cadastrado        |

  Cenário: Upload de documentos obrigatórios
    Dado que completei o cadastro básico com sucesso
    Quando faço upload dos documentos obrigatórios:
      | documento        | formato | tamanho |
      | CNH              | JPG     | 2MB     |
      | CRLV             | PNG     | 1.5MB   |
      | selfie com CNH   | JPG     | 1MB     |
    Então vejo a mensagem "Documentos enviados com sucesso"
    E meu status muda para "documentos_em_analise"

  Esquema do Cenário: Rejeição de upload de documento inválido
    Dado que estou na página de upload de documentos
    Quando tento fazer upload de arquivo com "<problema>"
    Então vejo a mensagem de erro "<mensagem_erro>"
    E o upload não é concluído

    Exemplos:
      | problema           | mensagem_erro                                  |
      | tamanho de 6MB     | Arquivo muito grande. Tamanho máximo: 5MB     |
      | formato TXT        | Formato não suportado. Use JPG, PNG ou PDF    |

  Cenário: Validação de força de senha
    Dado que eu estou na página de cadastro de motorista
    Quando eu preencho a senha "123456"
    Então eu vejo o indicador de força da senha como "Fraca"
    E eu vejo a mensagem "Senha deve ter pelo menos 8 caracteres, incluindo maiúscula, minúscula, número e símbolo"
    Quando eu altero a senha para "MinhaSenh@123"
    Então eu vejo o indicador de força da senha como "Forte"
    E a mensagem de erro desaparece

  Cenário: Confirmação de senha diferente
    Dado que eu estou na página de cadastro de motorista
    Quando eu preencho a senha "MinhaSenh@123"
    E eu preencho a confirmação de senha "MinhaSenh@456"
    E eu submeto o formulário
    Então eu vejo a mensagem de erro "Senhas não conferem"
    E os campos de senha são destacados

  Cenário: Upload de documentos obrigatórios
    Dado que eu completei o cadastro básico com sucesso
    E eu estou na página de upload de documentos
    Quando eu faço upload da foto da CNH em formato JPG com 2MB
    E eu faço upload da foto do CRLV em formato PNG com 1.5MB
    E eu faço upload de uma selfie com CNH em formato JPG com 1MB
    Então eu vejo "Documentos enviados com sucesso"
    E eu vejo o status "Documentos em análise"
    E eu recebo um email confirmando o recebimento dos documentos

  Cenário: Tentativa de upload com arquivo muito grande
    Dado que eu estou na página de upload de documentos
    Quando eu tento fazer upload de um arquivo de 6MB
    Então eu vejo a mensagem de erro "Arquivo muito grande. Tamanho máximo: 5MB"
    E o upload é cancelado

  Cenário: Tentativa de upload com formato inválido
    Dado que eu estou na página de upload de documentos
    Quando eu tento fazer upload de um arquivo em formato TXT
    Então eu vejo a mensagem de erro "Formato não suportado. Use JPG, PNG ou PDF"
    E o upload é cancelado

  Cenário: Verificação automática de documentos aprovada
    Dado que eu enviei todos os documentos obrigatórios
    E o sistema de verificação automática está funcionando
    Quando o sistema valida meus documentos
    E todos os documentos são aprovados automaticamente
    Então meu status muda para "aprovado"
    E eu recebo um email "Parabéns! Seu cadastro foi aprovado"
    E eu posso fazer login no sistema
    E eu vejo o tutorial inicial do motorista

  Cenário: Verificação de documentos rejeitada
    Dado que eu enviei documentos com problemas de qualidade
    Quando o sistema ou analista verifica meus documentos
    Então meu status muda para "rejeitado"
    E eu recebo um email detalhando os problemas encontrados
    E eu posso reenviar os documentos corrigidos
    E eu vejo quais documentos específicos foram rejeitados
