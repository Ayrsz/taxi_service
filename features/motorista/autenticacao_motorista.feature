# language: pt
Funcionalidade: Autenticação de Motorista
  Como um motorista cadastrado
  Eu quero fazer login no sistema
  Para acessar minhas funcionalidades de trabalho

  Contexto:
    Dado que existe um motorista cadastrado:
      | nome        | João Silva              |
      | email       | joao.silva@email.com    |
      | cpf         | 123.456.789-00          |
      | cnh         | 12345678901             |
      | status      | ativo                   |
      | senha       | MinhaSenh@123           |

  Cenário: Login bem-sucedido com credenciais válidas
    Quando realizo login com email "joao.silva@email.com" e senha "MinhaSenh@123"
    Então sou autenticado com sucesso
    E vejo o dashboard do motorista
    E vejo "João Silva" como motorista logado

  Esquema do Cenário: Tentativa de login com credenciais inválidas
    Quando realizo login com email "<email>" e senha "<senha>"
    Então vejo a mensagem de erro "<mensagem_erro>"
    E não sou autenticado

    Exemplos:
      | email                    | senha         | mensagem_erro              |
      | joao.silva@email.com     | senhaerrada   | Credenciais inválidas      |
      | inexistente@email.com    | qualquersenha | Credenciais inválidas      |
      | email_invalido           | MinhaSenh@123 | Formato de email inválido  |
      | joao.silva@email.com     |               | Senha é obrigatória        |
      |                          | MinhaSenh@123 | Email é obrigatório        |

  Cenário: Bloqueio temporário após múltiplas tentativas
    Dado que o email "joao.silva@email.com" não está bloqueado
    Quando realizo 5 tentativas de login com senha incorreta para "joao.silva@email.com"
    Então a conta é bloqueada temporariamente por 15 minutos
    E vejo a mensagem "Conta bloqueada temporariamente"
    E não consigo fazer login mesmo com credenciais corretas

  Esquema do Cenário: Login com status de conta específico
    Dado que existe um motorista com status "<status>":
      | email  | <email>             |
      | senha  | <senha>             |
    Quando realizo login com email "<email>" e senha "<senha>"
    Então vejo a mensagem "<mensagem>"
    E <resultado>

    Exemplos:
      | status    | email                 | senha         | mensagem                                  | resultado            |
      | inativo   | maria.santos@email.com| MinhaSenh@456 | Conta inativa. Entre em contato com o suporte | não sou autenticado |
      | suspenso  | pedro.lima@email.com  | MinhaSenh@789 | Conta suspensa. Entre em contato com o suporte | não sou autenticado |

  Cenário: Logout automático por inatividade
    Dado que estou logado como motorista
    Quando fico inativo por mais de 30 minutos
    E tento acessar uma funcionalidade do sistema
    Então sou redirecionado para a página de login
    E vejo a mensagem "Sessão expirada por inatividade"
    Quando eu realizo login com email "joao.silva@email.com" e senha "MinhaSenh@123"
    Então eu vejo a mensagem "Dispositivo não autorizado. Solicite autorização no suporte"
    E uma notificação é enviada para o email do motorista sobre tentativa de acesso

  Cenário: Logout automático por inatividade
    Dado que eu estou logado como motorista "João Silva"
    E eu fico inativo por mais de 30 minutos
    Quando eu tento acessar uma funcionalidade do sistema
    Então eu sou redirecionado para a página de login
    E eu vejo a mensagem "Sessão expirada por inatividade"

  Cenário: Recuperação de sessão após reconexão de internet
    Dado que eu estou logado como motorista "João Silva"
    E eu perco a conexão com a internet por 2 minutos
    E eu recupero a conexão com a internet
    Quando eu tento acessar o dashboard
    Então eu permaneço logado
    E eu vejo uma notificação "Conexão restaurada"
    E meus dados são sincronizados automaticamente
