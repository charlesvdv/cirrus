use sqlx::SqlitePool;

pub const ADMIN: &str = "admin";
pub const ADMIN_PASSWORD: &str = "MyAdminPassword!23";

pub async fn init_test_instance(db: SqlitePool) -> anyhow::Result<()> {
    let mut conn = db.acquire().await?;

    Ok(())
}
