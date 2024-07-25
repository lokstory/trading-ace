CREATE EXTENSION IF NOT EXISTS citext;

CREATE DOMAIN UINT256 NUMERIC(78)
    CHECK (value BETWEEN 0 AND 115792089237316195423570985008687907853269984665640564039457584007913129639935);

CREATE DOMAIN ETH_ADDRESS CITEXT
    CHECK (LENGTH(value) = 42);

CREATE TABLE public.block_info
(
    id              BIGINT         NOT NULL PRIMARY KEY,
    block_timestamp TIMESTAMPTZ(0) NOT NULL,
    created_at      TIMESTAMPTZ    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMPTZ    NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE public.address_log_state
(
    id                   BIGSERIAL PRIMARY KEY,
    contract_address     ETH_ADDRESS    NOT NULL UNIQUE,
    topic0               VARCHAR(66)    NULL,
    topic1               VARCHAR(66)    NULL,
    topic2               VARCHAR(66)    NULL,
    topic3               VARCHAR(66)    NULL,
    start_block_number   BIGINT         NOT NULL,
    current_block_number BIGINT         NOT NULL,
    current_block_time   TIMESTAMPTZ(0) NOT NULL,
    created_at           TIMESTAMPTZ    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at           TIMESTAMPTZ    NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE public.swap_event
(
    id               BIGSERIAL PRIMARY KEY,
    contract_address ETH_ADDRESS    NOT NULL,
    block_number     BIGINT         NOT NULL,
    log_index        BIGINT         NOT NULL,
    sender           ETH_ADDRESS    NOT NULL,
    to_address       ETH_ADDRESS    NOT NULL,
    amount0_in       UINT256        NOT NULL,
    amount1_in       UINT256        NOT NULL,
    amount0_out      UINT256        NOT NULL,
    amount1_out      UINT256        NOT NULL,
    block_time       TIMESTAMPTZ(0) NOT NULL,
    created_at       TIMESTAMPTZ    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT swap_event_contract_address_block_number_log_index_uk UNIQUE (contract_address, block_number, log_index)
);

CREATE INDEX swap_event_contract_address_block_time_index
    ON public.swap_event (contract_address, block_time);


CREATE TABLE public.point_task
(
    id                    BIGSERIAL PRIMARY KEY,
    name                  VARCHAR        NOT NULL,
    task_type             VARCHAR(32)    NOT NULL,
    contract_address      ETH_ADDRESS    NOT NULL,
    token_address         ETH_ADDRESS    NOT NULL,
    volume                UINT256        NULL,
    reward_point          BIGINT         NOT NULL,
    settlement_type       VARCHAR(32),
    status                VARCHAR(32)    NOT NULL,
    start_time            TIMESTAMPTZ(0) NOT NULL,
    end_time              TIMESTAMPTZ(0) NOT NULL,
    updated_at_block_time TIMESTAMPTZ(0) NOT NULL,
    created_at            TIMESTAMPTZ    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at            TIMESTAMPTZ    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CHECK (reward_point >= 0)
);

CREATE TABLE public.user_point_task_state
(
    id            BIGSERIAL PRIMARY KEY,
    address       ETH_ADDRESS NOT NULL,
    point_task_id BIGINT      NOT NULL,
    volume        UINT256     NOT NULL,
    point         BIGINT      NOT NULL,
    status        VARCHAR(32) NOT NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT user_point_task_state_point_task_id_fk
        FOREIGN KEY (point_task_id) REFERENCES public.point_task (id),
    CONSTRAINT user_point_task_state_address_point_task_id_uk UNIQUE (address, point_task_id)
);

CREATE TABLE public.user_point_history
(
    id            BIGSERIAL PRIMARY KEY,
    address       ETH_ADDRESS,
    point         BIGINT NOT NULL,
    point_task_id BIGINT NOT NULL,
    created_at    TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT user_point_history_point_task_id_fk
        FOREIGN KEY (point_task_id) REFERENCES public.point_task (id),
    CONSTRAINT user_point_history_address_point_task_id_uk UNIQUE (address, point_task_id)
);

CREATE OR REPLACE PROCEDURE public.update_account_trading_volume_point_task(
  i_task_id BIGINT
)
  LANGUAGE plpgsql
AS
$$
DECLARE
v_task_type                VARCHAR(32);
  v_contract_address         ETH_ADDRESS;
  v_token_address            ETH_ADDRESS;
  v_volume                   UINT256;
  v_reward_point             BIGINT;
  v_settlement_type          VARCHAR(32);
  v_status                   VARCHAR(32);
  v_new_status               VARCHAR(32);
  v_start_time               TIMESTAMPTZ(0);
  v_end_time                 TIMESTAMPTZ(0);
  v_updated_at_block_time    TIMESTAMPTZ(0);
  v_current_log_block_number BIGINT;
  v_current_log_block_time   TIMESTAMPTZ(0);
  v_max_block_time           TIMESTAMPTZ(0);
BEGIN
SELECT task_type,
       contract_address,
       token_address,
       volume,
       reward_point,
       settlement_type,
       status,
       start_time,
       end_time,
       updated_at_block_time
INTO v_task_type,
    v_contract_address,
    v_token_address,
    v_volume,
    v_reward_point,
    v_settlement_type,
    v_status,
    v_start_time,
    v_end_time,
    v_updated_at_block_time
FROM public.point_task
WHERE id = i_task_id
    FOR NO KEY UPDATE;

IF v_task_type <> 'ACCOUNT_TRADING_VOLUME' THEN
    RAISE EXCEPTION 'Invalid task type: %', v_task_type;
END IF;

  IF v_status <> 'CREATED' THEN
    RAISE EXCEPTION 'Invalid task status: %', v_status;
END IF;

SELECT current_block_number, current_block_time
INTO v_current_log_block_number, v_current_log_block_time
FROM public.address_log_state AS s
WHERE s.contract_address = v_contract_address;

IF v_current_log_block_number IS NULL THEN
    RAISE EXCEPTION 'Invalid address log state, address: %', v_contract_address;
END IF;

  IF v_current_log_block_time >= v_end_time THEN
    v_new_status = 'FINISHED';
    v_max_block_time = v_end_time - INTERVAL '1 seconds';
ELSE
    v_new_status = 'CREATED';
    v_max_block_time = v_current_log_block_time;
END IF;

  IF v_updated_at_block_time >= v_max_block_time THEN
    RETURN;
END IF;

INSERT INTO public.user_point_task_state (address, point_task_id, volume, status, point)
SELECT e.to_address, i_task_id, 0, 'CREATED', 0
FROM public.swap_event AS e
WHERE e.contract_address = v_contract_address
  AND e.block_time > v_updated_at_block_time
  AND e.block_time <= v_max_block_time
    ON CONFLICT (address, point_task_id) DO NOTHING;

UPDATE public.point_task
SET status                = v_new_status,
    updated_at_block_time = v_current_log_block_time,
    updated_at            = NOW()
WHERE id = i_task_id;

UPDATE public.user_point_task_state
SET volume     = volume + e.total,
    updated_at = NOW()
    FROM (SELECT e.to_address,
               SUM(e.amount0_in + e.amount0_out) AS total
        FROM public.swap_event AS e
        WHERE e.contract_address = v_contract_address
          AND e.block_time > v_updated_at_block_time
          AND e.block_time <= v_current_log_block_time
        GROUP BY e.to_address) AS e
WHERE address = e.to_address
  AND point_task_id = i_task_id;

UPDATE public.user_point_task_state AS s
SET status     = 'COMPLETED',
    point      = t.reward_point,
    updated_at = NOW()
    FROM public.point_task t
WHERE s.point_task_id = i_task_id
  AND s.status = 'CREATED'
  AND s.volume >= t.volume
  AND t.id = s.point_task_id;

INSERT INTO public.user_point_history (address, point, point_task_id)
SELECT s.address, s.point, s.point_task_id
FROM public.user_point_task_state AS s
WHERE s.point_task_id = i_task_id
  AND s.status = 'COMPLETED' AND s.point > 0
    ON CONFLICT (address, point_task_id) DO NOTHING;
END;
$$;

CREATE OR REPLACE PROCEDURE public.update_share_pool_point_task(
  i_task_id BIGINT
)
  LANGUAGE plpgsql
AS
$$
DECLARE
v_task_type                VARCHAR(32);
  v_contract_address         ETH_ADDRESS;
  v_token_address            ETH_ADDRESS;
  v_volume                   UINT256;
  v_reward_point             BIGINT;
  v_settlement_type          VARCHAR(32);
  v_status                   VARCHAR(32);
  v_start_time               TIMESTAMPTZ(0);
  v_end_time                 TIMESTAMPTZ(0);
  v_updated_at_block_time    TIMESTAMPTZ(0);
  v_current_log_block_number BIGINT;
  v_current_log_block_time   TIMESTAMPTZ(0);
  v_max_block_time           TIMESTAMPTZ(0);
BEGIN
SELECT task_type,
       contract_address,
       token_address,
       reward_point,
       settlement_type,
       status,
       start_time,
       end_time,
       updated_at_block_time
INTO v_task_type,
    v_contract_address,
    v_token_address,
    v_reward_point,
    v_settlement_type,
    v_status,
    v_start_time,
    v_end_time,
    v_updated_at_block_time
FROM public.point_task
WHERE id = i_task_id
    FOR NO KEY UPDATE;

IF v_task_type <> 'SHARE_POOL' THEN
    RAISE EXCEPTION 'Invalid task type: %s', v_task_type;
END IF;

  IF v_status <> 'CREATED' THEN
    RAISE EXCEPTION 'Invalid task status: %s', v_status;
END IF;

SELECT current_block_number, current_block_time
INTO v_current_log_block_number, v_current_log_block_time
FROM public.address_log_state AS s
WHERE s.contract_address = v_contract_address;

IF v_current_log_block_number IS NULL THEN
    RAISE EXCEPTION 'Invalid address log state, address: %s', v_contract_address;
END IF;

  IF v_current_log_block_time < v_end_time THEN
    RETURN;
END IF;

  v_max_block_time = v_end_time - INTERVAL '1 seconds';

  IF v_updated_at_block_time >= v_max_block_time THEN
    RETURN;
END IF;

INSERT INTO public.user_point_task_state (address, point_task_id, volume, status, point)
SELECT e.to_address, i_task_id, 0, 'CREATED', 0
FROM public.swap_event AS e
WHERE e.contract_address = v_contract_address
  AND e.block_time > v_updated_at_block_time
  AND e.block_time <= v_max_block_time
    ON CONFLICT (address, point_task_id) DO NOTHING;

UPDATE public.user_point_task_state
SET volume     = volume + e.total,
    updated_at = NOW()
    FROM (SELECT e.to_address,
               SUM(e.amount0_in + e.amount0_out) AS total
        FROM public.swap_event AS e
        WHERE e.contract_address = v_contract_address
          AND e.block_time > v_updated_at_block_time
          AND e.block_time <= v_current_log_block_time
        GROUP BY e.to_address) AS e
WHERE address = e.to_address
  AND point_task_id = i_task_id;

SELECT SUM(volume) INTO v_volume
FROM user_point_task_state
WHERE point_task_id = i_task_id;

IF v_volume IS NULL THEN
    v_volume = 0;
END IF;

UPDATE public.point_task
SET status                = 'FINISHED',
    volume                = v_volume,
    updated_at_block_time = v_current_log_block_time,
    updated_at            = NOW()
WHERE id = i_task_id;

IF v_volume = 0 THEN
    RETURN;
END IF;

UPDATE public.user_point_task_state AS s
SET point = s.volume * t.reward_point / v_volume
    FROM public.point_task t
WHERE s.point_task_id = i_task_id
  AND s.status = 'CREATED'
  AND t.id = s.point_task_id;

UPDATE public.user_point_task_state AS s
SET status     = 'COMPLETED',
    updated_at = NOW()
    FROM public.point_task t
WHERE s.point_task_id = i_task_id
  AND s.status = 'CREATED'
  AND s.point > 0;


INSERT INTO public.user_point_history (address, point, point_task_id)
SELECT s.address, s.point, s.point_task_id
FROM public.user_point_task_state AS s
WHERE s.point_task_id = i_task_id
  AND s.status = 'COMPLETED' AND s.point > 0
    ON CONFLICT (address, point_task_id) DO NOTHING;
END;
$$;
