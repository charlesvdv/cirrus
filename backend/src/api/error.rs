use axum::{response::IntoResponse, Json};
use hyper::StatusCode;
use serde::Serialize;

#[derive(thiserror::Error, Debug)]
pub enum Error {
    #[error("an error occured with the database")]
    Sqlx(#[from] sqlx::Error),
    #[error("an internal server error occured")]
    Anyhow(#[from] anyhow::Error),
    #[error("{0}")]
    BadRequest(String),
}

impl Error {
    fn status_code(&self) -> StatusCode {
        match *self {
            Self::Sqlx(_) | Self::Anyhow(_) => StatusCode::INTERNAL_SERVER_ERROR,
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
            Self::Anyhow(ref err) => {
                log::error!("Generic error: {:?}", err);
            }
            Self::Sqlx(ref err) => {
                log::error!("Sqlx error: {:?}", err);
            }
        };

        (self.status_code(), self.to_string()).into_response()
    }
}
