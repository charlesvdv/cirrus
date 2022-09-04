use axum::{response::IntoResponse, Json};
use hyper::StatusCode;
use serde::Serialize;

use crate::{auth::AuthError, users::UserError};

#[derive(thiserror::Error, Debug)]
pub enum Error {
    #[error("an error occured with the database")]
    Internal(anyhow::Error),
    #[error("an internal server error occured")]
    Database(#[from] sqlx::Error),
    #[error("{0}")]
    BadRequest(anyhow::Error),
    #[error("Forbidden")]
    Forbidden(anyhow::Error),
}

impl From<anyhow::Error> for Error {
    fn from(err: anyhow::Error) -> Self {
        if err.is::<AuthError>() {
            return Error::Forbidden(err);
        }

        if err.is::<UserError>() {
            return Error::BadRequest(err);
        }

        if err.is::<sqlx::Error>() {
            return Error::Database(err.downcast().unwrap());
        }

        Error::Internal(err)
    }
}

impl Error {
    fn status_code(&self) -> StatusCode {
        match *self {
            Self::Internal(_) | Self::Database(_) => StatusCode::INTERNAL_SERVER_ERROR,
            Self::Forbidden(_) => StatusCode::FORBIDDEN,
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
            Self::BadRequest(ref err) => {
                log::warn!("Bad request: {:?}", err);

                return (
                    StatusCode::BAD_REQUEST,
                    Json(ErrorResponse {
                        message: format!("{}", err),
                    }),
                )
                    .into_response();
            }
            Self::Database(ref err) => {
                log::error!("Generic error: {:?}", err);
            }
            Self::Internal(ref err) => {
                log::error!("Sqlx error: {:?}", err);
            }
            Self::Forbidden(ref err) => {
                log::debug!("Forbidden: {:?}", err)
            }
        };

        (self.status_code(), self.to_string()).into_response()
    }
}
