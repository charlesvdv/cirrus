use dotenv::dotenv;

use cirrus_backend::App;
use cirrus_backend::Config;

#[tokio::main]
async fn main() {
    dotenv().ok();

    let config = Config {
        host: String::from("127.0.0.1"),
        port: 8000,
        ui_assets_path: String::from("/path"),
    };

    let app = App::new(&config).unwrap();

    println!("Listening on {}:{}...", config.host, config.port);
    app.run().await.unwrap();
}
