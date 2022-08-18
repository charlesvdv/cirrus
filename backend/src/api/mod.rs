use axum::{http::StatusCode, response::Json, routing::get, Extension, Router};
use sqlx::SqlitePool;

use crate::instance::{self, Instance};

use self::ui::spa_router;

mod error;
mod ui;

pub use error::Error;

pub type Result<T, E = Error> = std::result::Result<T, E>;

pub fn build_api_router(ui_assets_path: &std::path::Path, db_pool: SqlitePool) -> Router {
    let api_routes = Router::new()
        .route("/health", get(health_handler))
        .route("/instance", get(instance_get))
        .layer(Extension(db_pool));

    Router::new()
        .merge(spa_router(ui_assets_path))
        .nest("/api", api_routes)
}

async fn health_handler() -> StatusCode {
    StatusCode::OK
}

async fn instance_get(Extension(db_pool): Extension<SqlitePool>) -> Result<Json<Instance>> {
    let mut tx = db_pool.begin().await?;
    let instance = instance::get(&mut tx).await?;
    tx.commit().await?;

    Ok(Json(instance))
}
