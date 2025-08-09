	SELECT usename AS username, COUNT(*) AS active_connections
FROM pg_stat_activity
WHERE 1=1 
-- and state = 'active'
GROUP BY usename
ORDER BY active_connections DESC;


SELECT 
    usename as user,
    datname as database,
    count(*) as total_connections,
    count(*) FILTER (WHERE state = 'idle') as idle_connections,
    count(*) FILTER (WHERE state = 'active') as active_connections,
    count(*) FILTER (WHERE state = 'idle in transaction') as idle_in_transaction
FROM pg_stat_activity 
WHERE usename IS NOT NULL
GROUP BY usename, datname
ORDER BY idle_connections DESC;

SELECT 
    usename as user,
    datname as database,
    count(*) as total_connections,
    count(*) FILTER (WHERE state = 'idle') as idle_connections,
    count(*) FILTER (WHERE state = 'active') as active_connections,
    count(*) FILTER (WHERE state = 'idle in transaction') as idle_in_transaction,
    date_trunc('hour', query_start) as hour_start,
    date_trunc('minute', query_start) as minute_start
FROM pg_stat_activity 
WHERE usename IS NOT NULL
    AND query_start IS NOT NULL
GROUP BY usename, datname, hour_start, minute_start
ORDER BY hour_start DESC, minute_start DESC, idle_connections DESC;



SELECT 
    usename as user,
    datname as database,
    count(*) as total_connections,
    count(*) FILTER (WHERE state = 'idle') as idle_connections,
    count(*) FILTER (WHERE state = 'active') as active_connections,
    count(*) FILTER (WHERE state = 'idle in transaction') as idle_in_transaction,
    date_trunc('minute', query_start) as query_minute
FROM pg_stat_activity 
WHERE usename IS NOT NULL
--     AND query_start >= '2025-08-07 00:35:00'
--     AND query_start <= '2025-08-09 23:59:00'
GROUP BY usename, datname, query_minute
ORDER BY query_minute DESC, idle_connections DESC;


SELECT 
    usename,
    datname,
    state,
    query_start,
    now() - query_start as duration,
    query
FROM pg_stat_activity 
WHERE datname = 'hdb_beneficiary'
    AND state = 'idle in transaction'
ORDER BY duration DESC;



SELECT 
    usename,
    datname,
    state,
    now() - state_change as idle_duration,
    count(*) as connection_count
FROM pg_stat_activity 
WHERE datname = 'hdb_beneficiary'
GROUP BY usename, datname, state, idle_duration
ORDER BY idle_duration DESC;
SHOW idle_session_timeout;
SHOW idle_in_transaction_session_timeout;
SHOW statement_timeout;



SELECT 
    COALESCE(usename, '[system process]') as user,
    COALESCE(datname, '[no database]') as database,
    backend_type,
    count(*) as total_connections,
    count(*) FILTER (WHERE state = 'idle') as idle,
    count(*) FILTER (WHERE state = 'active') as active,
    count(*) FILTER (WHERE state = 'idle in transaction') as idle_in_txn,
    count(*) FILTER (WHERE state = 'idle in transaction (aborted)') as idle_in_txn_aborted,
    count(*) FILTER (WHERE wait_event_type = 'Lock') as lock_waits,
    count(*) FILTER (WHERE wait_event_type = 'LWLock') as lwlock_waits,
    count(*) FILTER (WHERE wait_event_type = 'BufferPin') as buffer_pin_waits,
    count(*) FILTER (WHERE wait_event_type = 'Activity') as activity_waits,
    count(*) FILTER (WHERE wait_event_type = 'Client') as client_waits,
    count(*) FILTER (WHERE wait_event_type = 'Extension') as extension_waits,
    count(*) FILTER (WHERE wait_event_type = 'IPC') as ipc_waits,
    count(*) FILTER (WHERE wait_event_type = 'Timeout') as timeout_waits,
    count(*) FILTER (WHERE wait_event_type = 'IO') as io_waits,
    -- Most common wait events
    mode() WITHIN GROUP (ORDER BY wait_event) as most_common_wait_event,
    date_trunc('minute', COALESCE(query_start, backend_start)) as time_minute
FROM pg_stat_activity 
WHERE backend_start >= '2025-08-07 00:35:00'
    AND backend_start <= '2025-08-09 23:59:00'
GROUP BY 
    COALESCE(usename, '[system process]'), 
    COALESCE(datname, '[no database]'),
    backend_type,
    time_minute
ORDER BY time_minute DESC, total_connections DESC;


SELECT pid, usename, datname, state, backend_start, state_change, query_start, query
FROM pg_stat_activity
WHERE state != 'idle'
ORDER BY query_start;


SELECT count(*) as current_connections FROM pg_stat_activity, pg_settings WHERE name = 'max_connections';






