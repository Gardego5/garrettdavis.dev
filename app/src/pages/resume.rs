use axum::response::IntoResponse;
use maud::{html, Render};

use crate::components::{error::Result, template::template};

enum ItemDesc {
    P(&'static str),
    Ul(Vec<&'static str>),
}

impl Render for ItemDesc {
    fn render(&self) -> maud::Markup {
        match self {
            ItemDesc::P(desc) => maud::html! { p class="text-sm" { (desc) } },
            ItemDesc::Ul(list_items) => maud::html! {
                ul class="flex gap-x-3 gap-y-2 flex-wrap text-sm" {
                    @for item in list_items { li { (item) } }
                }
            },
        }
    }
}

struct Item {
    title: &'static str,
    subtitle: Option<&'static str>,
    desc: ItemDesc,
    annotation: Option<&'static str>,
}

impl Render for Item {
    fn render(&self) -> maud::Markup {
        maud::html! {
            div class="pb-3" {
                h3 class="text-blue-300 text-xl" {
                    (self.title)

                    @if let Some(date) = &self.annotation {
                        span class="inline-block float-right" { (date) }
                    }
                }

                @if let Some(subtitle) = &self.subtitle {
                    h4 class="text-slate-300 text-sm pl-2 -mt-1 mb-1 italic print:text-xs" {
                        (subtitle)
                    }
                }

                (self.desc)
            }
        }
    }
}

struct Section {
    title: &'static str,
    items: Vec<Item>,
}

impl Render for Section {
    fn render(&self) -> maud::Markup {
        maud::html! {
            h2 class="font-mono leading-7 tracking-lighter " { (self.title) }
            div { @for item in &self.items { (item) } }
        }
    }
}

pub async fn handler() -> Result<impl IntoResponse> {
    let sections = vec![
        Section {
            title: "Technical Skills",
            items: vec![
                Item {
                    title: "Programming Languages",
                    subtitle: None,
                    desc: ItemDesc::Ul(vec![
                        "TypeScript",
                        "JavaScript",
                        "Nix",
                        "Rust",
                        "HCL",
                        "Python",
                        "Shell",
                        "PHP",
                        "Java"
                    ]),
                    annotation: None,
                },
                Item {
                    title: "Frontend",
                    subtitle: None,
                    desc: ItemDesc::Ul(vec![
                        "NextJS",
                        "React",
                        "Vue",
                        "Nuxt",
                        "HTML",
                        "CSS",
                        "TailwindCSS",
                        "Bootstrap",
                        "MaterialUI",
                        "SASS"
                    ]),
                    annotation: None,
                },
                Item {
                    title: "Libraries",
                    subtitle: None,
                    desc: ItemDesc::Ul(vec!["Express", "NextAuth", "TRPC", "Axum", "Zod", "esbuild"]),
                    annotation: None,
                },
                Item {
                    title: "Databases",
                    subtitle: None,
                    desc: ItemDesc::Ul(vec!["MySQL", "DynamoDB", "SQLite", "PostgreSQL", "IndexedDB"]),
                    annotation: None
                },
                Item {
                    title: "CI/CD Tools",
                    subtitle: None,
                    desc: ItemDesc::Ul(vec!["Terraform", "AWS CodeBuild", "AWS CodePipeline", "GitHub Actions", "NixOS"]),
                    annotation: None
                },
            ],
        },
        Section {
            title: "Personal Skills",
            items: vec![
                Item {
                    title: "Teaching",
                    subtitle: None,
                    desc: ItemDesc::P("I seek out new technologies and look for how they might benefit my current work. After I find something useful and applicable, I love to share this with my team - and strive to help my whole group excel. For example while working as a contractor at University of Phoenix, I volunteered and then lead a talk on building AWS Lambda functions using Rust."),
                    annotation: None,
                },
                Item {
                    title: "Leadership",
                    subtitle: None,
                    desc: ItemDesc::P("I managed a team of 3-7 baristas while working at Starbucks. I lead our store to be the top in our district for customer connection by focusing on technical excellence and meaningful conversations with our customers."),
                    annotation: None,
                }
            ]
        },
        Section {
            title: "Experience",
            items: vec![
                Item {
                    title: "University of Phoenix",
                    subtitle: Some("Software Engineer I"),
                    desc: ItemDesc::Ul(vec![
                                       "Spearheaded transition from typescript to golang for backend microservices",
                    ]),
                    annotation: Some("10/23 - Present"),
                },
                Item {
                    title: "Cook Systems - University of Phoenix",
                    subtitle: Some("Contract Software Engineer"),
                    desc: ItemDesc::Ul(vec![
                        "Revamped existing authentication flow built with PHP to integrate with NextAuth and custom Django / EdX authentication system.",
                        "I built automatically generated sitemap for marketing website.",
                        "I built a tool to programatically test lighthouse scores across entire marketing site.",
                        "I built a custom replacement for Terraform Cloud using s3 backend and codepipeline, reducing time to create a new project's cicd infrastrucure from roughly 8 hours down to only 15 - 30 minutes.",
                    ]),
                    annotation: Some("10/22 - 10/23"),
                },
                Item {
                    title: "MSR-FSR",
                    subtitle: Some("Production Technician"),
                    desc: ItemDesc::Ul(vec![
                        "Performed detail oriented work in a cleanroom environment.",
                    ]),
                    annotation: Some("10/21 - 02/22")
                },
                Item {
                    title: "Starbucks",
                    subtitle: Some("Shift Supervisor"),
                    desc: ItemDesc::Ul(vec![
                        "Mentored multiple baristas toward promotion to supervisor by encouraging them to coach others.",
                        "I created opportunities for training and used commendation to encourage individual growth.",
                        "Grew our  team by training 18+ new hires while instilling Starbucks quality and customer service values.",
                        "I lead coffee tastings focusing on the science and history in each cup.",
                        "Managed a team of 3-7 people, promoting communication and teamwork.",
                    ]),
                    annotation: Some("08/18 - 10/21"),
                },
            ],
        },
        Section {
            title: "Education",
            items: vec![
                Item {
                    title: "Cook Systems FastTrack'D",
                    subtitle: None,
                    desc: ItemDesc::P("Concentrated Java Frameworks and developer tools training."),
                    annotation: Some("07/22 - 09/22"),
                },
                Item {
                    title: "Portland Community College",
                    subtitle: None,
                    desc: ItemDesc::P("Associate of General Studies. Emphasis on GIS, Cartography, and Math."),
                    annotation: Some("09/16 - 03/20"),
                },
            ]
        },
    ];

    Ok(template(
        html! {
            title { "Resume - Garrett Davis" }
            meta name="description" content="Garrett Davis' resume";
        },
        html! {
            div class="fixed top-0 right-0 flex gap-2 p-2 print:hidden" {
                button x-data x-on:click="window.print()" title="Print this page." {
                    iconify-icon icon="ph:printer" width="36" height="36" {}
                }
            }

            div class="grid grid-cols-[1fr_6fr] [&>*:nth-child(odd)]:text-right gap-4 print:text-black max-w-4xl m-auto" {
                div class="col-start-2 flex py-4 items-center" {
                    h1 {
                        span class="text-4xl font-mono" { "Garrett Davis" }
                        br;
                        span class="text-lg" { "Resume" }
                    }
                }

                div class="col-start-2" {
                    "I am a software engineer, with a passion for maintainable software, teaching, and learning new technologies. "
                }

                @for section in sections { (section) }
            }
        },
    ))
}
