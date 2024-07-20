use axum::response::IntoResponse;
use maud::html;

use crate::components::{
    layout::{header, margins},
    template::template,
};

pub async fn handler() -> impl IntoResponse {
    template(
        html! { title { "Garrett Davis" } },
        html! {
            (header("Garrett Davis"))
            (margins(html! {
                main {
                    h2 class="font-mono text-3xl flex flex-col items-end" {
                        span class="self-start" { "Idea!" }
                        span class="" { "Code." }
                        span class="" { "Create." }
                        span class="" { "Deploy." }
                        span class="" { "Repeat." }
                    }
                }
            }))
        },
    )
}
