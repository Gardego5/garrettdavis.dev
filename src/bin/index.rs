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

    let page = template(
        html! {
            meta name="description" content="Garrett Davis is a young software developer who cares deaply about creating great software for people.";
            title { "Garrett Davis" }
        },
        html! {
            (header("Garrett Davis"))
            (margins(html! {
                main {
                    p {
                        "hello! i'm garrett."
                    }

                    h2 { "some things i'm passionate about creating in my work" }
                    ul class="list-disc" {
                        @for (short, long) in [
                            ("a good user experience", html! {
                                "there's a lot that goes into creating a good user experience. "
                                "if we're not careful as developers " span.italic { "(+ team)" }
                                ", we could forget that they are who we're really building a project for. "
                            }),
                            ("reliable software", html! {
                                "part of what makes something a good user experience is " span.italic { "when it doesn't break. " }
                                "in order to do this, i do a few things: "
                            }),
                            ("maintainable software", html! {
                                "if you want your project to be useful long term, it's also important that it can be effectively maintained. "
                                "some things that contribute to maintainable projects are "
                                span.italic { "keeping your eye simple, " }
                                span.italic { "following accepted standards. " }
                            }),
                        ] {
                            li {
                                span class="italic font-bold" { (short) }
                                " " (long)
                            }
                        }
                    }

                    div class="flex items-center justify-center gap-4" {
                        @for (href, icon, label) in [
                            ("https://github.com/Gardego5", "mdi:github", "View Garrett's GitHub profile"),
                            ("https://www.linkedin.com/in/garrett-davis-8793a721b/", "mdi:linkedin", "View Garrett's LinkedIn Profile"),
                        ] { a target="_blank" rel="noreferrer noopener" href=(href) aria-label=(label)
                            class="border border-slate-500 bg-zinc-950 border-dotted aspect-square p-2 grid place-content-center" {
                            iconify-icon icon=(icon) width="28" height="28" {}
                        } }
                    }
                }
            }))
        },
    );

    Ok(HtmlResponse::new(page).into())
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
