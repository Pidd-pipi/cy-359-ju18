-- 数据库初始化脚本（幂等）：由 docker-entrypoint-initdb.d 自动执行
-- 说明：后端连接使用 POSTGRES_USER/POSTGRES_DB 环境变量创建的主用户与库，
-- 应用启动时 GORM AutoMigrate 自动创建业务表（见 backend/migrations/001_init.sql）。

-- 创建业务用户（若不存在）
DO $$
BEGIN
   IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = 'orienteering_user') THEN
      CREATE ROLE orienteering_user LOGIN PASSWORD 'orienteering_pwd';
   END IF;
END
$$;

-- 创建业务数据库（若不存在）
SELECT 'CREATE DATABASE orienteering_db OWNER orienteering_user'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'orienteering_db')\gexec

GRANT ALL PRIVILEGES ON DATABASE orienteering_db TO orienteering_user;
