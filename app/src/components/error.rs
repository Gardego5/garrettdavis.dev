use http::StatusCode;

use super::template::template;

pub type Result<T> = std::result::Result<T, Error>;

pub struct Error {
    status: axum::http::StatusCode,
    inner: Option<anyhow::Error>,
}

impl Error {
    pub fn with_status<E: Into<anyhow::Error>>(status: axum::http::StatusCode, err: E) -> Self {
        Self {
            status,
            inner: Some(err.into()),
        }
    }
}

impl Default for Error {
    fn default() -> Self {
        Self {
            status: axum::http::StatusCode::INTERNAL_SERVER_ERROR,
            inner: None,
        }
    }
}

impl<E> From<E> for Error
where
    E: Into<anyhow::Error>,
{
    fn from(value: E) -> Self {
        Self {
            status: StatusCode::INTERNAL_SERVER_ERROR,
            inner: Some(value.into()),
        }
    }
}

impl axum::response::IntoResponse for Error {
    fn into_response(self) -> axum::response::Response {
        (
            self.status,
            template(
                maud::html! {},
                maud::html! {
                    main {
                        h1 { "An error has occurred" }
                        @if let Some(err) = self.inner {
                            p { (err) }
                        }
                    }
                },
            ),
        )
            .into_response()
    }
}
