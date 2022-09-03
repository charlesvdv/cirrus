use axum::routing::{get, post};
use axum::Router;
use axum::{response::Json, Extension};
use sqlx::SqlitePool;

use crate::api::Result;
use crate::instance::InitInstance;
use crate::instance::{self, Instance};

pub fn router() -> Router {
    Router::new()
        .route("/instance", get(get_instance))
        .route("/instance/init", post(init_instance))
}

async fn get_instance(Extension(db_pool): Extension<SqlitePool>) -> Result<Json<Instance>> {
    let mut tx = db_pool.begin().await?;
    let instance = instance::get(&mut tx).await?;
    tx.commit().await?;

    Ok(Json(instance))
}

async fn init_instance(
    Extension(db_pool): Extension<SqlitePool>,
    Json(mut payload): Json<InitInstance>,
) -> Result<()> {
    let mut tx = db_pool.begin().await?;

    instance::init(&mut tx, &mut payload).await?;

    tx.commit().await?;

    Ok(())
}
