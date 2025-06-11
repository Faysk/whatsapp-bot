-- 📦 Tabela de dispositivos registrados (WhatsMeow)
CREATE TABLE IF NOT EXISTS whatsmeow_device (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    device_id BYTEA NOT NULL,
    user_id TEXT NOT NULL,
    agent TEXT NOT NULL,
    platform TEXT NOT NULL,
    version TEXT NOT NULL,
    full_id TEXT NOT NULL UNIQUE,
    push_name TEXT,
    business_name TEXT,
    business_desc TEXT,
    business_categories TEXT[],
    business_profile_pic_url TEXT
);

-- 💾 Sessões persistidas do WhatsApp
CREATE TABLE IF NOT EXISTS sessions (
    id TEXT PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    data BYTEA NOT NULL
);

-- 🔐 Identidades Signal (usuário + dispositivo)
CREATE TABLE IF NOT EXISTS identities (
    user_id TEXT NOT NULL,
    device_id BYTEA NOT NULL,
    identity BYTEA NOT NULL,
    PRIMARY KEY (user_id, device_id)
);

-- 🔑 Pre-keys (Signal)
CREATE TABLE IF NOT EXISTS prekeys (
    id SERIAL PRIMARY KEY,
    key_id INTEGER NOT NULL CHECK (key_id > 0),
    key_data BYTEA NOT NULL
);

-- 🔐 Signed Pre-keys (Signal)
CREATE TABLE IF NOT EXISTS signed_prekeys (
    id SERIAL PRIMARY KEY,
    key_id INTEGER NOT NULL CHECK (key_id > 0),
    key_data BYTEA NOT NULL,
    signature BYTEA NOT NULL
);

-- 🔁 App State Sync Keys
CREATE TABLE IF NOT EXISTS app_state_sync_keys (
    key_id TEXT PRIMARY KEY,
    key_data BYTEA NOT NULL,
    fingerprint BYTEA NOT NULL
);

-- 🔄 Versões do estado da aplicação
CREATE TABLE IF NOT EXISTS app_state_versions (
    key_id TEXT NOT NULL,
    version BYTEA NOT NULL,
    PRIMARY KEY (key_id, version)
);
