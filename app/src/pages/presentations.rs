use std::io;

use anyhow::anyhow;
use axum::{
    extract::{Path, State},
    response::IntoResponse,
};
use http::StatusCode;
use maud::{html, PreEscaped, Render};
use serde::{Deserialize, Serialize};
use tailwind_fuse::tw_merge;

use crate::components::{
    error::{Error, Result},
    markdown::{MarkdownWithFrontMatter, RawDocument},
    template::template,
};

const DELIMITER: &'static str = "===slide===\n";

pub async fn handler(
    State(config): State<std::sync::Arc<crate::state::config::State>>,
    Path((presenation,)): Path<(String,)>,
) -> Result<impl IntoResponse> {
    let file_path = format!("{}/presentations/{presenation}.md", config.data_dir);
    let str = std::fs::read_to_string(file_path).map_err(|err| match err.kind() {
        io::ErrorKind::NotFound => Error::default()
            .with_status(StatusCode::NOT_FOUND)
            .with_message(html!(p class="max-w-xs mx-auto text-center italic" {
                "It looks like I haven't made any presentation with that name..."
            })),
        _ => Error::from(err),
    })?;

    let slides = str
        .split(DELIMITER)
        .map(RawDocument)
        .map(Slide::try_from)
        .collect::<std::result::Result<Vec<_>, _>>()?;

    let slides_length = slides.len();

    Ok(template(
        html!(
            meta name="description" content="";
            style { (PreEscaped(include_str!("./presentations.rs.css"))) }
            script { (PreEscaped(include_str!("./presentations.rs.js"))) }
            title { }
        ),
        html!(main
            x-data={"{"
                "slide:0,"
                "islide:'1',"
                "modal:false,"
                "max:" (slides_length) ","
                "increment() { this.slide < this.max && this.slide++ },"
                "decrement() { this.slide > 0 && this.slide-- },"
            "}"}
            "@keydown.arrow-down.window.prevent"="increment"
            "@keydown.arrow-up.window.prevent"="decrement"
            x-effect="window.document.getElementById(slide.toString())?.scrollIntoView()" {

            @for (idx, slide) in slides.into_iter().enumerate() {
                div #(idx+1) class="relative h-screen" {
                    (slide)

                    button class="absolute top-8 left-8 text-xs italic p-1" "@click"="$refs.modal.showModal()" {
                        (idx+1) "/" (slides_length)
                    }

                    @if idx > 0 {
                        button class="w-6 h-6 text-slate-500 absolute right-8 top-8"
                            "@click"="decrement" {
                            iconify-icon icon="ph:caret-up" width="24" height="24" {}
                        }
                    }

                    @if idx+1 < slides_length {
                        button class="w-6 h-6 text-slate-500 absolute right-8 bottom-8"
                            "@click"="increment" {
                            iconify-icon icon="ph:caret-down" width="24" height="24" {}
                        }
                    }
                }
            }

            dialog class="bg-black/50 p-3 rounded-lg" x-ref="modal" {
                form method="dialog" {
                    input class="bg-transparent text-slate-300 text-center text-lg"
                        "@beforeinput"={"(NUMBERS_ONLY.test($event.data) &&"
                            "!($event.target.value+$event.data) ||"
                            "inRange(+$event.target.value+$event.data, 0, " (slides_length) " ) )"
                            "|| $event.preventDefault()"}
                        autofocus type="text" x-model="slide" {}
                }
            }
        }),
    ))
}

#[derive(Debug, Clone, Copy, Serialize, Deserialize, Default)]
enum SlideLayout {
    #[default]
    Default,
    Centered,
}

impl SlideLayout {
    fn class(&self) -> &'static str {
        match self {
            Self::Default => "",
            Self::Centered => "justify-center items-center text-center",
        }
    }
}

#[derive(Debug, Serialize, Deserialize)]
struct SlideMetadata {
    title: String,
    layout: Option<SlideLayout>,
}
type Slide = MarkdownWithFrontMatter<SlideMetadata>;

impl Render for Slide {
    fn render(&self) -> maud::Markup {
        let class = tw_merge!(
            "p-20 inset-4 absolute border border-slate-500 slide grid grid-rows-[auto_minmax(0,1fr)] gap-16",
            self.matter.layout.unwrap_or_default().class()
        );

        html!(section class=(class) {
            h2 { (self.matter.title) }

            div class="grid gap-8" {
                (self.markup)
            }
        })
    }
}
