use axum::{
    Router,
    http::{Response, Request, uri::Uri}, handler::Handler
};
use axum_extra::routing::SpaRouter;
use hyper::{client::{HttpConnector, Client}, Body};

static UI_DEV_SERVER_ENV: &str = "CIRRUS_UI_DEV_SERVER";
lazy_static::lazy_static! {
    static ref CLIENT: Client<HttpConnector, Body> = Client::<HttpConnector, Body>::new();
}

/// Route web UI request to the correct location.
pub fn spa_router(assets_dir: &std::path::Path) -> Router<Body> {
    if std::env::var(UI_DEV_SERVER_ENV).is_ok() {
        Router::new().fallback(reverse_proxy.into_service())
    } else {
        SpaRouter::new("/assets", assets_dir).into()
    }
}

async fn reverse_proxy(mut req: Request<Body>) -> Response<Body> {
    let path = req.uri().path();
    let path_query = req
        .uri()
        .path_and_query()
        .map(|v| v.as_str())
        .unwrap_or(path);

    let server = std::env::var(UI_DEV_SERVER_ENV).unwrap();
    let uri = format!("{}{}", server, path_query);

    println!("uri: {}", uri);

    *req.uri_mut() = Uri::try_from(uri).unwrap();

    CLIENT.request(req).await.unwrap()
}