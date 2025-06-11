-- sql/psql_init.sql
-- Migrações WhatsMeow (v1 → v11), idempotente, atômico e completo

BEGIN;

-------------------------------
-- v1: Tabela de versão e schema inicial
-------------------------------
CREATE TABLE IF NOT EXISTS whatsmeow_version (
  version INTEGER PRIMARY KEY
);
INSERT INTO whatsmeow_version (version)
  VALUES (1)
  ON CONFLICT DO NOTHING;

CREATE TABLE IF NOT EXISTS whatsmeow_device (
  jid TEXT PRIMARY KEY,
  registration_id BIGINT NOT NULL
    CHECK (registration_id >= 0 AND registration_id < 4294967296),
  noise_key BYTEA NOT NULL CHECK (length(noise_key) = 32),
  identity_key BYTEA NOT NULL CHECK (length(identity_key) = 32),
  signed_pre_key BYTEA NOT NULL CHECK (length(signed_pre_key) = 32),
  signed_pre_key_id INTEGER NOT NULL
    CHECK (signed_pre_key_id >= 0 AND signed_pre_key_id < 16777216),
  signed_pre_key_sig BYTEA NOT NULL CHECK (length(signed_pre_key_sig) = 64),
  adv_key BYTEA NOT NULL,
  adv_details BYTEA NOT NULL,
  adv_account_sig BYTEA NOT NULL CHECK (length(adv_account_sig) = 64),
  adv_device_sig BYTEA NOT NULL CHECK (length(adv_device_sig) = 64),
  platform TEXT NOT NULL DEFAULT '',
  business_name TEXT NOT NULL DEFAULT '',
  push_name TEXT NOT NULL DEFAULT ''
);

CREATE TABLE IF NOT EXISTS whatsmeow_identity_keys (
  our_jid TEXT NOT NULL,
  their_id TEXT NOT NULL,
  identity BYTEA NOT NULL CHECK (length(identity) = 32),
  PRIMARY KEY (our_jid, their_id),
  FOREIGN KEY (our_jid)
    REFERENCES whatsmeow_device(jid)
      ON DELETE CASCADE
      ON UPDATE CASCADE
);

-------------------------------
-- v2: coluna adv_account_sig_key
-------------------------------
ALTER TABLE IF EXISTS whatsmeow_device
  ADD COLUMN IF NOT EXISTS adv_account_sig_key BYTEA
    CHECK (length(adv_account_sig_key) = 32);
INSERT INTO whatsmeow_version (version)
  VALUES (2)
  ON CONFLICT DO NOTHING;

-------------------------------
-- v3: tabela de message_secrets
-------------------------------
CREATE TABLE IF NOT EXISTS whatsmeow_message_secrets (
  our_jid TEXT NOT NULL,
  chat_jid TEXT NOT NULL,
  sender_jid TEXT NOT NULL,
  message_id TEXT NOT NULL,
  key BYTEA NOT NULL,
  PRIMARY KEY (our_jid, chat_jid, sender_jid, message_id),
  FOREIGN KEY (our_jid)
    REFERENCES whatsmeow_device(jid)
      ON DELETE CASCADE
      ON UPDATE CASCADE
);
INSERT INTO whatsmeow_version (version)
  VALUES (3)
  ON CONFLICT DO NOTHING;

-------------------------------
-- v4: coluna lid e índice único
-------------------------------
ALTER TABLE IF EXISTS whatsmeow_device
  ADD COLUMN IF NOT EXISTS lid TEXT NOT NULL DEFAULT '';

DO $$
BEGIN
  IF NOT EXISTS (
    SELECT 1
      FROM pg_constraint
     WHERE conname = 'uq_device_lid'
  ) THEN
    ALTER TABLE whatsmeow_device
      ADD CONSTRAINT uq_device_lid UNIQUE (lid);
  END IF;
END
$$;
INSERT INTO whatsmeow_version (version)
  VALUES (4)
  ON CONFLICT DO NOTHING;

-------------------------------
-- v5: coluna facebook_uuid
-------------------------------
ALTER TABLE IF EXISTS whatsmeow_device
  ADD COLUMN IF NOT EXISTS facebook_uuid TEXT DEFAULT '';
INSERT INTO whatsmeow_version (version)
  VALUES (5)
  ON CONFLICT DO NOTHING;

