use maud::{html, Markup, Render};

pub fn margins<T: Render>(content: T) -> Markup {
    html!(div class="p-4 m-auto max-w-3xl" { (content) })
}

pub enum Header<'a> {
    Floating(Option<&'a str>),
    Static(Option<&'a str>),
}

impl Render for Header<'_> {
    fn render(&self) -> Markup {
        html!(header
            ."overflow-hidden border-b border-slate-500 bg-zinc-950 pt-3 top-0 z-10"
            .(if let Header::Floating(_) = self { "sticky" } else { "relative" }) {

            @let header_text_class = "mx-auto max-w-3xl px-4 pb-1 md:py-1";
            @if let Self::Floating(Some(title)) | Self::Static(Some(title)) = self {
                div .(header_text_class)."relative flex items-start justify-between gap-4" {
                    h1 ."text-xl md:text-3xl" { (title) }
                    a ."md:absolute md:left-[calc(100%-1rem)] pr-2 lg:left-2 md:top-0 md:-translate-x-full text-right text-sm lg:text-base" href="/" {
                        "Garrett" br; "Davis"
                    }
                }
            } @else {
                div .(header_text_class)."text-xl md:text-3xl" { a href="/" { "Garrett Davis" } }
            }

            nav ."relative -left-1 mx-auto mb-1 flex justify-end max-w-3xl flex-shrink items-center gap-4 px-3 text-sm "
                ."after:absolute after:-left-[100%] after:-z-10 after:h-[1px] after:w-[1000vw] after:bg-slate-500 md:text-base [&>*]:bg-zinc-950 [&>*]:px-2"
            {
                @for (href, name) in [
                    ("/blog", "blog"),
                    ("/contact", "contact"),
                    ("/notes", "notes"),
                ] { a href=(href) { (name) } }
            }
        })
    }
}
