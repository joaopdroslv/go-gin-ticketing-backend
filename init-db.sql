
CREATE DATABASE IF NOT EXISTS main;

USE main;

CREATE TABLE IF NOT EXISTS user_statuses (
    id              BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name            VARCHAR(32) NOT NULL UNIQUE,
    description     TEXT DEFAULT NULL,

    created_at      DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

INSERT INTO user_statuses (name, description) VALUES
    ("active", NULL),
    (
        "inactive",
        "A inactive user cannot perform any action."
    ),
    (
        "email_confirmation",
        "When a user's account is awaiting email confirmation, it's not considered active yet."
    ),
    (
        "deleted_account",
        "When a user requests the deletion of their account, it's considered inactive."
    )
;

SET @active_user_status_id := (
    SELECT id FROM user_statuses WHERE name = "active" LIMIT 1
);

CREATE TABLE IF NOT EXISTS users (
    id              BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    status_id       BIGINT UNSIGNED NOT NULL,

    name            VARCHAR(128) NOT NULL,
    birthdate       DATE NOT NULL,
    -- Authentication fields
    email           VARCHAR(128) UNIQUE NOT NULL,
    password_hash   VARCHAR(255) NOT NULL,

    created_at      DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- System root user
INSERT INTO users (name, birthdate, status_id, email, password_hash)
VALUES (
    "system",
    "2000-01-01",
    @active_user_status_id,
    "system@system.com",
    -- This bcrypt(12) hash translates to "systemsystem123"
    "$2a$12$sZ.BjwbUgXAigyfepBLH7uUXijODjjRUMEGEKRKCitjAN8yciNjhe"
);

SET @system_user_id := (
    SELECT id FROM users WHERE email = "system@system.com" LIMIT 1
);

CREATE TABLE IF NOT EXISTS roles (
    id              BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,

    name            VARCHAR(64) NOT NULL UNIQUE,
    description     TEXT DEFAULT NULL,

    created_at      DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

INSERT INTO roles (name, description) VALUES
    (
        "system",
        "A super user role, can perform any action."
    ),
    (
        "common",
        "Common user, can perform basic ticket operations, like create, read, update and close."
    )
;

SET @system_role_id := (
    SELECT id FROM roles WHERE name = "system" LIMIT 1
);

SET @common_role_id := (
    SELECT id FROM roles WHERE name = "common" LIMIT 1
);

CREATE TABLE IF NOT EXISTS permissions (
    id              BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,

    name            VARCHAR(100) NOT NULL UNIQUE,
    description     TEXT DEFAULT NULL,

    created_at      DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

INSERT INTO permissions (name, description) VALUES
    (
        "user:read",
        "The role has permission to read users info."
    ),
    (
        "user:list",
        "The role has permission to list and read multiple users info."
    ),
    (
        "user:create",
        "The role has permission to create a new user."
    ),
    (
        "user:update",
        "The role has permission to update a user."
    ),
    (
        "user:delete",
        "The role has permission to delete a user."
    ),
    -- (
    --     "user:inactivate",
    --     "The role has permission to inactivate a user."
    -- ),
    (
        "ticket:read",
        "The role has permission to read tickets info."
    ),
    (
        "ticket:list",
        "The role has permission to list and read multiple tickets info."
    ),
    (
        "ticket:create",
        "The role has permission to open a new ticket."
    ),
    (
        "ticket:update",
        "The role has permission to update a ticket."
    ),
    (
        "ticket:close",
        "The role has permission to close a ticket."
    )
;

SET @user_read   := (SELECT id FROM permissions WHERE name = "user:read");
SET @user_list   := (SELECT id FROM permissions WHERE name = "user:list");
SET @user_create := (SELECT id FROM permissions WHERE name = "user:create");
SET @user_update := (SELECT id FROM permissions WHERE name = "user:update");
SET @user_delete := (SELECT id FROM permissions WHERE name = "user:delete");

SET @ticket_read   := (SELECT id FROM permissions WHERE name = "ticket:read");
SET @ticket_list   := (SELECT id FROM permissions WHERE name = "ticket:list");
SET @ticket_create := (SELECT id FROM permissions WHERE name = "ticket:create");
SET @ticket_update := (SELECT id FROM permissions WHERE name = "ticket:update");
SET @ticket_close  := (SELECT id FROM permissions WHERE name = "ticket:close");

CREATE TABLE IF NOT EXISTS role_permissions (
    id              BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    role_id         BIGINT UNSIGNED NOT NULL,
    permission_id   BIGINT UNSIGNED NOT NULL,

    created_at      DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    UNIQUE KEY uk_role_permission (role_id, permission_id),

    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
    FOREIGN KEY (permission_id) REFERENCES permissions(id) ON DELETE CASCADE
);

INSERT INTO role_permissions (role_id, permission_id) VALUES
    (@common_role_id, @ticket_read),
    (@common_role_id, @ticket_list),
    (@common_role_id, @ticket_create),
    (@common_role_id, @ticket_update),
    (@common_role_id, @ticket_close),
    (@system_role_id, @user_read),
    (@system_role_id, @user_list),
    (@system_role_id, @user_create),
    (@system_role_id, @user_update),
    (@system_role_id, @user_delete),
    (@system_role_id, @ticket_read),
    (@system_role_id, @ticket_list),
    (@system_role_id, @ticket_create),
    (@system_role_id, @ticket_update),
    (@system_role_id, @ticket_close)
;

CREATE TABLE IF NOT EXISTS role_inheritance (
    id              BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    parent_role_id  BIGINT UNSIGNED NOT NULL,
    child_role_id   BIGINT UNSIGNED NOT NULL,

    created_at      DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    UNIQUE KEY uk_role_inheritance (parent_role_id, child_role_id),

    FOREIGN KEY (parent_role_id) REFERENCES roles(id) ON DELETE CASCADE,
    FOREIGN KEY (child_role_id) REFERENCES roles(id) ON DELETE CASCADE
);

-- INSERT INTO role_inheritance (parent_role_id, child_role_id) VALUES (@common_role_id, @system_role_id);

CREATE TABLE IF NOT EXISTS user_roles (
    id              BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id         BIGINT UNSIGNED NOT NULL,
    role_id         BIGINT UNSIGNED NOT NULL,
    scope_id        BIGINT UNSIGNED DEFAULT NULL, -- Tenant/Project scope, for now it's gonna be null (global)

    created_at      DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    UNIQUE KEY uk_user_role_scope (user_id, role_id, scope_id),

    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE
);

INSERT INTO user_roles (user_id, role_id)
VALUES (@system_user_id, @system_role_id);
