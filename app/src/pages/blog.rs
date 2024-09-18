use std::{collections::HashMap, fs::DirEntry, str::FromStr};

use axum::{
    extract::{Path, State},
    response::IntoResponse,
};
use maud::html;
use serde::Deserialize;

use crate::components::{
    error::Result,
    layout,
    markdown::{Document, Frontmatter},
    template::template,
};

fn gather_files<P: AsRef<std::path::Path>>(
    root_path: P,
) -> anyhow::Result<HashMap<std::path::PathBuf, DirEntry>> {
    Ok(std::fs::read_dir(root_path)?
        .filter_map(std::result::Result::ok)
        .map(|dir_entry| {
            let path = dir_entry.path();
            if path.is_file() {
                Ok(vec![(path, dir_entry)])
            } else if path.is_dir() {
                gather_files(path).map(|m| m.into_iter().collect())
            } else {
                Err(anyhow::anyhow!("Not a file or directory: {path:?}"))
            }
        })
        .try_fold(HashMap::new(), |mut acc, item| -> anyhow::Result<_> {
            acc.extend(item?);
            Ok(acc)
        })?)
}

#[derive(Debug, Deserialize)]
#[serde(rename_all = "camelCase")]
struct BlogPostFrontmatter {
    title: String,
    //slug: String,
    author: String,
    live: bool,
    created_at: chrono::DateTime<chrono::offset::Utc>,
    updated_at: chrono::DateTime<chrono::offset::Utc>,
}
impl Frontmatter for BlogPostFrontmatter {
    fn prefix() -> String {
        "blog".into()
    }
}
type BlogPost = Document<BlogPostFrontmatter>;

pub async fn handle_get(
    State(config): State<std::sync::Arc<crate::state::config::State>>,
) -> Result<impl IntoResponse> {
    let file_path = format!("{}/blog", config.data_dir);

    let entries = gather_files(file_path)?
        .into_iter()
        .filter_map(|(path, _dir_entry)| {
            std::fs::read_to_string(&path)
                .map_err(anyhow::Error::from)
                .and_then(|content| BlogPost::from_str(&content).map_err(Into::into))
                .ok()
                .filter(|post| post.metadata.live)
                .map(|post| (path, post.metadata))
        })
        .collect::<HashMap<_, _>>();

    Ok(template(
        html!( title { "Garrett Davis" } ),
        html!((layout::Header::Floating(None))(layout::margins(html!(
            main class="grid gap-4" {
                @for (path, post) in entries {
                    div class="border border-slate-500 rounded p-4" {
                        h2 class="text-2xl" {
                            a href={"/blog/" (path.file_stem().unwrap().to_str().unwrap())} {
                                (post.title)
                            }
                        }

                        p { (post.author) }

                        table.grid.text-sm { tbody {
                            tr {
                                td { "created at:" }
                                td { time datetime=(post.created_at.to_rfc3339()) {
                                    (post.created_at.format("%d/%m/%Y"))
                                } }
                            }
                            tr {
                                td { "updated at:" }
                                td { time datetime=(post.updated_at.to_rfc3339()) {
                                    (post.updated_at.format("%d/%m/%Y"))
                                } }
                            }
                        } }
                    }
                }
            }
        )))),
    ))
}

pub async fn handle_get_slug(
    Path((slug,)): Path<(String,)>,
    State(config): State<std::sync::Arc<crate::state::config::State>>,
) -> Result<impl IntoResponse> {
    let file_path = format!("{}/blog/{}.md", config.data_dir, slug);

    let blog = match std::fs::read_to_string(file_path) {
        Err(err) => match err.kind() {
            std::io::ErrorKind::NotFound => Err(crate::components::error::Error::from(err)
                .with_status(axum::http::StatusCode::NOT_FOUND)),
            _ => return Err(err.into()),
        }?,
        Ok(content) => BlogPost::from_str(&content)?,
    };

    Ok(template(
        html!(
            title { (blog.metadata.title) }
            @if let Some(description) = blog.parsed_entity.excerpt.as_deref() {
                meta name="description" content=(description);
            }
        ),
        html!((layout::Header::Floating(Some(&blog.metadata.title)))(
            layout::margins(html!(main.markdown { (blog.md()) } ))
        )),
    ))
}
