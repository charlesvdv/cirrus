#[derive(Clone)]
pub struct Config {
    pub host: String,
    pub port: u16,
    pub ui_assets_path: String,
    pub database: DatabaseConfig,
}

#[derive(Clone)]
pub struct DatabaseConfig {
    pub url: String,
}