## Exchange V1

### Prerequisites

**[Install docker compose](https://docs.docker.com/compose/install/)**

### Run Services

```shell
docker compose up
```

### Worker

Modify Ethereum JSON RPC URL and task settings in `./build/worker/config/worket.dev.yml`

### API Service

The default port is `3000`. This can be changed in the `./docker-compose.yml` file.

#### List User Point Task States

`GET /users/{address}/point-task-states`

```shell
curl 'http://localhost:3000/users/0xF4ACDAC048C14c5E49BbEDe0C72444d806A75Cde/point-task-states'
```

Response body
```json
{
  "status": "SUCCESS",
  "data": [
    {
      "address": "0x67CeA36eEB36Ace126A3Ca6E21405258130CF33C",
      "point": 100,
      "volume": "2660469101",
      "completed": true,
      "updated_at": "2024-07-25T16:54:29.617409Z",
      "task_name": "Onboarding Task",
      "task_type": "ACCOUNT_TRADING_VOLUME",
      "contract_address": "0xB4e16d0168e52d35CaCD2c6185b44281Ec28C9Dc",
      "token_address": "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48",
      "task_volume": "1000000000",
      "task_reward_point": 100,
      "task_status": "FINISHED",
      "start_time": "2024-07-01T00:00:00Z",
      "end_time": "2024-07-01T02:00:00Z"
    },
    {
      "address": "0x67CeA36eEB36Ace126A3Ca6E21405258130CF33C",
      "point": 29,
      "volume": "2660469101",
      "completed": true,
      "updated_at": "2024-07-25T16:54:29.627846Z",
      "task_name": "Share Pool Task",
      "task_type": "SHARE_POOL",
      "contract_address": "0xB4e16d0168e52d35CaCD2c6185b44281Ec28C9Dc",
      "token_address": "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48",
      "task_volume": "927847611814",
      "task_reward_point": 10000,
      "task_status": "FINISHED",
      "start_time": "2024-07-01T00:00:00Z",
      "end_time": "2024-07-01T02:00:00Z"
    }
  ]
}
```

#### List User Point Histories

`GET /users/{address}/point-histories`

```shell
curl 'http://localhost:3000/users/0xF4ACDAC048C14c5E49BbEDe0C72444d806A75Cde/point-histories'
```

Response body
```json
{
    "status": "SUCCESS",
    "data": [
        {
            "address": "0xF4ACDAC048C14c5E49BbEDe0C72444d806A75Cde",
            "point": 100,
            "task_name": "Onboarding Task",
            "task_type": "ACCOUNT_TRADING_VOLUME",
            "created_at": "2024-07-25T16:54:29.617409Z"
        },
        {
            "address": "0xF4ACDAC048C14c5E49BbEDe0C72444d806A75Cde",
            "point": 13,
            "task_name": "Share Pool Task",
            "task_type": "SHARE_POOL",
            "created_at": "2024-07-25T16:54:29.627846Z"
        }
    ]
}
```