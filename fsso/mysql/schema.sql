--
-- Create the tables needed for fsso.
--
-- Run this with the mysql command:
--
--    $ mysql -u <mysqluser> -p <dbname> < schema.sql
--
-- Or from the mysql shell:
--
--    mysql> source schema.sql
--


--
-- One entry per registered user.  Additional fields may be added as needed.
--
CREATE TABLE `fsso_members` (
  `id` serial,
  `email` varchar(255) NOT NULL,
  `fullname` varchar(50) NOT NULL,
  `shortname` varchar(50) NOT NULL,
  `is_active` boolean NOT NULL DEFAULT 1,
  `roles` int(10) unsigned NOT NULL DEFAULT 0,
  `created_at` bigint(20) NOT NULL,
  `active_at` bigint(20) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

--
-- One entry per registered user who has email/password authentication enabled.
--
CREATE TABLE `fsso_auth_email` (
  `member_id` bigint(20) unsigned NOT NULL PRIMARY KEY,
  `email` varchar(255) NOT NULL UNIQUE,
  `pwhash` varchar(64) NOT NULL,
  `pwchanged_at` bigint(20) NOT NULL,
  `is_primary` boolean NOT NULL DEFAULT 1,
  CONSTRAINT `fsso_auth_email_ibfk_1` FOREIGN KEY (`member_id`) REFERENCES `fsso_members` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

--
-- One entry per registered user who has Google authentication enabled.
--
CREATE TABLE `fsso_auth_goog` (
  `member_id` bigint(20) unsigned NOT NULL PRIMARY KEY,
  `uid` varchar(32) NOT NULL UNIQUE,
  `email` varchar(255) NOT NULL,
  `is_primary` boolean NOT NULL DEFAULT 1,
  KEY `email` (`email`),
  CONSTRAINT `fsso_auth_goog_ibfk_1` FOREIGN KEY (`member_id`) REFERENCES `fsso_members` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

--
-- One entry per registered user who has Facebook authentication enabled.
--
CREATE TABLE `fsso_auth_fb` (
  `member_id` bigint(20) unsigned NOT NULL PRIMARY KEY,
  `uid` varchar(32) NOT NULL UNIQUE,
  `email` varchar(255) NOT NULL,
  `is_primary` boolean NOT NULL DEFAULT 1,
  KEY `email` (`email`),
  CONSTRAINT `fsso_auth_fb_ibfk_1` FOREIGN KEY (`member_id`) REFERENCES `fsso_members` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

--
-- One entry per sign-in session.  One user may have more than one session active.
-- is_session is 1 for cooke-based sessions and 0 for token-based sessions.
--
CREATE TABLE `fsso_active` (
  `atoken` varchar(32) NOT NULL PRIMARY KEY,
  `member_id` bigint(20) unsigned NOT NULL,
  `active_at` bigint(20) NOT NULL,
  `useragent` varchar(50) NOT NULL,
  `ip` varchar(50) NOT NULL,
  `is_session` boolean NOT NULL,
  `data` text NOT NULL,
  KEY `member_id` (`member_id`),
  CONSTRAINT `fsso_active_ibfk_1` FOREIGN KEY (`member_id`) REFERENCES `fsso_members` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

--
-- One entry per refresh token, limit one per user.  Refresh tokens are only for
-- re-establishing cookie-based sessions and may only be used once.
--
CREATE TABLE `fsso_refresh` (
  `rtoken` varchar(32) NOT NULL PRIMARY KEY,
  `member_id` bigint(20) unsigned NOT NULL UNIQUE,
  `expires_at` bigint(20) NOT NULL,
  CONSTRAINT `fsso_refresh_ibfk_1` FOREIGN KEY (`member_id`) REFERENCES `fsso_members` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
