use cirrus_backend::DatabaseConfig;
use clap::Parser;
use dotenv::dotenv;

use cirrus_backend::App;
use cirrus_backend::Config;

#[derive(Parser)]
#[clap(name = "cirrus", version, about)]
#[clap(propagate_version = true)]
struct Cli {
    /// Specify port exposed
    #[clap(short, long)]
    port: Option<u16>,
}

#[tokio::main]
async fn main() {
    dotenv().ok();

    let cli = Cli::parse();

    let config = Config {
        host: String::from("127.0.0.1"),
        port: cli.port.unwrap_or(8000),
        ui_assets_path: String::from("/path"),
        database: DatabaseConfig {
            url: std::env::var("DATABASE_URL").unwrap_or(String::from("cirrus.sqlite")),
        }
    };

    let app = App::new(&config).await.unwrap();

    println!("Listening on {}:{}...", config.host, config.port);
    app.run().await.unwrap();
}
