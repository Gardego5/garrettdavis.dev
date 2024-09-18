use std::io;

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
    layout::{margins, Header},
    markdown::{Document, Frontmatter, NestedDocument},
    template::template,
};

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

    let pres: NestedDocument<Presentation, Slide> = str.parse()?;

    Ok(template(
        html!(
            meta name="description" content="";
            style { (PreEscaped(include_str!("./presentations.rs.css"))) }
            script { (PreEscaped(include_str!("./presentations.rs.js"))) }
            title { }
        ),
        pres,
    ))
}

#[derive(Debug, Clone, Copy, Serialize, Deserialize, Default)]
#[serde(tag = "layout")]
enum Slide {
    #[default]
    Default,
    Centered,
}

impl Frontmatter for Slide {
    fn prefix() -> String {
        "slide".into()
    }
}

impl Render for Document<Slide> {
    fn render(&self) -> maud::Markup {
        static CLASS: &'static str = "p-20 inset-4 absolute border border-slate-500 slide grid grid-rows-[auto_minmax(0,1fr)] gap-16";
        let layout_class = match self.metadata {
            Slide::Centered => "justify-center items-center text-center",
            _ => "",
        };

        let content = match self.metadata {
            _ => html!(div class="grid gap-8" { (self.md()) }),
        };

        html!(div class=(tw_merge!(CLASS, layout_class)) { (content) })
    }
}

#[derive(Debug, Serialize, Deserialize)]
struct Presentation {}
impl Frontmatter for Presentation {
    fn prefix() -> String {
        "presentation".into()
    }
}

impl Render for NestedDocument<Presentation, Slide> {
    fn render(&self) -> maud::Markup {
        let slides_length = self.children.len();

        html!(main
            x-data={"{"
                "slide:0,"
                "islide:'1',"
                "modal:false,"
                "max:" (slides_length) ","
                "increment() { this.slide < this.max && this.slide++ },"
                "decrement() { this.slide > 1 && this.slide-- },"
            "}"}

            // hotkeys
            "@keydown.arrow-down.window.prevent"="increment"
            "@keydown.arrow-up.window.prevent"="decrement"
            "@keydown.space.window.prevent"="increment"
            "@keydown.home.window.prevent"="slide=1"
            "@keydown.end.window.prevent"={"slide=" (slides_length)}

            x-effect="window.document.getElementById(slide.toString())?.scrollIntoView()" {

            div class="transition-all delay-500 duration-1000" style="min-height: 60vh;" x-init="$el.style.minHeight = '100vh'" {
                (Header::Static(None))

                (margins(html!(header class="markdown" {
                    (self.md())
                })))
            }

            @for (idx, slide) in self.children.iter().enumerate() {
                section #(idx+1) class="relative h-screen" {
                    (slide)

                    button class="absolute top-8 left-8 text-slate-500 text-xs italic p-1"
                        "@click"="$refs.modal.showModal()" {
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
        })
    }
}
