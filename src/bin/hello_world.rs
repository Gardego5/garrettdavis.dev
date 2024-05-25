use aws_lambda_events::{
    apigw::{ApiGatewayProxyRequest, ApiGatewayProxyResponse},
    encodings::Body,
};
use base64::Engine;
use http::{header::CONTENT_TYPE, HeaderMap};
use lambda_runtime::{service_fn, Error, LambdaEvent};
use maud::Render;

use garrettdavis_dev_cargo_lambda::Base64Body;

async fn handler(
    event: LambdaEvent<ApiGatewayProxyRequest>,
) -> Result<ApiGatewayProxyResponse, Error> {
    let page = Base64Body(maud::html! { (maud::DOCTYPE) html {
        head {
            title {"Test"}
        }
        body {
            @if let Some(body) = event.payload.body {
                @let body = base64::engine::general_purpose::STANDARD.decode(body).unwrap();
                @let body = String::from_utf8(body).unwrap();
                p { "Body: " (body) }
            }
            table {
                caption { "query string parameters" }
                thead {
                    tr { th { "key" } th { "value" } }
                }
                tbody {
                    @for (k, v) in event.payload.query_string_parameters.iter() {
                        tr { td { (k) } td { (v) } }
                    }
                }
            }
            table {
                caption { "multi value query string parameters" }
                thead {
                    tr { th { "key" } th { "value" } }
                }
                tbody {
                    @for (k, v) in event.payload.multi_value_query_string_parameters.iter() {
                        tr { td { (k) } td { (v) } }
                    }
                }
            }
        }
    }});

    Ok(ApiGatewayProxyResponse {
        body: Some(Body::Text(page.render().into())),
        headers: {
            let mut headers = HeaderMap::new();
            headers.insert(CONTENT_TYPE, "text/html".parse()?);
            headers
        },
        status_code: 200,
        is_base64_encoded: true,
        multi_value_headers: HeaderMap::new(),
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
