use dotenv::dotenv;
use clap::{Subcommand, Parser};

use cirrus_backend::App;
use cirrus_backend::Config;

#[derive(Parser)]
#[clap(name = "cirrus", version, about)]
#[clap(propagate_version = true)]
struct Cli {
    #[clap(subcommand)]
    command: Option<Commands>,

    /// Specify port exposed
    #[clap(short, long)]
    port: Option<u16>,
}

#[derive(Subcommand)]
enum Commands {
    /// Initialize cirrus
    Init{},
}

#[tokio::main]
async fn main() {
    dotenv().ok();

    let cli = Cli::parse();

    let config = Config {
        host: String::from("127.0.0.1"),
        port: cli.port.unwrap_or(8000),
        ui_assets_path: String::from("/path"),
    };

    match &cli.command {
        Some(Commands::Init{}) => init(&config).await,
        None => run(&config).await,
    }

}

async fn init(config: &Config) {

}

async fn run(config: &Config) {
    let app = App::new(&config).unwrap();

    println!("Listening on {}:{}...", config.host, config.port);
    app.run().await.unwrap();
}