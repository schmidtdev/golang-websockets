## Visão Geral

Este repositório contém um serviço WebSocket baseado em Golang. O serviço escuta conexões WebSocket no endpoint `/ws` e processa as mensagens de acordo.

## Estrutura da Mensagem

A struct `Message` possui as seguintes propriedades:

- **channel**: O nome da tabela do banco de dados (por exemplo, `webgr_pedidos`).
- **type**: O tipo da mensagem. Valores possíveis são:
  - `subscription`
  - `unsubscription`
  - `get`
  - `insert`
  - `update`
  - `broadcast`
- **content**: O conteúdo da mensagem.

## Endpoints

- **/ws**: O endpoint WebSocket para receber mensagens.

## Uso

Para usar este serviço WebSocket, conecte-se ao endpoint `/ws` e envie mensagens no formato definido pela struct `Message`. O serviço irá lidar com as mensagens com base no seu tipo e canal.

## Exemplo

Aqui está um exemplo de uma mensagem:

```json
{
  "channel": "webgr_pedidos",
  "type": "insert",
  "content": "Seu conteúdo aqui"
}
```

## Licença

Este projeto está licenciado sob a Licença MIT.