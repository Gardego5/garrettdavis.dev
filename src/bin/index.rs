use aws_lambda_events::{
    apigw::{ApiGatewayV2httpRequest, ApiGatewayV2httpResponse},
    encodings::Body,
};
use http::{header::CONTENT_TYPE, HeaderMap};
use lambda_runtime::{service_fn, Error, LambdaEvent};
use maud::{html, Render};

use garrettdavis_dev_cargo_lambda::{components::template::template, Base64Body};

async fn handler(
    _event: LambdaEvent<ApiGatewayV2httpRequest>,
) -> Result<ApiGatewayV2httpResponse, Error> {
    let head = html! {};
    let body = html! {
        h1 { "Hi there!" }
    };

    let page = Base64Body(template(head, body));

    Ok(ApiGatewayV2httpResponse {
        status_code: 200,
        headers: {
            let mut headers = HeaderMap::new();
            headers.insert(CONTENT_TYPE, "text/html".parse()?);
            headers
        },
        multi_value_headers: HeaderMap::new(),
        body: Some(Body::Text(page.render().into())),
        is_base64_encoded: true,
        cookies: vec![],
    })
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
