use std::borrow::BorrowMut;

use axum::http::{Request, StatusCode};
use serde_json::{json, Value};
use tower::ServiceExt;

use cirrus_backend::instance::Instance;
use cirrus_backend::{api::build_api_router, auth::Session};

use crate::common::{self, RequestBuilderExt};

#[sqlx::test]
async fn test_auth_success(db: sqlx::SqlitePool) {
    common::init_test_instance(db.clone()).await.unwrap();

    let mut app = build_api_router(std::path::Path::new(""), db);

    let mut resp = app
        .borrow_mut()
        .oneshot(Request::post("/api/users/login").json(json!({
            "name": common::ADMIN,
            "password": common::ADMIN_PASSWORD,
        })))
        .await
        .unwrap();

    assert_eq!(StatusCode::OK, resp.status());
    let _ = common::to_json::<Session>(&mut resp).await;
}
