use maud::{html, Markup, Render, DOCTYPE};

pub fn template<H: Render, B: Render>(head: H, body: B) -> Markup {
    html!((DOCTYPE) html lang="en" {
        head {
            (head)

            // Meta Tags
            meta charset="utf-8";
            meta name="viewport" content="width=device-width, initial-scale=1";

            // Tailwind CSS
            link rel="stylesheet" href="/static/css/style.css";

            // Local Third-Party Scripts
            script src="/static/3p/js/htmx.2.0.1.min.js" {}
            @for src in [
                "/static/3p/js/iconify-icon.2.1.0.min.js",
                "/static/3p/js/alpinejs-morph.3.14.1.min.js",
                "/static/3p/js/alpinejs.3.14.1.min.js",
            ] {
                // Even though we want to defer the execution of these
                // scripts, we don't want to delay it's loading.
                link rel="preload" as="script" href=(src);
                script defer src=(src) {}
            }

            link rel="stylesheet" href="/static/fonts/truetype/IosevkaGarrettDavisDev.css";
        }

        body class="box-border bg-zinc-900 text-zinc-50" hx-boost="true" { (body) }
    })
}
