use chrono::{DateTime, Duration, Utc};
use rand::distributions::{Alphanumeric, DistString};
use serde::{Deserialize, Serialize};

use super::{
    user::{PasswordHash, UserId},
    IdentityError, Password,
};

#[derive(Deserialize, Debug)]
pub struct LoginUser {
    name: String,
    password: String,
}

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
}

pub async fn login(
    conn: &mut sqlx::SqliteConnection,
    login: LoginUser,
) -> Result<Session, IdentityError> {
    let user = sqlx::query!(
        "SELECT id, name, password FROM user WHERE name = ?",
        login.name
    )
    .fetch_optional(&mut *conn)
    .await?
    .ok_or_else(|| IdentityError::UserNameOrPasswordInvalid)?;

    let user_password = PasswordHash::from_str(user.password);
    if user_password.verify_password(&login.password).is_err() {
        return Err(IdentityError::UserNameOrPasswordInvalid);
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
