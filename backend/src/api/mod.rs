use axum::{
    http::StatusCode,
    routing::{get, post},
    Extension, Router,
};
use sqlx::SqlitePool;

use self::ui::spa_router;

mod error;
mod instance;
mod ui;

pub use error::Error;

pub type Result<T, E = Error> = std::result::Result<T, E>;

pub fn build_api_router(ui_assets_path: &std::path::Path, db_pool: SqlitePool) -> Router {
    let api_routes = Router::new()
        .route("/health", get(health_handler))
        .route("/instance", get(instance::instance_get))
        .route("/instance/init", post(instance::instance_init))
        .layer(Extension(db_pool));

    Router::new()
        .merge(spa_router(ui_assets_path))
        .nest("/api", api_routes)
}

async fn health_handler() -> StatusCode {
    StatusCode::OK
}
