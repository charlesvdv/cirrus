use cirrus_backend::{Config, App};

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
    };

    let app = App::new(&config).unwrap();
    let port = app.port();

    let _ = tokio::spawn(app.run());

    TestApp {
        address: format!("http://127.0.0.1:{}", port),
    }
}
