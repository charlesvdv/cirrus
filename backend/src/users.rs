use std::str::FromStr;

use anyhow::{bail, Result};
use argon2::{
    password_hash::{rand_core::OsRng, PasswordHasher, SaltString},
    Argon2, PasswordVerifier,
};

use super::role::Role;
use serde::{Deserialize, Serialize};

#[derive(thiserror::Error, Debug)]
pub enum UserError {
    #[error("User name is empty")]
    UserNameEmpty,
    #[error("User name too long")]
    UserNameTooLong,
}

#[derive(Serialize, Deserialize, sqlx::Type, Debug)]
#[sqlx(transparent)]
pub struct UserId(i64);

impl UserId {
    pub fn new(id: i64) -> UserId {
        UserId(id)
    }
}

pub struct User {
    id: UserId,
    name: String,
    password: PasswordHash,
    role: Role,
}

impl User {
    fn from_new_user(new_user: &NewUser) -> Result<Self> {
        if new_user.name.is_empty() {
            bail!(UserError::UserNameEmpty);
        }

        if new_user.name.len() > 64 {
            bail!(UserError::UserNameTooLong);
        }

        let role = new_user.role.as_ref().unwrap_or_else(|| {
            log::info!("defaulting to default role for user: {}", new_user.name);
            &Role::User
        });

        Ok(User {
            id: UserId(-1),
            name: new_user.name.clone(),
            password: new_user.password.hash()?,
            role: role.clone(),
        })
    }
}

#[derive(Deserialize, Debug)]
pub struct NewUser {
    pub name: String,
    pub password: Password,
    pub role: Option<Role>,
}

#[derive(Deserialize)]
pub struct Password(String);

impl Password {
    pub fn hash(&self) -> Result<PasswordHash> {
        let salt = SaltString::generate(&mut OsRng);
        let hash = Argon2::default()
            .hash_password(self.0.as_bytes(), &salt)?
            .to_string();

        Ok(PasswordHash(hash))
    }
}

impl std::fmt::Debug for Password {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(f, "*****")
    }
}

impl FromStr for Password {
    type Err = std::convert::Infallible;

    fn from_str(s: &str) -> Result<Self, Self::Err> {
        Ok(Password(s.into()))
    }
}

#[derive(sqlx::Type)]
#[sqlx(transparent)]
pub struct PasswordHash(String);

impl PasswordHash {
    pub fn from_str(pwd: String) -> PasswordHash {
        PasswordHash(pwd)
    }

    pub fn verify_password(&self, password: &str) -> Result<()> {
        let hash = argon2::PasswordHash::new(&self.0)?;
        Argon2::default().verify_password(password.as_bytes(), &hash)?;

        Ok(())
    }
}

pub async fn get(conn: &mut sqlx::SqliteConnection, id: &UserId) -> Result<Option<User>> {
    let record = sqlx::query!(
        r#"
        SELECT user.id, user.name, user.password, role.name as role FROM user
        INNER JOIN role
        ON user.id = role.id
        WHERE user.id = ?
    "#,
        id.0
    )
    .fetch_optional(&mut *conn)
    .await?;

    if let Some(record) = record {
        Ok(Some(User {
            id: UserId(record.id),
            name: record.name,
            password: PasswordHash(record.password),
            role: record.role.parse()?,
        }))
    } else {
        Ok(None)
    }
}

pub async fn create(conn: &mut sqlx::SqliteConnection, new_user: &NewUser) -> Result<User> {
    let mut user = User::from_new_user(new_user)?;

    let role_name = user.role.to_string();
    let result = sqlx::query!(
        "INSERT INTO user(name, password, role) VALUES(?, ?, (SELECT id FROM role WHERE name = ?))",
        user.name,
        user.password,
        role_name,
    )
    .execute(&mut *conn)
    .await?;

    user.id = UserId(result.last_insert_rowid());

    Ok(user)
}
