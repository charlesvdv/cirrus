use axum::{
    Router,
    routing::get,
    http::StatusCode, Extension,
};
use sqlx::SqlitePool;

use self::ui::spa_router;

mod ui;

pub fn build_api_router(ui_assets_path: &std::path::Path, db_pool: SqlitePool) -> Router {
    let api_routes = Router::new()
        .route("/health", get(health_handler))
        .layer(Extension(db_pool));

    Router::new()
        .merge(spa_router(ui_assets_path))
        .nest("/api", api_routes)
}



async fn health_handler() -> StatusCode {
    StatusCode::OK
}