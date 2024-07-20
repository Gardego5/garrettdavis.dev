use axum::response::IntoResponse;
use maud::html;

use crate::components::{
    error::Result,
    layout::{header, margins},
    markdown::RawDocument,
    notes::Note,
    template::template,
};

const DATA: &'static str = r#"---
title: Test Note.
date: 2024-06-07T23:39:26.705Z
---

# This is a simple test.

Hopefully, we will even have a bit of formatting.

<a href="/">Home</a>

<button>Hi!</button>
"#;

pub async fn handler() -> Result<impl IntoResponse> {
    let note: Note = RawDocument(DATA).try_into()?;

    Ok(template(
        html! { title { "Garrett Davis" } },
        html! {
            (header("Garrett Davis"))
            (margins(html! {
                main {
                    (note)
                }
            }))
        },
    ))
}
