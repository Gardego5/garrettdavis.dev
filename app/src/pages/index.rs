use axum::response::IntoResponse;
use maud::html;

use crate::components::{
    layout::{margins, Header},
    template::template,
};

pub async fn handler() -> impl IntoResponse {
    template(
        html!(
            meta name="description" content="Garrett Davis is a young software developer who cares deaply about creating great software for people.";
            title { "Garrett Davis" }
        ),
        html!((Header::Floating(None))(margins(
            html!(main class="text-lg [&>*]:mb-2 [&>2]:last:mb-0" {
                p {
                    "hello! i'm garrett. "
                    "i'm learning a lot. "
                    "and i plan to keep at it. "
                }

                div {
                    "i love making software that is: "
                    ul ."inline md:inline-flex md:gap-5 items-center" {
                        @for (adjective, class) in [
                            ("enjoyable to use", "text-yellow-200"),
                            ("reliable", "text-red-300"),
                            ("maintainable", "text-blue-200"),
                            ("simple", "text-green-300")
                        ] { " " li class={
                            "relative after:absolute after:-right-3 after:top-1/2 after:-translate-y-1/2 "
                            "inline after:w-1 after:h-1 after:rounded-full after:bg-zinc-50 "
                            "after:content-none md:after:content-['_'] last:after:content-none " (class)
                        } { (adjective) } }
                    }
                    "."
                }

                div."flex items-center justify-center gap-4" {
                    div {
                        a target="_blank" rel="noreferrer noopener"
                            href="https://github.com/Gardego5"
                            aria-label="View Garrett's GitHub profile"
                            class="border border-slate-500 bg-zinc-950 border-dotted p-2 pl-2 sm:pl-0.5 flex text-baseline h-8 items-center text-sm gap-1 text-nowrap"
                        {
                            iconify-icon icon="mdi:github" width="32" height="32" class="sm:scale-75" {}
                            span."hidden sm:inline" { "github" }
                        }
                    }

                    span."text-sm text-center" {
                        span."text-nowrap" { "<--" }
                        " find me here "
                        span."text-nowrap" { "-->" }
                    }

                    div {
                        a target="_blank" rel="noreferrer noopener"
                            href="https://www.linkedin.com/in/garrett-davis-8793a721b/"
                            aria-label="View Garrett's LinkedIn Profile"
                            class="border border-slate-500 bg-zinc-950 border-dotted p-2 pr-2 sm:pr-0.5 flex text-baseline h-8 items-center text-sm gap-1 text-nowrap"
                        {
                            span."hidden sm:inline" { "linkedin" }
                            iconify-icon icon="mdi:linkedin" width="32" height="32" class="sm:scale-75" {}
                        }
                    }
                }

                hr;

                div {
                    "a blog thing"
                }
            })
        ))),
    )
}
