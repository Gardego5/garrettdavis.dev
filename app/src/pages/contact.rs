use axum::{response::IntoResponse, Form};
use axum_extra::extract::CookieJar;
use http::{HeaderMap, StatusCode};
use maud::{html, Markup, DOCTYPE};

use serde::{Deserialize, Serialize};

use crate::components::{
    error::{Error, Result},
    layout::{header, margins},
    template::template,
};

const COOKIE_NAME: &'static str = "contact-form-sumbitted";

pub async fn handle_get(jar: CookieJar) -> impl IntoResponse {
    let head = html! { title { "Contact Garrett" } };
    let header = header("Garrett Davis");

    if let Some(cookie) = jar.get(COOKIE_NAME) {
        let message: Message = serde_html_form::from_str(cookie.value()).unwrap_or_default();
        return template(
            head,
            html! {
                (header)
                (margins(message.thanks_message()))
            },
        );
    }

    template(
        html! { title { "Contact Garrett" } },
        html! {
            (header)
            (margins(html! {
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
            }))
        },
    )
}

#[derive(Debug, PartialEq, Serialize, Deserialize, Default)]
pub struct Message {
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

#[axum::debug_handler]
pub async fn handle_post(
    headers: HeaderMap,
    jar: CookieJar,
    body: Form<Message>,
) -> Result<impl IntoResponse> {
    Ok((
        jar.add(
            serde_html_form::to_string(&body.0)
                .map_err(|err| Error::with_status(StatusCode::INTERNAL_SERVER_ERROR, err))?,
        ),
        if headers
            .get("hx-request")
            .map(|val| val.as_bytes() == b"true")
            .unwrap_or(false)
        {
            body.thanks_message()
        } else {
            template(
                html! { title { "Thanks! I'll be in touch soon" } },
                html! {
                    (header("Contact Garrett"))
                    (margins(body.thanks_message()))
                },
            )
        },
    ))
}
