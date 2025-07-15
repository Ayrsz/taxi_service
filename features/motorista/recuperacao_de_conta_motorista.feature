# language: pt
Funcionalidade: Recuperação de Conta do Motorista
  Como um motorista que esqueceu suas credenciais
  Eu quero recuperar o acesso à minha conta
  Para voltar a trabalhar e acessar o sistema

  Contexto:
    Dado que existe um motorista cadastrado com os dados:
      | nome   | José Santos             |
      | email  | jose.santos@email.com   |
      | cpf    | 987.654.321-00          |
      | status | ativo                   |
    E o sistema de envio de emails está funcionando

  Cenário: Solicitar recuperação com email válido
    Dado que estou na página "Recuperação de Conta"
    Quando submeto o formulário com email "jose.santos@email.com"
    Então vejo a mensagem "Instruções de recuperação enviadas para jose.santos@email.com"
    E um email de recuperação é enviado para "jose.santos@email.com"
    E o email contém um link válido por 1 hora

  Cenário: Tentativa de recuperação com email não cadastrado
    Dado que estou na página "Recuperação de Conta"
    Quando submeto o formulário com email "naoexiste@email.com"
    Então vejo a mensagem de erro "Este email não está cadastrado em nosso sistema"
    E nenhum email de recuperação é enviado

  Esquema do Cenário: Validação de formato de email
    Dado que estou na página "Recuperação de Conta"
    Quando submeto o formulário com email "<email_invalido>"
    Então vejo a mensagem de erro "<mensagem_erro>"
    E nenhum email é enviado

    Exemplos:
      | email_invalido  | mensagem_erro                     |
      |                 | Email é obrigatório               |
      | email_invalido  | Formato de email inválido         |
      | @email.com      | Formato de email inválido         |
      | email@          | Formato de email inválido         |

  Cenário: Redefinir senha com sucesso
    Dado que solicitei recuperação de conta há 30 minutos
    E acessei o link de recuperação válido
    Quando preencho a nova senha "NovaSenha@123"
    E confirmo a nova senha "NovaSenha@123"
    E submeto o formulário
    Então vejo a mensagem "Senha redefinida com sucesso"
    E posso fazer login com a nova senha "NovaSenha@123"
    E o link de recuperação é invalidado

  Esquema do Cenário: Validação na redefinição de senha
    Dado que estou na página "Redefinir Senha" através de link válido
    Quando preencho a nova senha "<nova_senha>"
    E confirmo a nova senha "<confirmacao>"
    E submeto o formulário
    Então vejo a mensagem de erro "<mensagem_erro>"

    Exemplos:
      | nova_senha     | confirmacao    | mensagem_erro                                    |
      | 123456         | 123456         | Senha deve ter pelo menos 8 caracteres e incluir maiúscula, minúscula, número e símbolo |
      | NovaSenha@123  | NovaSenha@456  | Senhas não conferem                              |
      |                | NovaSenha@123  | Nova senha é obrigatória                         |
      | NovaSenha@123  |                | Confirmação de senha é obrigatória               |

  Cenário: Tentativa de acesso com link expirado
    Dado que solicitei recuperação de conta há 2 horas
    Quando acesso o link de recuperação expirado
    Então vejo a mensagem "Este link de recuperação expirou"
    E vejo a opção de solicitar um novo link

  Cenário: Tentativa de reutilizar link já usado
    Dado que redefini minha senha usando um link de recuperação
    Quando tento acessar o mesmo link novamente
    Então vejo a mensagem "Este link já foi utilizado e não é mais válido"
    
  Cenário: Recuperação para conta com status específico
    Dado que existe um motorista com status "<status>":
      | email  | <email>             |
      | status | <status>            |
    Quando solicito recuperação para "<email>"
    Então vejo a mensagem "<mensagem>"
    E <acao_email>

    Exemplos:
      | status    | email                 | mensagem                                                    | acao_email                    |
      | inativo   | maria.silva@email.com | Esta conta está desativada. Entre em contato com o suporte  | nenhum email é enviado        |
      | suspenso  | pedro.lima@email.com  | Esta conta está suspensa. Entre em contato com o suporte    | nenhum email é enviado        |

  Cenário: Múltiplas solicitações de recuperação
    Dado que solicitei recuperação de conta há 3 minutos
    Quando solicito recuperação novamente com o mesmo email
    Então vejo a mensagem "Um email de recuperação já foi enviado recentemente"
    E vejo que devo aguardar antes de solicitar novamente