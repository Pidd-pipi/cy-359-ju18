-- 城市定向越野活动平台初始化脚本（与 GORM AutoMigrate 对应）
-- 说明：应用启动时由 GORM AutoMigrate 自动建表；本脚本供手动初始化数据库结构使用。

CREATE TABLE IF NOT EXISTS users (
    id            BIGSERIAL PRIMARY KEY,
    username      VARCHAR(50)  NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    nickname      VARCHAR(50)  DEFAULT '',
    email         VARCHAR(100) DEFAULT '',
    phone         VARCHAR(20)  DEFAULT '',
    role          VARCHAR(20)  NOT NULL DEFAULT 'user',
    points        INTEGER      NOT NULL DEFAULT 0,
    disabled      BOOLEAN      NOT NULL DEFAULT FALSE,
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS activities (
    id                     BIGSERIAL PRIMARY KEY,
    title                  VARCHAR(120) NOT NULL,
    description            TEXT,
    difficulty             VARCHAR(20)  NOT NULL DEFAULT 'adult',
    duration_minutes       INTEGER      NOT NULL DEFAULT 120,
    equipment_requirement  TEXT,
    start_time             TIMESTAMPTZ  NOT NULL,
    end_time               TIMESTAMPTZ  NOT NULL,
    status                 VARCHAR(20)  NOT NULL DEFAULT 'draft',
    creator_id             BIGINT       NOT NULL,
    start_lat              DOUBLE PRECISION DEFAULT 0,
    start_lng              DOUBLE PRECISION DEFAULT 0,
    end_lat                DOUBLE PRECISION DEFAULT 0,
    end_lng                DOUBLE PRECISION DEFAULT 0,
    address                VARCHAR(255) DEFAULT '',
    max_teams              INTEGER      NOT NULL DEFAULT 50,
    created_at             TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at             TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_activities_status ON activities(status);
CREATE INDEX IF NOT EXISTS idx_activities_creator ON activities(creator_id);

CREATE TABLE IF NOT EXISTS checkpoints (
    id              BIGSERIAL PRIMARY KEY,
    activity_id     BIGINT NOT NULL,
    name            VARCHAR(80) NOT NULL,
    sequence        INTEGER NOT NULL,
    lat             DOUBLE PRECISION DEFAULT 0,
    lng             DOUBLE PRECISION DEFAULT 0,
    clue            VARCHAR(255) DEFAULT '',
    task_type       VARCHAR(20) NOT NULL DEFAULT 'none',
    task_content    TEXT,
    expected_answer VARCHAR(255) DEFAULT '',
    radius_meters   INTEGER NOT NULL DEFAULT 200,
    qr_code         VARCHAR(64) NOT NULL UNIQUE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_checkpoints_activity ON checkpoints(activity_id);

CREATE TABLE IF NOT EXISTS teams (
    id         BIGSERIAL PRIMARY KEY,
    name       VARCHAR(80) NOT NULL UNIQUE,
    slogan     VARCHAR(255) DEFAULT '',
    captain_id BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS team_members (
    id        BIGSERIAL PRIMARY KEY,
    team_id   BIGINT NOT NULL,
    user_id   BIGINT NOT NULL,
    role      VARCHAR(20) NOT NULL DEFAULT 'member',
    joined_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (team_id, user_id)
);

CREATE TABLE IF NOT EXISTS registrations (
    id            BIGSERIAL PRIMARY KEY,
    team_id       BIGINT NOT NULL,
    activity_id   BIGINT NOT NULL,
    status        VARCHAR(20) NOT NULL DEFAULT 'pending',
    start_time    TIMESTAMPTZ,
    finish_time   TIMESTAMPTZ,
    total_seconds INTEGER NOT NULL DEFAULT 0,
    registered_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (team_id, activity_id)
);
CREATE INDEX IF NOT EXISTS idx_registrations_status ON registrations(status);
-- 候补排队/名额统计：按活动 + 状态 + 提交时间取最早候补
CREATE INDEX IF NOT EXISTS idx_registrations_activity_status_time
    ON registrations(activity_id, status, registered_at);

CREATE TABLE IF NOT EXISTS checkin_records (
    id            BIGSERIAL PRIMARY KEY,
    activity_id   BIGINT NOT NULL,
    checkpoint_id BIGINT NOT NULL,
    team_id       BIGINT NOT NULL,
    user_id       BIGINT NOT NULL,
    checkin_type  VARCHAR(20) NOT NULL,
    latitude      DOUBLE PRECISION DEFAULT 0,
    longitude     DOUBLE PRECISION DEFAULT 0,
    answer        VARCHAR(255) DEFAULT '',
    photo_url     VARCHAR(255) DEFAULT '',
    result        VARCHAR(20) NOT NULL DEFAULT 'none',
    points_earned INTEGER NOT NULL DEFAULT 0,
    checked_in_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_checkin_activity ON checkin_records(activity_id);
CREATE INDEX IF NOT EXISTS idx_checkin_team ON checkin_records(team_id);

CREATE TABLE IF NOT EXISTS products (
    id          BIGSERIAL PRIMARY KEY,
    name        VARCHAR(120) NOT NULL,
    description TEXT,
    type        VARCHAR(20) NOT NULL DEFAULT 'equipment',
    points_cost INTEGER NOT NULL,
    stock       INTEGER NOT NULL,
    status      VARCHAR(20) NOT NULL DEFAULT 'on',
    image_url   VARCHAR(255) DEFAULT '',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS redemptions (
    id           BIGSERIAL PRIMARY KEY,
    user_id      BIGINT NOT NULL,
    product_id   BIGINT NOT NULL,
    quantity     INTEGER NOT NULL DEFAULT 1,
    points_cost  INTEGER NOT NULL,
    total_points INTEGER NOT NULL,
    status       VARCHAR(20) NOT NULL DEFAULT 'pending',
    redeemed_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    completed_at TIMESTAMPTZ,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_redemptions_status ON redemptions(status);

CREATE TABLE IF NOT EXISTS favorites (
    id          BIGSERIAL PRIMARY KEY,
    user_id     BIGINT NOT NULL,
    activity_id BIGINT NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (user_id, activity_id)
);

CREATE TABLE IF NOT EXISTS audit_logs (
    id            BIGSERIAL PRIMARY KEY,
    user_id       BIGINT NOT NULL,
    username      VARCHAR(50) DEFAULT '',
    action        VARCHAR(50) NOT NULL,
    resource_type VARCHAR(50) NOT NULL,
    resource_id   VARCHAR(64) DEFAULT '',
    detail        TEXT,
    ip            VARCHAR(50) DEFAULT '',
    request_id    VARCHAR(64) DEFAULT '',
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_audit_logs_request ON audit_logs(request_id);
