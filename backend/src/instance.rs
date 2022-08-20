use serde::{Deserialize, Serialize};

const INSTANCE_ID: i64 = 0;

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

pub async fn get(conn: &mut sqlx::SqliteConnection) -> Result<Instance, sqlx::Error> {
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

pub async fn update(
    conn: &mut sqlx::SqliteConnection,
    instance: &Instance,
) -> Result<(), sqlx::Error> {
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
