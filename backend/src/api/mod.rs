use axum::{
    http::StatusCode,
    routing::{get, post},
    Extension, Router,
};
use axum_extra::routing::SpaRouter;
use sqlx::SqlitePool;
use tower_http::trace::TraceLayer;

mod error;
mod instance;

pub use error::Error;

pub type Result<T, E = Error> = std::result::Result<T, E>;

pub fn build_api_router(ui_assets_path: &std::path::Path, db_pool: SqlitePool) -> Router {
    let api_routes = Router::new()
        .route("/health", get(health_handler))
        .route("/instance", get(instance::instance_get))
        .route("/instance/init", post(instance::instance_init))
        .layer(Extension(db_pool))
        .layer(TraceLayer::new_for_http());

    Router::new()
        .merge(SpaRouter::new("/assets", ui_assets_path))
        .nest("/api", api_routes)
}

async fn health_handler() -> StatusCode {
    StatusCode::OK
}
