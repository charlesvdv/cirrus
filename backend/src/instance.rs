use anyhow::{bail, Result};
use serde::{Deserialize, Serialize};

use crate::users::{self, NewUser};

const INSTANCE_ID: i64 = 0;

#[derive(thiserror::Error, Debug)]
pub enum InstanceError {
    #[error("Instance already initialized")]
    AlreadyInitialized,
}

/// Describe global information about the instance
#[derive(Default, Clone, Serialize, Deserialize)]
pub struct Instance {
    is_initialized: bool,
}

impl Instance {
    pub fn is_initialized(&self) -> bool {
        self.is_initialized
    }

    pub fn mark_as_initialized(&mut self) {
        self.is_initialized = true;
    }
}

#[derive(Deserialize, Debug)]
pub struct InitInstance {
    pub admin: NewUser,
}

pub async fn init(conn: &mut sqlx::SqliteConnection, input: &mut InitInstance) -> Result<Instance> {
    let mut instance = get(conn).await?;
    if instance.is_initialized() {
        bail!(InstanceError::AlreadyInitialized);
    }

    input.admin.force_as_admin();

    users::create(conn, &input.admin).await?;

    instance.mark_as_initialized();
    update(conn, &instance).await?;

    Ok(instance)
}

pub async fn get(conn: &mut sqlx::SqliteConnection) -> Result<Instance> {
    let maybe_record = sqlx::query!(
        "SELECT is_initialized FROM instance WHERE id = ?",
        INSTANCE_ID
    )
    .fetch_optional(&mut *conn)
    .await?;

    if let Some(record) = maybe_record {
        Ok(Instance {
            is_initialized: record.is_initialized == 1,
        })
    } else {
        let instance = Instance::default();
        sqlx::query!(
            "INSERT INTO instance (id, is_initialized) VALUES(?, ?)",
            INSTANCE_ID,
            instance.is_initialized
        )
        .execute(&mut *conn)
        .await?;
        Ok(instance)
    }
}

pub async fn update(conn: &mut sqlx::SqliteConnection, instance: &Instance) -> Result<()> {
    let is_initialized_int = if instance.is_initialized { 1 } else { 0 };
    sqlx::query!(
        "UPDATE instance SET is_initialized = ? WHERE id = ?",
        is_initialized_int,
        INSTANCE_ID
    )
    .execute(&mut *conn)
    .await?;

    Ok(())
}
