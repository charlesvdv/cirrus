use anyhow::Result;
use argon2::{
    password_hash::{rand_core::OsRng, PasswordHasher, SaltString},
    Argon2,
};

use super::{IdentityError, Role};
use serde::Deserialize;

pub struct UserId(i64);

pub struct User {
    id: UserId,
    name: String,
    password: PasswordHash,
    role: Role,
}

impl User {
    fn from_new_user(new_user: &NewUser) -> Result<Self, IdentityError> {
        if new_user.name.is_empty() {
            return Err(IdentityError::UserNameEmpty);
        }

        if new_user.name.len() > 64 {
            return Err(IdentityError::UserNameTooLong);
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
    name: String,
    password: Password,
    role: Option<Role>,
}

impl NewUser {
    pub fn force_as_admin(&mut self) {
        self.role = Some(Role::Administrator);
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

pub async fn create_user(
    conn: &mut sqlx::SqliteConnection,
    new_user: &NewUser,
) -> Result<User, IdentityError> {
    let mut user = User::from_new_user(new_user)?;

    let role_name = user.role.to_string();
    let result = sqlx::query!(
        "INSERT INTO user(name, password, role) VALUES(?, ?, (SELECT id FROM role WHERE name = ?))",
        user.name,
        user.password.0,
        role_name,
    )
    .execute(&mut *conn)
    .await?;

    user.id = UserId(result.last_insert_rowid());

    Ok(user)
}
