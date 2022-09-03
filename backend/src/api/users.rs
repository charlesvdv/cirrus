use super::Result;
use crate::identity::login::{authenticate, LoginUser, Session};
use axum::{Extension, Json};
use sqlx::SqlitePool;

pub async fn login(
    Extension(db_pool): Extension<SqlitePool>,
    Json(payload): Json<LoginUser>,
) -> Result<Json<Session>> {
    let mut tx = db_pool.begin().await?;
    let session = authenticate(&mut tx, payload).await?;
    tx.commit().await?;

    Ok(Json(session))
}
