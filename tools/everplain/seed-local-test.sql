-- Supply service_key and upstream_key as psql variables from private runtime values.
-- Synthetic fixture only. Refuses to touch any other database.
DO $$ BEGIN
 IF current_database() <> 'everplain_gateway_test' THEN
  RAISE EXCEPTION 'This fixture is restricted to everplain_gateway_test';
 END IF;
END $$;
INSERT INTO users(email,password_hash,role,balance,concurrency,status)
VALUES('synthetic-service@everplain.test','!synthetic-no-password-login','user',0,12,'active')
ON CONFLICT DO NOTHING;
INSERT INTO accounts(name,platform,type,credentials,concurrency,status)
SELECT 'everplain-local-synthetic','anthropic','apikey',
 jsonb_build_object('api_key', :'upstream_key', 'base_url', 'https://127.0.0.1:18082', 'model_mapping', jsonb_build_object('everplain-test','claude-sonnet-4-20250514')),12,'active'
WHERE NOT EXISTS(SELECT 1 FROM accounts WHERE name='everplain-local-synthetic');
INSERT INTO api_keys(user_id,key,name,group_id,status)
SELECT u.id,:'service_key','everplain-local-synthetic',g.id,'active'
FROM users u CROSS JOIN LATERAL(SELECT id FROM groups WHERE platform='anthropic' AND deleted_at IS NULL ORDER BY id LIMIT 1)g
WHERE u.email='synthetic-service@everplain.test'
ON CONFLICT DO NOTHING;
INSERT INTO account_groups(account_id,group_id)
SELECT a.id,g.id FROM accounts a CROSS JOIN LATERAL(SELECT id FROM groups WHERE platform='anthropic' AND deleted_at IS NULL ORDER BY id LIMIT 1)g
WHERE a.name='everplain-local-synthetic'
ON CONFLICT DO NOTHING;
