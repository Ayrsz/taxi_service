# language: pt
Funcionalidade: Gerenciamento de Perfil do Motorista
  Como um motorista autenticado
  Eu quero gerenciar meu perfil e dados pessoais
  Para manter minhas informações sempre atualizadas

  Contexto:
    Dado que estou autenticado como motorista:
      | nome           | João Silva               |
      | email          | joao.silva@email.com     |
      | cpf            | 123.456.789-00           |
      | telefone       | (11) 99999-9999          |
      | data_nascimento| 15/03/1990               |
      | cnh            | 12345678901              |
      | status         | ativo                    |
      | avaliacao      | 4.8                      |

  Cenário: Visualizar perfil completo
    Quando acesso meu perfil
    Então vejo meus dados pessoais cadastrados
    E vejo minha foto de perfil
    E vejo estatísticas sobre minhas corridas

  Cenário: Editar dados pessoais básicos com sucesso
    Dado que estou na página de edição do perfil
    Quando atualizo as seguintes informações:
      | nome     | João Santos Silva    |
      | telefone | (11) 88888-8888      |
    E salvo as alterações
    Então vejo a mensagem "Perfil atualizado com sucesso"
    E meus dados são atualizados no sistema

  Esquema do Cenário: Validação de campos editáveis
    Dado que estou na página de edição do perfil
    Quando preencho o campo "<campo>" com "<valor_invalido>"
    E tento salvar as alterações
    Então vejo a mensagem de erro "<mensagem_erro>"

    Exemplos:
      | campo    | valor_invalido   | mensagem_erro                |
      | nome     |                  | Nome é obrigatório           |
      | telefone | 999999999        | Formato de telefone inválido |
      | telefone | (11) 9999-999    | Telefone incompleto          |

  Cenário: Alterar senha com sucesso
    Dado que estou na página de alteração de senha
    Quando preencho:
      | senha_atual   | MinhaSenh@123  |
      | nova_senha    | NovaSenha@456  |
      | confirmacao   | NovaSenha@456  |
    E salvo a alteração
    Então vejo a mensagem "Senha alterada com sucesso"
    E posso fazer login com a nova senha

  Esquema do Cenário: Validação na alteração de senha
    Dado que estou na página de alteração de senha
    Quando preencho os campos de senha:
      | senha_atual   | <senha_atual>  |
      | nova_senha    | <nova_senha>   |
      | confirmacao   | <confirmacao>  |
    E tento salvar a alteração
    Então vejo a mensagem de erro "<mensagem_erro>"

    Exemplos:
      | senha_atual    | nova_senha     | confirmacao    | mensagem_erro                                |
      | senha_errada   | NovaSenha@456  | NovaSenha@456  | Senha atual incorreta                        |
      | MinhaSenh@123  | 123456         | 123456         | Nova senha deve atender aos critérios mínimos|
      | MinhaSenh@123  | NovaSenha@456  | NovaSenha@789  | Confirmação de senha não confere             |

  Cenário: Upload de nova foto de perfil válida
    Dado que estou na página de edição do perfil
    Quando faço upload de uma nova foto em formato JPG com 1MB
    E salvo a nova foto
    Então vejo a mensagem "Foto de perfil atualizada com sucesso"

  Esquema do Cenário: Upload de foto inválida
    Dado que estou na página de edição do perfil
    Quando tento fazer upload de arquivo "<tipo_arquivo>"
    Então vejo a mensagem de erro "<mensagem_erro>"

    Exemplos:
      | tipo_arquivo           | mensagem_erro                                     |
      | 8MB em formato JPG     | Foto muito grande. Tamanho máximo: 5MB           |
      | 1MB em formato TXT     | Formato não suportado. Use JPG, PNG ou WEBP      |

  Cenário: Desativar conta temporariamente
    Dado que estou na página do meu perfil
    Quando solicito desativar minha conta temporariamente
    E seleciono o motivo "Férias"
    E confirmo a desativação
    Então meu status muda para "inativo_temporario"
    E vejo a mensagem "Conta desativada temporariamente"

  Cenário: Reativar conta
    Dado que minha conta está com status "inativo_temporario"
    Quando solicito reativar minha conta
    E confirmo a reativação
    Então minha conta volta ao status "ativo"
    E vejo a mensagem "Conta reativada com sucesso"

  Cenário: Solicitar exclusão permanente da conta
    Dado que estou na página do meu perfil
    Quando solicito excluir minha conta permanentemente
    E confirmo a exclusão fornecendo a justificativa "Mudança de profissão"
    Então recebo um email com instruções para confirmar a exclusão
    E minha conta fica com status "aguardando_exclusao"
  Cenário: Atualizar documentos vencidos
    Dado que minha CNH está próxima do vencimento (30 dias)
    E eu estou na minha área do motorista
    Então eu vejo um alerta "Sua CNH vence em 25 dias. Atualize seus documentos"
    Quando eu clico em "Atualizar Documentos"
    Então eu sou direcionado para a página de upload de documentos
    E eu posso fazer upload da nova CNH
    E eu vejo o status "Aguardando aprovação dos novos documentos"

  Cenário: Visualizar histórico de alterações
    Dado que eu estou na página do meu perfil
    Quando eu acesso o histórico de alterações
    Então eu vejo uma lista cronológica das modificações:
      | data       | alteracao                    | usuario     |
      | 10/01/2025 | Telefone alterado            | João Silva  |
      | 05/01/2025 | Foto de perfil atualizada    | João Silva  |
      | 01/01/2025 | Senha alterada               | João Silva  |
    E eu vejo apenas as alterações dos últimos 90 dias

  Cenário: Desativar conta temporariamente
    Dado que eu estou na página do meu perfil
    Quando eu clico em "Desativar Conta Temporariamente"
    E eu seleciono o motivo "Férias"
    E eu confirmo a desativação
    Então minha conta fica com status "inativo_temporario"
    E eu não recebo novas solicitações de corrida
    E eu vejo a mensagem "Conta desativada temporariamente. Você pode reativá-la a qualquer momento"
    E eu vejo o botão "Reativar Conta"

  Cenário: Reativar conta
    Dado que minha conta está com status "inativo_temporario"
    E eu estou na minha área do motorista
    Quando eu clico em "Reativar Conta"
    E eu confirmo a reativação
    Então minha conta volta ao status "ativo"
    E eu volto a receber solicitações de corrida
    E eu vejo a mensagem "Conta reativada com sucesso"

  Cenário: Solicitar exclusão permanente da conta
    Dado que eu estou na página do meu perfil
    Quando eu clico em "Excluir Conta Permanentemente"
    Então eu vejo um aviso detalhado sobre as consequências da exclusão
    E eu vejo um formulário para informar o motivo da exclusão
    Quando eu preencho o motivo "Mudança de profissão"
    E eu confirmo digitando "EXCLUIR PERMANENTEMENTE"
    E eu clico em "Confirmar Exclusão"
    Então eu vejo a mensagem "Solicitação de exclusão enviada"
    E eu recebo um email com instruções para confirmar a exclusão
    E minha conta fica com status "aguardando_exclusao"

  Cenário: Cancelar solicitação de exclusão
    Dado que eu solicitei exclusão da conta há 2 dias
    E minha conta está com status "aguardando_exclusao"
    Quando eu acesso o link de cancelamento no email recebido
    Então eu vejo a opção "Cancelar Solicitação de Exclusão"
    Quando eu cancelo a solicitação
    Então minha conta volta ao status "ativo"
    E eu vejo a mensagem "Solicitação de exclusão cancelada"

  Cenário: Tentativa de acesso a perfil de outro motorista
    Dado que existe outro motorista com ID "999"
    Quando eu tento acessar diretamente a URL do perfil do motorista "999"
    Então eu vejo a mensagem de erro "Acesso negado"
    E eu sou redirecionado para o meu próprio perfil
    E o sistema registra a tentativa de acesso não autorizado

  Cenário: Sincronização de dados entre dispositivos
    Dado que eu altero meu telefone no dispositivo mobile
    Quando eu acesso meu perfil no navegador web
    Então eu vejo o telefone atualizado também no navegador
    E eu vejo a data/hora da última sincronização
