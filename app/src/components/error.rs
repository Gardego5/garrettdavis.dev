use maud::html;

use crate::components::layout::{margins, Header};

use super::template::template;

pub type Result<T> = std::result::Result<T, Error>;

#[derive(Debug)]
pub struct Error {
    status: axum::http::StatusCode,
    inner: Option<anyhow::Error>,
    message: Option<maud::Markup>,
}

impl Error {
    pub fn with_message(self, message: maud::Markup) -> Self {
        Self {
            message: Some(message),
            ..self
        }
    }

    pub fn with_status(self, status: axum::http::StatusCode) -> Self {
        Self { status, ..self }
    }
}

impl Default for Error {
    fn default() -> Self {
        Self {
            status: axum::http::StatusCode::INTERNAL_SERVER_ERROR,
            inner: None,
            message: Some(html!(p { "An unknown error has occured." })),
        }
    }
}

impl<E> From<E> for Error
where
    E: Into<anyhow::Error>,
{
    fn from(value: E) -> Self {
        Self {
            status: axum::http::StatusCode::INTERNAL_SERVER_ERROR,
            inner: Some(value.into()),
            message: None,
        }
    }
}

impl axum::response::IntoResponse for Error {
    fn into_response(self) -> axum::response::Response {
        let head = html!();
        let body = html! {
            (Header::Floating("Error".into()))

            (margins(html! {
                img class="m-auto w-[75%] rounded-full" alt="latte with orchid"
                    src="/static/image/coffee_flower.jpg" {}

                div class="mt-4" {
                    @if let Some(msg) = self.message {
                        (msg)
                    } @else if let Some(err) = self.inner {
                        p { (err) }
                    } @else {
                        p { "An unknown error occurred." }
                    }
                }
            }))
        };

        (self.status, template(head, body)).into_response()
    }
}
