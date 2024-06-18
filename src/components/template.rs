use maud::{html, Markup, Render, DOCTYPE};

pub fn template<H: Render, B: Render>(head_content: H, body_content: B) -> Markup {
    html! { (DOCTYPE)
        html lang="en" {
            head {
                (head_content)

                // Meta Tags
                meta charset="utf-8";
                meta name="viewport" content="width=device-width, initial-scale=1";

                // Tailwind CSS
                link rel="stylesheet" href="/static/style.css";

                // Local Third-Party Scripts
                script src="/static/3p/js/htmx.min.js" {}
                // Even though we want to defer the execution of these scripts, we don't want to delay it's loading.
                link rel="preload" as="script" href="/static/3p/js/iconify-icon.min.js";
                link rel="preload" as="script" href="/static/3p/js/alpinejs-morph.min.js";
                link rel="preload" as="script" href="/static/3p/js/alpinejs.min.js";
                script defer src="/static/3p/js/iconify-icon.min.js" {}
                script defer src="/static/3p/js/alpinejs-morph.min.js" {}
                script defer src="/static/3p/js/alpinejs.min.js" {}

                link rel="stylesheet" href="static/fonts/IosevkaGarrettDavisDev.css";
            }

            body class="box-border bg-zinc-900 text-zinc-50" hx-boost="true" { (body_content) }
        }
    }
}
