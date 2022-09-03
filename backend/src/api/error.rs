use axum::{response::IntoResponse, Json};
use hyper::StatusCode;
use serde::Serialize;

use crate::identity::IdentityError;

#[derive(thiserror::Error, Debug)]
pub enum Error {
    #[error("an error occured with the database")]
    Other(#[from] sqlx::Error),
    #[error("an internal server error occured")]
    Database(#[from] anyhow::Error),
    #[error("{0}")]
    BadRequest(String),
    #[error("Forbidden")]
    Forbidden,
}

impl From<IdentityError> for Error {
    fn from(err: IdentityError) -> Self {
        match err {
            IdentityError::Other(e) => Error::Database(e),
            IdentityError::Database(e) => Error::Other(e),
            IdentityError::InvalidSessionToken => Error::Forbidden,
            e => Error::BadRequest(e.to_string()),
        }
    }
}

impl Error {
    fn status_code(&self) -> StatusCode {
        match *self {
            Self::Other(_) | Self::Database(_) => StatusCode::INTERNAL_SERVER_ERROR,
            Self::Forbidden => StatusCode::FORBIDDEN,
            Self::BadRequest(_) => StatusCode::BAD_REQUEST,
        }
    }
}

#[derive(Serialize)]
struct ErrorResponse {
    message: String,
}

impl IntoResponse for Error {
    fn into_response(self) -> axum::response::Response {
        match self {
            Self::BadRequest(ref message) => {
                log::warn!("Bad request: {:?}", message);

                return (
                    StatusCode::BAD_REQUEST,
                    Json(ErrorResponse {
                        message: message.clone(),
                    }),
                )
                    .into_response();
            }
            Self::Database(ref err) => {
                log::error!("Generic error: {:?}", err);
            }
            Self::Other(ref err) => {
                log::error!("Sqlx error: {:?}", err);
            }
            Self::Forbidden => {
                log::debug!("Forbidden")
            }
        };

        (self.status_code(), self.to_string()).into_response()
    }
}
