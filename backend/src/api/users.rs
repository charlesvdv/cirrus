use crate::auth::{self, LoginUser, Session};

use super::Result;
use axum::{routing::post, Extension, Json, Router};
use sqlx::SqlitePool;

pub fn router() -> Router {
    Router::new().route("/users/login", post(login))
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
