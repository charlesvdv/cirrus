use axum::response::IntoResponse;
use hyper::StatusCode;

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

impl IntoResponse for Error {
    fn into_response(self) -> axum::response::Response {
        match self {
            Self::Anyhow(_) => {
                // TODO log error
            }
            Self::Sqlx(_) => {
                // TODO log error
            }
            Self::BadRequest(ref message) => {}
        };

        (self.status_code(), self.to_string()).into_response()
    }
}
