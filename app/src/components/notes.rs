use maud::{html, Render};
use serde::Deserialize;

use super::markdown::MarkdownWithFrontMatter;

#[derive(Debug, Deserialize)]
pub struct NoteFrontMatter {
    pub title: String,
    #[serde(with = "time::serde::iso8601")]
    pub date: time::OffsetDateTime,
}

pub type Note = MarkdownWithFrontMatter<NoteFrontMatter>;

impl Render for Note {
    fn render(&self) -> maud::Markup {
        html! { section class="relative grid gap-2 p-6 bg-zinc-950 after:content-[''] after:absolute after:-z-10 after:-inset-[1px] after:bg-[linear-gradient(135deg,theme(colors.slate.500)_20px,transparent_20px,transparent_calc(100%-20px),theme(colors.slate.500)_calc(100%-20px))]" {
            div class="flex justify-between items-center" {
                h2 class="text-xl text-bold" { (self.matter.title) }
                span class="italic text-sm" { (self.matter.date) }
            }

            div class="markdown" { (self.markup) }
        } }
    }
}
