use std::net::TcpListener;
use std::pin::Pin;

mod config;
mod api;

pub use config::Config;

pub struct App {
    port: u16,
    fut: Pin<Box<dyn std::future::Future<Output = Result<(), hyper::Error>> + Send>>,
}

impl App {
    pub fn new(config: &Config) -> anyhow::Result<App> {
        let router = api::build_api_router();

        let address = format!("{}:{}", config.host, config.port);
        let listener = TcpListener::bind(&address)?;
        let port = listener.local_addr().unwrap().port();
        let fut = axum::Server::from_tcp(listener)?
            .serve(router.into_make_service());

        Ok(App { port, fut: Box::pin(fut) })
    }

    pub fn port(&self) -> u16 {
        return self.port;
    }

    pub async fn run(self) -> anyhow::Result<()> {
        self.fut.await?;
        Ok(())
    }
}