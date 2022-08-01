mod config;
mod api;

pub use config::Config;

pub struct App {
    config: Config,
}

impl App {
    pub fn new(config: &Config) -> App {
        App{
            config: config.clone(),
        }
    }

    pub async fn run(self) -> anyhow::Result<()> {
        let router = api::build_router();

        let address = format!("{}:{}", self.config.host, self.config.port);
        axum::Server::bind(&address.parse().unwrap())
            .serve(router.into_make_service())
            .await?;
        Ok(())
    }
}