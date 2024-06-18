use aws_lambda_events::apigw::{ApiGatewayV2httpRequest, ApiGatewayV2httpResponse};
use http::Method;
use lambda_runtime::{service_fn, Error, LambdaEvent};
use maud::{html, Markup, DOCTYPE};

use garrettdavis_dev_cargo_lambda::{
    components::{
        layout::{header, margins},
        template::template,
    },
    HtmlResponse,
};
use serde::Deserialize;

const COOKIE_NAME: &'static str = "contact-form-sumbitted";

async fn handle_get(
    event: LambdaEvent<ApiGatewayV2httpRequest>,
) -> Result<ApiGatewayV2httpResponse, Error> {
    if let Some(payload) = event
        .payload
        .cookies
        .unwrap_or(vec![])
        .iter()
        .map(|cookie| cookie.strip_prefix(format!("{COOKIE_NAME}=").as_str()))
        .flatten()
        .take(1)
        .collect::<Vec<&str>>()
        .first()
    {
        let message: Message = serde_html_form::from_str(payload).unwrap_or_default();
        let resp = HtmlResponse::new(template(
            html!(),
            html! {
                (header("Garrett Davis"))
                (margins(message.thanks_message()))
            },
        ));
        return Ok(resp.into());
    }

    let main = html! {
        form class="relative grid gap-2 rounded border border-slate-500 bg-gray-800 p-4 md:grid-cols-3"
            hx-post="/contact" hx-swap="outerHTML"
        {
            h2 class="text-2xl font-semibold md:col-span-4" { "I'd love to hear from you!" }
            input class="px-2 py-1 rounded border bg-zinc-900 border-slate-500"
                required type="text" name="name" placeholder="Name";
            input class="px-2 py-1 rounded border bg-zinc-900 border-slate-500 col-span-2"
                required type="email" name="email" placeholder="Email";
            textarea class="px-2 py-1 rounded border bg-zinc-900 border-slate-500 col-span-full min-h-[20vh] leading-5 md:col-span-4"
                required name="message" placeholder="Send me a message :)" {}
            button class="absolute -bottom-1.5 right-8 rounded border border-slate-500 bg-zinc-900 px-4 py-1 text-sm hover:bg-slate-800"
                type="submit" { "Send" };
        }
    };

    let resp = HtmlResponse::new(template(
        html!(),
        html! {
            (header("Contact Garrett"))
            (margins(main))
        },
    ));

    Ok(resp.into())
}

#[derive(Debug, PartialEq, Deserialize, Default)]
struct Message {
    pub name: Option<String>,
    pub email: Option<String>,
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
                        p { "Name: " (self.name.as_deref().unwrap_or("Missing")) }
                        p { "Email: " @match &self.email {
                            Some(email) => a href={"mailto:" (email)} { (email) },
                            None => span { "email missing" },
                        } }
                        p { "Message: " (self.message) }
                    }
                }
            }
        }
    }

    fn thanks_message(&self) -> Markup {
        html! {
            div class="rounded border border-slate-500 bg-gray-800 p-4 text-center" {
                "Thanks for reaching out"
                @if let Some(name) = &self.name { ", " (name) }
                "! I'll get back to you soon :)"
            }
        }
    }
}

async fn handle_post(
    event: LambdaEvent<ApiGatewayV2httpRequest>,
) -> Result<ApiGatewayV2httpResponse, Error> {
    let payload = event.payload.body.unwrap_or(String::new());
    let message: Message = match serde_html_form::from_str(&payload) {
        Ok(message) => message,
        Err(err) => {
            let mut resp = HtmlResponse::new(template(
                html! {title {"Error"}},
                html! {
                    div class="text-red-500" { (err) }
                },
            ));
            *resp.status_mut() = 400;
            return Ok(resp.into());
        }
    };

    let cookie = format!("{COOKIE_NAME}={payload}");

    if event
        .payload
        .headers
        .get("hx-request")
        .map(|val| val == "true")
        .unwrap_or(false)
    {
        let mut resp = HtmlResponse::new(message.thanks_message());
        resp.cookies_mut().push(cookie);
        Ok(resp.into())
    } else {
        let mut resp = HtmlResponse::new(template(
            html! { title { "Thanks! I'll be in touch soon" } },
            html! {
                (header("Contact Garrett"))
                (margins(message.thanks_message()))
            },
        ));
        resp.cookies_mut().push(cookie);
        Ok(resp.into())
    }
}

async fn handler(
    event: LambdaEvent<ApiGatewayV2httpRequest>,
) -> Result<ApiGatewayV2httpResponse, Error> {
    match event.payload.request_context.http.method {
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
