use aws_lambda_events::{
    apigw::{ApiGatewayV2httpRequest, ApiGatewayV2httpResponse},
    encodings::Body,
};
use http::{header::CONTENT_TYPE, HeaderMap, Method};
use lambda_runtime::{service_fn, Error, LambdaEvent};
use maud::{html, Markup, Render, DOCTYPE};

use garrettdavis_dev_cargo_lambda::{
    components::{
        layout::{header, margins},
        template::template,
    },
    Base64Body,
};
use serde::Deserialize;

async fn handle_get(
    _event: LambdaEvent<ApiGatewayV2httpRequest>,
) -> Result<ApiGatewayV2httpResponse, Error> {
    let main = html! {
        form class="relative grid gap-2 rounded border border-slate-500 bg-gray-800 p-4 md:grid-cols-4"
            hx-post="/contact" hx-swap="outerHTML"
        {
            h2 class="text-2xl font-semibold md:col-span-4" { "I'd love to hear from you!" }
            input class="px-2 py-1 rounded border bg-zinc-900 border-slate-500"
                required type="text" name="first" placeholder="First Name";
            input class="px-2 py-1 rounded border bg-zinc-900 border-slate-500"
                required type="text" name="last" placeholder="Last Name";
            input class="px-2 py-1 rounded border bg-zinc-900 border-slate-500 col-span-2"
                required type="email" name="email" placeholder="Email";
            textarea class="px-2 py-1 rounded border bg-zinc-900 border-slate-500 col-span-full min-h-[20vh] leading-5 md:col-span-4"
                required name="message" placeholder="Send me a message :)" {}
            button class="absolute -bottom-1.5 right-8 rounded border border-slate-500 bg-zinc-900 px-4 py-1 text-sm hover:bg-slate-800"
                type="submit" { "Send" };
        }
    };

    let page = Base64Body(template(
        html!(),
        html! {
            (header("Contact Garrett"))
            (margins(main))
        },
    ));

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

#[derive(Debug, PartialEq, Deserialize)]
struct Message {
    pub first: String,
    pub last: String,
    pub email: String,
    pub message: String,
}

impl Message {
    fn automated_email(&self) -> Markup {
        html! {
            (DOCTYPE)
            html lang="en" {
                head {
                    meta http-equiv="Content-Type" content="text/html; charset=utf-8";
                    meta name="viewport" content="width=device-width";
                    title { "New message from contact form" }
                }
                body {
                    div id="email" style="width: 600px" {
                        p { "You have a new message from the contact form on garrettdavis.dev!" }
                        p { "Name: " (self.first) " " (self.last) }
                        p {
                            "Email: "
                            a href={"mailto:" (self.email)} {
                                (self.email)
                            }
                        }
                        p { "Message: " (self.message) }
                    }
                }
            }
        }
    }

    fn thanks_message(&self) -> Markup {
        html! {
            div class="rounded border border-slate-500 bg-gray-800 p-4 text-center" {
                "Thanks for reaching out, " (self.first) "! I'll get back to you soon :)"
            }
        }
    }
}

async fn handle_post(
    event: LambdaEvent<ApiGatewayV2httpRequest>,
) -> Result<ApiGatewayV2httpResponse, Error> {
    let message: Message =
        match serde_html_form::from_str(&event.payload.body.unwrap_or(String::new())) {
            Ok(message) => message,
            Err(err) => {
                return Ok(ApiGatewayV2httpResponse {
                    status_code: 400,
                    body: Some(Body::Text(
                        template(
                            html! {title {"Error"}},
                            html! {
                                div class="text-red-500" { (err) }
                            },
                        )
                        .render()
                        .into(),
                    )),
                    ..Default::default()
                })
            }
        };

    let page = Base64Body(template(
        html! { title { "Thanks! I'll be in touch soon" } },
        html! {
            (header("Contact Garrett"))
            (margins(message.thanks_message()))
        },
    ));

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

async fn handler(
    event: LambdaEvent<ApiGatewayV2httpRequest>,
) -> Result<ApiGatewayV2httpResponse, Error> {
    match event.payload.http_method {
        Method::GET => handle_get(event).await,
        Method::POST => handle_post(event).await,
        _ => unreachable!("this endpoint is only attached on GET and POST methods"),
    }
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
