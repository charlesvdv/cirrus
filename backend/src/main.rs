use cirrus_backend::App;
use cirrus_backend::Config;

#[tokio::main]
async fn main() {
    let config = Config {
        host: String::from("127.0.0.1"),
        port: 8000,
    };
    let app = App::new(&config);

    println!("Listening on {}:{}...", config.host, config.port);
    app.run().await.unwrap();
}
