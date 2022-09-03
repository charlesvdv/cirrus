use anyhow::{bail, Result};
use chrono::{DateTime, Duration, TimeZone, Utc};
use rand::distributions::{Alphanumeric, DistString};
use serde::{Deserialize, Serialize};

use super::users::{PasswordHash, UserId};

#[derive(thiserror::Error, Debug)]
pub enum AuthError {
    #[error("User name or password is invalid")]
    UserNameOrPasswordInvalid,
    #[error("Invalid session token")]
    InvalidSessionToken,
}

#[derive(Deserialize, Debug)]
pub struct LoginUser {
    name: String,
    password: String,
}

#[derive(Serialize, Debug)]
pub struct Session {
    token: String,
    user_id: UserId,
    expired_at: DateTime<Utc>,
}

impl Session {
    fn generate(user_id: UserId) -> Session {
        let token = Alphanumeric.sample_string(&mut rand::thread_rng(), 64);

        Session {
            token,
            user_id,
            expired_at: Utc::now() + Duration::minutes(30),
        }
    }

    fn verify_session(&self, token: &str) -> Result<(), String> {
        if self.token != token {
            log::error!("Token are not matching. Not normal...");
            return Err(String::from("Token are not matching??"));
        }

        if Utc::now() < self.expired_at {
            Ok(())
        } else {
            Err(format!(
                "Token has expired since {}",
                (Utc::now() - self.expired_at).to_string()
            ))
        }
    }
}

pub async fn authenticate(conn: &mut sqlx::SqliteConnection, login: LoginUser) -> Result<Session> {
    let user = sqlx::query!(
        "SELECT id, name, password FROM user WHERE name = ?",
        login.name
    )
    .fetch_optional(&mut *conn)
    .await?
    .ok_or_else(|| {
        log::debug!("User not found: {}", login.name);
        AuthError::UserNameOrPasswordInvalid
    })?;

    let user_password = PasswordHash::from_str(user.password);
    if user_password.verify_password(&login.password).is_err() {
        log::debug!("Password are not matching");
        bail!(AuthError::UserNameOrPasswordInvalid);
    }

    // TODO: why do we need a unwrap() here? sqlx bug?
    let session = Session::generate(UserId::new(user.id.unwrap()));

    let expired_at = session.expired_at.timestamp();
    sqlx::query!(
        "INSERT INTO session(token, user, expired_at) VALUES(?, ?, ?)",
        session.token,
        session.user_id,
        expired_at
    )
    .execute(&mut *conn)
    .await?;

    Ok(session)
}

pub async fn verify_session_token(conn: &mut sqlx::SqliteConnection, token: String) -> Result<()> {
    let session = sqlx::query!(
        "SELECT token, user, expired_at FROM session WHERE token = ?",
        token
    )
    .fetch_optional(&mut *conn)
    .await?
    .ok_or_else(|| {
        log::debug!("Session token not found");
        AuthError::InvalidSessionToken
    })?;

    let session = Session {
        token: session.token,
        user_id: UserId::new(session.user),
        expired_at: Utc.timestamp(session.expired_at, 0),
    };

    session.verify_session(&token).map_err(|e| {
        log::debug!("Failed to verify token: {}", e);
        AuthError::InvalidSessionToken
    })?;

    Ok(())
}

pub async fn delete_expired_sessions(conn: &mut sqlx::SqliteConnection) -> anyhow::Result<()> {
    let current_time = Utc::now().timestamp();
    sqlx::query!("DELETE FROM session WHERE expired_at < ?", current_time)
        .execute(&mut *conn)
        .await?;

    Ok(())
}
