use std::borrow::BorrowMut;

use axum::http::{Request, StatusCode};
use tower::ServiceExt;

use cirrus_backend::api::build_api_router;
use cirrus_backend::instance::Instance;

mod common;
use common::RequestBuilderExt;

#[sqlx::test]
async fn test_instance_get(db: sqlx::SqlitePool) {
    let mut app = build_api_router(std::path::Path::new(""), db);

    let mut resp = app
        .borrow_mut()
        .oneshot(Request::get("/api/instance").empty_body())
        .await
        .unwrap();

    assert_eq!(resp.status(), StatusCode::OK);
    let instance = common::to_json::<Instance>(&mut resp).await;
    assert!(!instance.is_initialized());
}
