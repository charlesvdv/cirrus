mod helpers;

#[tokio::test]
async fn app_starts_correctly() {
    let app = helpers::spawn_app().await;

    let response = reqwest::get(&format!("{}/api/health", app.address))
        .await
        .expect("Failed to execute health check request");

    assert!(response.status().is_success());
}


