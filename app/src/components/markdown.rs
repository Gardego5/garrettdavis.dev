use std::{
    collections::{HashMap, HashSet},
    fmt::Display,
    sync::OnceLock,
};

use gray_matter::{engine::YAML, Matter};
use maud::{PreEscaped, Render};
use pulldown_cmark::{html, Event, LinkType, Options, Parser, Tag, TagEnd};
use serde::de::DeserializeOwned;

const OPTIONS: Options = Options::empty()
    .union(Options::ENABLE_HEADING_ATTRIBUTES)
    .union(Options::ENABLE_SMART_PUNCTUATION)
    .union(Options::ENABLE_YAML_STYLE_METADATA_BLOCKS);

fn tag_attributes_config() -> &'static HashMap<&'static str, HashSet<&'static str>> {
    static HASHMAP: OnceLock<HashMap<&'static str, HashSet<&'static str>>> = OnceLock::new();
    HASHMAP.get_or_init(|| {
        let id_set = HashSet::from(["id"]);
        HashMap::from([
            ("code", HashSet::from(["class"])),
            ("h1", id_set.clone()),
            ("h2", id_set.clone()),
            ("h3", id_set.clone()),
            ("h4", id_set.clone()),
            ("h5", id_set.clone()),
            ("h6", id_set),
            ("a", HashSet::from(["href"])),
        ])
    })
}

pub struct Markdown<T>(T);

impl<T> Markdown<T> {
    fn h_anchor<'a: 'b, 'b>(id: &'a str, title: &'a str) -> [Event<'b>; 3] {
        [
            Event::Start(Tag::Link {
                link_type: LinkType::Inline,
                dest_url: format!("#{id}").into(),
                title: title.into(),
                id: "".into(),
            }),
            Event::Text("#".into()),
            Event::End(TagEnd::Link),
        ]
    }
}

impl<T: AsRef<str>> Render for Markdown<T> {
    fn render(&self) -> maud::Markup {
        let mut unsafe_html = String::new();
        let mut in_heading = Vec::new();
        let mut heading_text = String::new();
        let parser = Parser::new_ext(self.0.as_ref(), OPTIONS).filter_map(|event| match event {
            // heading links and auto ids
            Event::Start(Tag::Heading { .. }) => {
                in_heading.push(event);
                None
            }
            Event::End(TagEnd::Heading(_)) => match in_heading.first() {
                Some(Event::Start(Tag::Heading { id: Some(id), .. })) => {
                    let (open_tag, contents) = in_heading.split_at(1);
                    let mut result = String::new();
                    html::push_html(&mut result, open_tag.iter().map(ToOwned::to_owned));
                    html::push_html(&mut result, Self::h_anchor(id, &heading_text).into_iter());
                    html::push_html(&mut result, contents.iter().map(ToOwned::to_owned));
                    html::push_html(&mut result, std::iter::once(event)); // close tag

                    in_heading.clear();
                    heading_text.clear();

                    Some(Event::Html(result.into()))
                }
                Some(Event::Start(Tag::Heading {
                    level,
                    classes,
                    id: None,
                    ..
                })) => {
                    heading_text = heading_text
                        .chars()
                        .filter_map(|ch| match ch {
                            'a'..='z' | 'A'..='Z' | '0'..='9' => Some(ch),
                            ' ' | '.' => Some('-'),
                            _ => None,
                        })
                        .collect::<String>()
                        .to_lowercase();

                    let mut result = format!("<{level} id=\"{heading_text}\"");

                    if classes.is_empty() {
                        result.push('>');
                    } else {
                        result.push_str(" class=\"");
                        for class in classes {
                            result.push_str(class.as_ref());
                            result.push(' ');
                        }
                        result.push_str("\">");
                    }

                    result.push_str(format!("<a href=\"#{heading_text}\">#</a>").as_str());
                    html::push_html(
                        &mut result,
                        in_heading.iter().skip(1).map(ToOwned::to_owned),
                    );
                    result.push_str(format!("</{level}>").as_str());

                    in_heading.clear();
                    heading_text.clear();

                    Some(Event::Html(result.into()))
                }
                _ => unreachable!(
                    "Expected first event in heading to be heading open tag. Found: {:?}",
                    event
                ),
            },
            Event::Text(text) if in_heading.is_empty() == false => {
                heading_text.push_str(text.as_ref());
                in_heading.push(Event::Text(text));
                None
            }
            Event::Code(text) if in_heading.is_empty() == false => {
                heading_text.push_str(text.as_ref());
                in_heading.push(Event::Code(text));
                None
            }
            event if in_heading.is_empty() == false => {
                in_heading.push(event);
                None
            }

            _ => Some(event),
        });

        html::push_html(&mut unsafe_html, parser);
        let safe_html = ammonia::Builder::default()
            .tag_attributes(tag_attributes_config().to_owned())
            .clean(&unsafe_html)
            .to_string();

        PreEscaped(safe_html)
    }
}

pub struct RawDocument<T>(pub T);

#[derive(Debug, thiserror::Error)]
pub enum MarkdownWithFrontMatterError {
    FrontMatterMissing,
    FrontMatterInvalid(serde_json::error::Error),
}

impl Display for MarkdownWithFrontMatterError {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        match self {
            Self::FrontMatterMissing => write!(f, "missing frontmatter"),
            Self::FrontMatterInvalid(err) => write!(f, "invalid frontmatter {err}"),
        }
    }
}

pub struct MarkdownWithFrontMatter<M> {
    pub matter: M,
    pub markup: Markdown<String>,
}

impl<T: AsRef<str>, M: DeserializeOwned> TryFrom<RawDocument<T>> for MarkdownWithFrontMatter<M> {
    type Error = MarkdownWithFrontMatterError;
    fn try_from(value: RawDocument<T>) -> Result<Self, Self::Error> {
        let matter = Matter::<YAML>::new();
        let result = matter.parse(value.0.as_ref());

        Ok(Self {
            matter: result
                .data
                .ok_or(Self::Error::FrontMatterMissing)?
                .deserialize()
                .map_err(|err| Self::Error::FrontMatterInvalid(err))?,
            markup: Markdown(result.content),
        })
    }
}
