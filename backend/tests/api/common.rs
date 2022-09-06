use axum::body::Body;
use axum::http::{request, HeaderValue, Request};
use axum::response::Response;
use cirrus_backend::instance::{self, InitInstance};
use cirrus_backend::users::NewUser;
use hyper::body;
use sqlx::SqlitePool;

pub trait RequestBuilderExt {
    fn json(self, json: serde_json::Value) -> Request<Body>;

    fn empty_body(self) -> Request<Body>;

    fn with_auth(self, token: &str) -> request::Builder;
}

impl RequestBuilderExt for request::Builder {
    fn json(self, json: serde_json::Value) -> Request<Body> {
        self.header("Content-Type", "application/json")
            .body(Body::from(json.to_string()))
            .expect("failed to build request")
    }

    fn empty_body(self) -> Request<Body> {
        self.body(Body::empty()).expect("failed to build request")
    }

    fn with_auth(mut self, token: &str) -> request::Builder {
        let headers = self.headers_mut().unwrap();
        headers.insert(
            "Authorization",
            HeaderValue::from_str(&format!("Bearer {}", token)).unwrap(),
        );
        self
    }
}

pub async fn to_json<T: serde::de::DeserializeOwned>(response: &mut Response) -> T {
    let bytes = body::to_bytes(response.body_mut()).await.unwrap();
    serde_json::from_slice(&bytes).unwrap()
}

pub const ADMIN: &str = "admin";
pub const ADMIN_PASSWORD: &str = "MyAdminPassword!23";

pub async fn init_test_instance(db: SqlitePool) -> anyhow::Result<()> {
    let mut conn = db.acquire().await?;

    instance::init(
        &mut conn,
        &mut InitInstance {
            admin: NewUser {
                name: String::from(ADMIN),
                password: ADMIN_PASSWORD.parse().unwrap(),
                role: None,
            },
        },
    )
    .await?;

    Ok(())
}
