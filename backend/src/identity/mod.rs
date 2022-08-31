mod login;
mod role;
mod user;

pub use role::Role;
pub(crate) use user::Password;
pub use user::{create_user, NewUser, User};

#[derive(thiserror::Error, Debug)]
pub enum IdentityError {
    #[error("User name cannot be empty")]
    UserNameEmpty,
    #[error("User name is too long")]
    UserNameTooLong,
    #[error("User name or password is invalid")]
    UserNameOrPasswordInvalid,

    #[error("Unexpected database error")]
    Database(#[from] sqlx::Error),
    #[error("Unexpected error")]
    Other(#[from] anyhow::Error),
}
