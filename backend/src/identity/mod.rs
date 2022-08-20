use anyhow::Result;
use argon2::{
    password_hash::{rand_core::OsRng, PasswordHasher, SaltString},
    Argon2,
};

use role::Role;
use serde::Deserialize;

mod role;

pub struct UserId(i64);

pub struct User {
    id: UserId,
    name: String,
    password: PasswordHash,
    role: Role,
}

#[derive(Deserialize, Debug)]
pub struct CreateUser {
    name: String,
    password: Password,
    role: Option<Role>,
}

impl CreateUser {
    pub fn force_as_admin(&mut self) {
        self.role = Some(Role::Admin);
    }
}

#[derive(Deserialize)]
struct Password(String);

impl Password {
    fn hash(&self) -> Result<PasswordHash> {
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

struct PasswordHash(String);

pub async fn get_user(conn: &mut sqlx::SqliteConnection, id: &UserId) -> Result<Option<User>> {
    let record = sqlx::query!(
        r#"
        SELECT id, name, password, role FROM user
        WHERE id = ?
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

pub async fn create_user(conn: &mut sqlx::SqliteConnection, user: &CreateUser) -> Result<User> {
    let role = user.role.as_ref().unwrap_or(&Role::User);
    let role_str = role.to_string();
    let hashed_password = user.password.hash()?;

    let result = sqlx::query!(
        "INSERT INTO user(name, password, role) VALUES(?, ?, ?)",
        user.name,
        hashed_password.0,
        role_str
    )
    .execute(&mut *conn)
    .await?;

    Ok(User {
        id: UserId(result.last_insert_rowid()),
        name: user.name.clone(),
        password: hashed_password,
        role: role.clone(),
    })
}
