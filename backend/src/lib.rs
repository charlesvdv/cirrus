use std::pin::Pin;
use std::{net::TcpListener, str::FromStr};

pub mod api;
mod config;
pub mod instance;

pub use config::*;
use sqlx::sqlite::{SqliteConnectOptions, SqlitePoolOptions};

pub struct App {
    port: u16,
    fut: Pin<Box<dyn std::future::Future<Output = Result<(), hyper::Error>> + Send>>,
}

impl App {
    pub async fn new(config: &Config) -> anyhow::Result<App> {
        let db_options =
            SqliteConnectOptions::from_str(&config.database.url)?.create_if_missing(true);
        let db_pool = SqlitePoolOptions::new().connect_with(db_options).await?;

        sqlx::migrate!().run(&db_pool).await?;

        let ui_assets_path = std::path::Path::new(&config.ui_assets_path);
        let router = api::build_api_router(ui_assets_path, db_pool);

        let address = format!("{}:{}", config.host, config.port);
        let listener = TcpListener::bind(&address)?;
        let port = listener.local_addr()?.port();
        let fut = axum::Server::from_tcp(listener)?.serve(router.into_make_service());

        Ok(App {
            port,
            fut: Box::pin(fut),
        })
    }

    pub fn port(&self) -> u16 {
        self.port
    }

    pub async fn run(self) -> anyhow::Result<()> {
        self.fut.await?;
        Ok(())
    }
}
