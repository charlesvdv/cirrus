use axum::{response::Json, Extension};
use serde::Deserialize;
use sqlx::SqlitePool;

use crate::api::Error;
use crate::api::Result;
use crate::identity::{self, CreateUser};
use crate::instance::{self, Instance};

pub async fn instance_get(Extension(db_pool): Extension<SqlitePool>) -> Result<Json<Instance>> {
    let mut tx = db_pool.begin().await?;
    let instance = instance::get(&mut tx).await?;
    tx.commit().await?;

    Ok(Json(instance))
}

#[derive(Deserialize, Debug)]
pub struct InitInstance {
    admin: CreateUser,
}

pub async fn instance_init(
    Extension(db_pool): Extension<SqlitePool>,
    Json(payload): Json<InitInstance>,
) -> Result<()> {
    let mut tx = db_pool.begin().await?;

    let mut instance = instance::get(&mut tx).await?;
    if instance.is_initialized() {
        return Err(Error::BadRequest(String::from(
            "Instance is already initialized",
        )));
    }

    let mut admin = payload.admin;
    admin.force_as_admin();

    identity::create_user(&mut tx, &admin).await?;

    instance.mark_as_initialized();
    instance::update(&mut tx, &instance).await?;

    tx.commit().await?;

    Ok(())
}
