
CREATE DATABASE IF NOT EXISTS main;

USE main;

CREATE TABLE IF NOT EXISTS user_statuses (
    id              BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name            VARCHAR(32) NOT NULL UNIQUE,
    description     TEXT DEFAULT NULL,
    created_at      DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

INSERT INTO main.user_statuses (
    name,
    description
) VALUES
    ("active", NULL),
    (
        "inactive",
        "A inactive user cannot perform any action."
    ),
    (
        "email_confirmation",
        "When a user's account is awaiting the email confirmation, it's not considered active yet."
    ),
    (
        "password_creation",
        "When a user's account is awaiting the creation of the password, it's not considered active yet."
    ),
    (
        "deleted_account",
        "When a user requests the deletion of their account, it's considered inactive."
    )
;

SET @active_user_status_id := (SELECT id FROM main.user_statuses WHERE name = "active" LIMIT 1);

CREATE TABLE IF NOT EXISTS user_credentials (
    id              BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    email           VARCHAR(128) UNIQUE NOT NULL,
    password_hash   VARCHAR(255) DEFAULT NULL,
    created_at      DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS users (
    id                  BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_credential_id  BIGINT UNSIGNED NOT NULL,
    user_status_id      BIGINT UNSIGNED NULL,
    name                VARCHAR(128) NOT NULL,
    birthdate           DATE NOT NULL,
    created_at          DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at          DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    CONSTRAINT fk_user_user_credential FOREIGN KEY (user_credential_id) REFERENCES main.user_credentials(id) ON DELETE CASCADE,
    CONSTRAINT fk_user_user_status FOREIGN KEY (user_status_id) REFERENCES main.user_statuses(id) ON DELETE SET NULL,

    UNIQUE KEY uk_user_user_credential (user_credential_id),
    INDEX idx_user_with_status (id, user_status_id)
);

INSERT INTO main.user_credentials (
    email,
    password_hash
) VALUES
    (
        "system@system.com",
        "$2a$12$sZ.BjwbUgXAigyfepBLH7uUXijODjjRUMEGEKRKCitjAN8yciNjhe" -- bcrypt(12) == "systemsystem123"
    )
;

SET @system_user_credential_id := (SELECT id FROM main.user_credentials WHERE email = "system@system.com" LIMIT 1);

INSERT INTO main.users (
    name,
    birthdate,
    user_credential_id,
    user_status_id
) VALUES
    (
        "system",
        "2000-01-01",
        @system_user_credential_id,
        @active_user_status_id
    )
;

SET @system_user_id := (
    SELECT id
    FROM main.users
    WHERE user_credential_id = @system_user_credential_id
    LIMIT 1
);

CREATE TABLE IF NOT EXISTS permissions (
    id              BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name            VARCHAR(100) NOT NULL UNIQUE,
    description     TEXT DEFAULT NULL,
    created_at      DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

INSERT INTO main.permissions (
    name,
    description
) VALUES
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

SET @user_read   := (SELECT id FROM main.permissions WHERE name = "user:read");
SET @user_list   := (SELECT id FROM main.permissions WHERE name = "user:list");
SET @user_create := (SELECT id FROM main.permissions WHERE name = "user:create");
SET @user_update := (SELECT id FROM main.permissions WHERE name = "user:update");
SET @user_delete := (SELECT id FROM main.permissions WHERE name = "user:delete");

SET @ticket_read   := (SELECT id FROM main.permissions WHERE name = "ticket:read");
SET @ticket_list   := (SELECT id FROM main.permissions WHERE name = "ticket:list");
SET @ticket_create := (SELECT id FROM main.permissions WHERE name = "ticket:create");
SET @ticket_update := (SELECT id FROM main.permissions WHERE name = "ticket:update");
SET @ticket_close  := (SELECT id FROM main.permissions WHERE name = "ticket:close");

CREATE TABLE IF NOT EXISTS roles (
    id              BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name            VARCHAR(64) NOT NULL UNIQUE,
    description     TEXT DEFAULT NULL,
    created_at      DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

INSERT INTO main.roles (
    name,
    description
) VALUES
    (
        "system",
        "A super user role, can perform any action."
    ),
    (
        "common",
        "Common user, can perform basic ticket operations, like create, read, update and close."
    )
;

SET @system_role_id := (SELECT id FROM main.roles WHERE name = "system" LIMIT 1);
SET @common_role_id := (SELECT id FROM main.roles WHERE name = "common" LIMIT 1);

CREATE TABLE IF NOT EXISTS user_roles (
    id              BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id         BIGINT UNSIGNED NOT NULL,
    role_id         BIGINT UNSIGNED NOT NULL,
    created_at      DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    CONSTRAINT fk_user_role_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_user_role_role FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,

    UNIQUE KEY uk_user_role (user_id, role_id)
);

INSERT INTO main.user_roles (user_id, role_id) VALUES (@system_user_id, @system_role_id);

CREATE TABLE IF NOT EXISTS role_permissions (
    id              BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    role_id         BIGINT UNSIGNED NOT NULL,
    permission_id   BIGINT UNSIGNED NOT NULL,
    created_at      DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    CONSTRAINT fk_role_permission_role FOREIGN KEY (role_id) REFERENCES main.roles(id) ON DELETE CASCADE,
    CONSTRAINT fk_role_permission_permission FOREIGN KEY (permission_id) REFERENCES main.permissions(id) ON DELETE CASCADE,

    UNIQUE KEY uk_role_permission (role_id, permission_id)
);

INSERT INTO main.role_permissions (
    role_id,
    permission_id
) VALUES
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
