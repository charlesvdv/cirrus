use cirrus_backend::{Config, App, DatabaseConfig};

pub struct TestApp {
    pub address: String,
}

impl TestApp {
}

pub async fn spawn_app() -> TestApp {
    let config = Config {
        host: String::from("127.0.0.1"),
        port: 0,
        ui_assets_path: String::from(""),
        database: DatabaseConfig {
            url: String::from(":memory:"),
        }
    };

    let app = App::new(&config).await.unwrap();
    let port = app.port();

    let _ = tokio::spawn(app.run());

    TestApp {
        address: format!("http://127.0.0.1:{}", port),
    }
}
