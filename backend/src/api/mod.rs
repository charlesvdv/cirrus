use axum::{
    Router,
    routing::get,
    http::StatusCode,
};

use self::ui::spa_router;

mod ui;

pub fn build_api_router(ui_assets_path: &std::path::Path) -> Router {
    let api_routes = Router::new()
        .route("/health", get(health_handler));

    Router::new()
        .merge(spa_router(ui_assets_path))
        .nest("/api", api_routes)
}



async fn health_handler() -> StatusCode {
    StatusCode::OK
}