use crate::{
    auth::{self, LoginUser, Session},
    users::{self, User, UserId},
};

use super::{AuthChecker, Result};
use anyhow::anyhow;
use axum::{
    routing::{get, post},
    Extension, Json, Router,
};
use sqlx::SqlitePool;
use tower_http::auth::AsyncRequireAuthorizationLayer;

pub fn router(db_pool: SqlitePool) -> Router {
    let private_router = Router::new()
        .route("/users/current", get(current_user))
        .layer(AsyncRequireAuthorizationLayer::new(AuthChecker::new(
            db_pool,
        )));

    Router::new()
        .route("/users/login", post(login))
        .merge(private_router)
}

async fn login(
    Extension(db_pool): Extension<SqlitePool>,
    Json(payload): Json<LoginUser>,
) -> Result<Json<Session>> {
    let mut tx = db_pool.begin().await?;
    let session = auth::authenticate(&mut tx, payload).await?;
    tx.commit().await?;

    Ok(Json(session))
}

#[axum_macros::debug_handler]
async fn current_user(
    Extension(db_pool): Extension<SqlitePool>,
    Extension(user_id): Extension<UserId>,
) -> Result<Json<User>> {
    let mut conn = db_pool.acquire().await?;
    let user = users::get(&mut conn, &user_id)
        .await?
        .ok_or_else(|| anyhow!("Failed to retrieve users"))?;

    Ok(Json(user))
}