-------------------------------
-- v6: tabela de mapeamento lid ↔ pn
-------------------------------
CREATE TABLE IF NOT EXISTS whatsmeow_lid_map (
  lid TEXT NOT NULL,
  pn TEXT NOT NULL,
  PRIMARY KEY (lid, pn),
  UNIQUE (lid)
);
INSERT INTO whatsmeow_version (version)
  VALUES (6)
  ON CONFLICT DO NOTHING;

-------------------------------
-- v7: tabela de pre_keys
-------------------------------
CREATE TABLE IF NOT EXISTS whatsmeow_pre_keys (
  jid TEXT NOT NULL,
  key_id INTEGER NOT NULL,
  key BYTEA NOT NULL,
  uploaded BOOLEAN NOT NULL DEFAULT false,
  PRIMARY KEY (jid, key_id),
  FOREIGN KEY (jid)
    REFERENCES whatsmeow_device(jid)
      ON DELETE CASCADE
      ON UPDATE CASCADE
);
INSERT INTO whatsmeow_version (version)
  VALUES (7)
  ON CONFLICT DO NOTHING;

-------------------------------
-- v8: tabela de sessions
-------------------------------
CREATE TABLE IF NOT EXISTS whatsmeow_sessions (
  our_jid TEXT NOT NULL,
  their_id TEXT NOT NULL,
  session BYTEA NOT NULL,
  PRIMARY KEY (our_jid, their_id),
  FOREIGN KEY (our_jid)
    REFERENCES whatsmeow_device(jid)
      ON DELETE CASCADE
      ON UPDATE CASCADE
);
INSERT INTO whatsmeow_version (version)
  VALUES (8)
  ON CONFLICT DO NOTHING;

-------------------------------
-- v9: tabela de app_state_version
-------------------------------
CREATE TABLE IF NOT EXISTS whatsmeow_app_state_version (
  jid TEXT NOT NULL,
  name TEXT NOT NULL,
  version INTEGER NOT NULL,
  hash BYTEA NOT NULL,
  PRIMARY KEY (jid, name),
  FOREIGN KEY (jid)
    REFERENCES whatsmeow_device(jid)
      ON DELETE CASCADE
      ON UPDATE CASCADE
);
INSERT INTO whatsmeow_version (version)
  VALUES (9)
  ON CONFLICT DO NOTHING;

-------------------------------
-- v10: tabelas adicionais de sincronização de estado
-------------------------------
CREATE TABLE IF NOT EXISTS whatsmeow_sender_keys (
  our_jid    TEXT NOT NULL,
  chat_id    TEXT NOT NULL,
  sender_id  TEXT NOT NULL,
  sender_key BYTEA NOT NULL,
  PRIMARY KEY (our_jid, chat_id, sender_id),
  FOREIGN KEY (our_jid)
    REFERENCES whatsmeow_device(jid)
      ON DELETE CASCADE
      ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS whatsmeow_contacts (
  our_jid       TEXT NOT NULL,
  their_jid     TEXT NOT NULL,
  first_name    TEXT,
  full_name     TEXT,
  push_name     TEXT,
  business_name TEXT,
  PRIMARY KEY (our_jid, their_jid),
  FOREIGN KEY (our_jid)
    REFERENCES whatsmeow_device(jid)
      ON DELETE CASCADE
      ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS whatsmeow_app_state_sync_keys (
  jid         TEXT NOT NULL,
  key_id      BIGINT NOT NULL,
  key_data    BYTEA NOT NULL,
  timestamp   BIGINT NOT NULL,
  fingerprint BYTEA NOT NULL,
  PRIMARY KEY (jid, key_id),
  FOREIGN KEY (jid)
    REFERENCES whatsmeow_device(jid)
      ON DELETE CASCADE
      ON UPDATE CASCADE
);
INSERT INTO whatsmeow_version (version)
  VALUES (10)
  ON CONFLICT DO NOTHING;

-------------------------------
-- v11: tabela de privacy_tokens
-------------------------------
CREATE TABLE IF NOT EXISTS whatsmeow_privacy_tokens (
  our_jid    TEXT NOT NULL,
  their_jid  TEXT NOT NULL,
  token      BYTEA NOT NULL,
  timestamp  BIGINT NOT NULL,
  PRIMARY KEY (our_jid, their_jid),
  FOREIGN KEY (our_jid)
    REFERENCES whatsmeow_device(jid)
      ON DELETE CASCADE
      ON UPDATE CASCADE
);
INSERT INTO whatsmeow_version (version)
  VALUES (11)
  ON CONFLICT DO NOTHING;

COMMIT;
