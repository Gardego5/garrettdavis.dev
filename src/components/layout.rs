use maud::{html, Markup, Render};

pub fn margins<T: Render>(content: T) -> Markup {
    html! { main ."p-4".m-auto."max-w-3xl" { (content) } }
}

pub fn header(title: &str) -> Markup {
    html! {
        header class="sticky top-0 z-10 overflow-hidden border-b border-slate-500 bg-zinc-950 pt-3" {
            h1 class="mx-auto max-w-3xl px-4 pb-1 font-mono text-xl md:py-1 md:text-3xl" { (title) }
            nav class="relative -left-1 mx-auto mb-1 flex max-w-3xl flex-shrink items-center gap-4 px-3 text-sm after:absolute after:-left-[100%] after:-z-10 after:h-[1px] after:w-[1000vw] after:bg-slate-500 md:text-base [&>*]:bg-zinc-950 [&>*]:px-2" {
                @for (href, name) in [
                    ("/", "home"),
                    ("/blog", "blog"),
                    ("/projects", "projects"),
                    ("/contact", "contact"),
                ] {
                    a href=(href) { (name) }
                }
            }
        }
    }
}
