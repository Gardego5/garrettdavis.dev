use aws_lambda_events::apigw::{ApiGatewayV2httpRequest, ApiGatewayV2httpResponse};
use http::{HeaderMap, Method};
use lambda_runtime::{service_fn, Error, LambdaEvent};
use maud::html;

use garrettdavis_dev_cargo_lambda::{
    components::{
        layout::{header, margins},
        markdown::RawDocument,
        notes::Note,
        template::template,
    },
    HtmlResponse,
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

async fn handler(
    event: LambdaEvent<ApiGatewayV2httpRequest>,
) -> Result<ApiGatewayV2httpResponse, Error> {
    if event.payload.request_context.http.method != Method::GET {
        return Ok(ApiGatewayV2httpResponse {
            status_code: 405,
            headers: HeaderMap::new(),
            multi_value_headers: HeaderMap::new(),
            body: None,
            is_base64_encoded: false,
            cookies: vec![],
        });
    }

    let note: Note = match RawDocument(DATA).try_into() {
        Ok(note) => note,
        Err(err) => {
            let mut resp = HtmlResponse::new(template(html! {}, html! { (err) }));
            *resp.status_mut() = 500;
            return Ok(resp.into());
        }
    };

    let resp = HtmlResponse::new(template(
        html! { title { "Garrett Davis" } },
        html! {
            (header("Garrett Davis"))
            (margins(html! {
                main {
                    (note)
                }
            }))
        },
    ));

    Ok(resp.into())
}

#[tokio::main]
async fn main() -> Result<(), Error> {
    tracing_subscriber::fmt()
        .with_max_level(tracing::Level::INFO)
        .with_target(false)
        .without_time()
        .init();

    lambda_runtime::run(service_fn(handler)).await
}
