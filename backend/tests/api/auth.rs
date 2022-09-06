use std::borrow::BorrowMut;

use axum::http::{Request, StatusCode};
use serde_json::json;
use tower::ServiceExt;

use cirrus_backend::{api::build_api_router, auth::Session, users::User};

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
    let session = common::to_json::<Session>(&mut resp).await;

    let mut resp = app
        .borrow_mut()
        .oneshot(
            Request::get("/api/users/current")
                .with_auth(session.token())
                .empty_body(),
        )
        .await
        .unwrap();

    assert_eq!(StatusCode::OK, resp.status());
    let user = common::to_json::<User>(&mut resp).await;
    assert_eq!(common::ADMIN, user.name());
}

#[sqlx::test]
async fn test_auth_unauthorized(db: sqlx::SqlitePool) {
    common::init_test_instance(db.clone()).await.unwrap();
    let mut app = build_api_router(std::path::Path::new(""), db);

    let resp = app
        .borrow_mut()
        .oneshot(
            Request::get("/api/users/current")
                .with_auth("invalid session")
                .empty_body(),
        )
        .await
        .unwrap();

    assert_eq!(StatusCode::UNAUTHORIZED, resp.status());
}
