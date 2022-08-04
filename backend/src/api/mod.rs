use axum::{Router, routing::get, http::StatusCode};

pub fn build_api_router() -> Router {
    let api_routes = Router::new()
        .route("/health", get(health_handler));

    Router::new().nest("/api", api_routes)
}

async fn health_handler() -> StatusCode {
    StatusCode::OK
}