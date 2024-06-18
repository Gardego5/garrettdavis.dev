use aws_lambda_events::apigw::{ApiGatewayV2httpRequest, ApiGatewayV2httpResponse};
use http::Method;
use lambda_runtime::{service_fn, Error, LambdaEvent};
use maud::html;

use garrettdavis_dev_cargo_lambda::{
    components::{
        layout::{header, margins},
        template::template,
    },
    HtmlResponse,
};

async fn handler(
    event: LambdaEvent<ApiGatewayV2httpRequest>,
) -> Result<ApiGatewayV2httpResponse, Error> {
    if event.payload.request_context.http.method != Method::GET {
        let mut resp = HtmlResponse::new(template(html! {}, html! { "Method Not Allowed" }));
        *resp.status_mut() = 405;
        return Ok(resp.into());
    }

    let resp = HtmlResponse::new(template(
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
