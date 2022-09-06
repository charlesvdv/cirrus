use anyhow::{anyhow, bail, Result};
use futures_util::future::BoxFuture;
use hyper::{Request, Response, StatusCode};
use sqlx::SqlitePool;
use tower_http::auth::AsyncAuthorizeRequest;

use crate::{auth, users::UserId};

#[derive(Clone)]
pub struct AuthChecker {
    db_pool: SqlitePool,
}

impl AuthChecker {
    pub fn new(db_pool: SqlitePool) -> Self {
        AuthChecker { db_pool }
    }
}

impl<B> AsyncAuthorizeRequest<B> for AuthChecker
where
    B: Send + Sync + 'static,
{
    type RequestBody = B;
    type ResponseBody = axum::body::BoxBody;
    type Future = BoxFuture<'static, Result<Request<B>, Response<Self::ResponseBody>>>;

    fn authorize(&mut self, mut request: hyper::Request<B>) -> Self::Future {
        let db_pool = self.db_pool.clone();
        Box::pin(async {
            match check_authorization(db_pool, &request).await {
                Ok(user_id) => {
                    request.extensions_mut().insert(user_id);
                    Ok(request)
                }
                Err(err) => {
                    log::debug!("Failed to authorize request: {:?}", err);
                    let response = Response::builder()
                        .status(StatusCode::UNAUTHORIZED)
                        .body(axum::body::BoxBody::default())
                        .unwrap();
                    Err(response)
                }
            }
        })
    }
}

async fn check_authorization<B>(
    db_pool: SqlitePool,
    request: &hyper::Request<B>,
) -> Result<UserId> {
    let token = extract_session_token(request)?;

    let mut conn = db_pool.acquire().await?;
    auth::verify_session_token(&mut conn, token).await
}

fn extract_session_token<B>(request: &hyper::Request<B>) -> Result<&str> {
    if let Some(auth_header) = request.headers().get("Authorization") {
        if let Ok(value) = auth_header.to_str() {
            parse_authorize_value(value)
                .ok_or_else(|| anyhow!("Invalid Authorization header authentication scheme"))
        } else {
            bail!("Invalid Authorization header")
        }
    } else {
        bail!("Authorization header not found in request")
    }
}

fn parse_authorize_value(header_value: &str) -> Option<&str> {
    header_value.strip_prefix("Bearer ")
}
