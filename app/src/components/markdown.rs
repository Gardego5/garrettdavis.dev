use std::{
    collections::{HashMap, HashSet},
    fmt::{Debug, Display},
    ops::Deref,
    str::FromStr,
};

use gray_matter::{
    engine::{Engine, YAML},
    Matter, ParsedEntity,
};
use maud::{PreEscaped, Render};
use pulldown_cmark::{html, CodeBlockKind, Event, LinkType, Options, Parser, Tag, TagEnd};
use serde::de::DeserializeOwned;

#[derive(Debug, thiserror::Error)]
pub enum DocumentError {
    FrontMatterMissing,
    FrontMatterInvalid(serde_json::error::Error),
}

impl Display for DocumentError {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        match self {
            Self::FrontMatterMissing => write!(f, "missing frontmatter"),
            Self::FrontMatterInvalid(err) => write!(f, "invalid frontmatter {err}"),
        }
    }
}

pub trait Frontmatter: DeserializeOwned {
    type Engine: Engine = YAML;
    fn prefix() -> String;
    fn tag() -> String {
        std::any::type_name::<Self::Engine>()
            .split("::")
            .last()
            .unwrap_or("")
            .to_lowercase()
    }
    fn delimiter() -> String {
        format!("```{} :{}:", Self::tag(), Self::prefix())
    }
}

fn matter<FM: Frontmatter>() -> Matter<FM::Engine> {
    let mut config = Matter::new();
    config.delimiter = FM::delimiter();
    config.close_delimiter = Some("```".into());
    config.excerpt_delimiter = Some("#".into());
    config
}

pub struct Document<M: Frontmatter> {
    pub metadata: M,
    pub parsed_entity: ParsedEntity,
}

struct DocumentMarkdown<'a, FM: Frontmatter>(&'a Document<FM>);

impl<FM: Frontmatter> FromStr for Document<FM> {
    type Err = DocumentError;
    fn from_str(s: &str) -> Result<Self, Self::Err> {
        let config = matter::<FM>();

        let parsed_entity = config.parse(s);
        Ok(Self {
            metadata: parsed_entity
                .data
                .clone()
                .ok_or(DocumentError::FrontMatterMissing)?
                .deserialize()
                .map_err(|err| DocumentError::FrontMatterInvalid(err))?,
            parsed_entity,
        })
    }
}

impl<FM: Frontmatter + Debug> Debug for Document<FM> {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        f.debug_struct("Document")
            .field("metadata", &self.metadata)
            .field("parsed_entity", &self.parsed_entity)
            .finish()
    }
}

impl<FM: Frontmatter> Document<FM> {
    pub fn md(&self) -> impl Render + '_ {
        DocumentMarkdown(self)
    }
}

impl<'a, FM: Frontmatter> DocumentMarkdown<'a, FM> {
    const OPTIONS: Options = Options::empty()
        .union(Options::ENABLE_HEADING_ATTRIBUTES)
        .union(Options::ENABLE_SMART_PUNCTUATION)
        .union(Options::ENABLE_YAML_STYLE_METADATA_BLOCKS);

    fn tag_attributes_config() -> HashMap<&'static str, HashSet<&'static str>> {
        let id_set = HashSet::from(["id"]);
        HashMap::from([
            ("code", HashSet::from(["class"])),
            ("pre", HashSet::from(["class"])),
            ("span", HashSet::from(["class"])),
            ("h1", id_set.clone()),
            ("h2", id_set.clone()),
            ("h3", id_set.clone()),
            ("h4", id_set.clone()),
            ("h5", id_set.clone()),
            ("h6", id_set),
            ("a", HashSet::from(["href"])),
        ])
    }

    fn h_anchor<'r: 'b, 'b>(id: &'r str, title: &'r str) -> [Event<'b>; 3] {
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

impl<'a, FM: Frontmatter> Render for DocumentMarkdown<'a, FM> {
    fn render(&self) -> maud::Markup {
        let mut unsafe_html = String::new();
        let mut in_heading = Vec::new();
        let mut heading_text = String::new();
        let mut codeblock_lang = Some(String::new());
        let parser =
            Parser::new_ext(&self.0.parsed_entity.content, Self::OPTIONS).filter_map(|event| {
                match event {
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
                            html::push_html(
                                &mut result,
                                Self::h_anchor(id, &heading_text).into_iter(),
                            );
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

                    // codeblock highlighting
                    Event::Start(Tag::CodeBlock(CodeBlockKind::Fenced(info_string))) => {
                        codeblock_lang =
                            Some(info_string.split(' ').next().unwrap_or("js").to_string());
                        Some(Event::Start(Tag::CodeBlock(CodeBlockKind::Fenced(
                            info_string,
                        ))))
                    }
                    Event::Text(text) if let Some(codeblock_lang) = codeblock_lang.take() => {
                        Some(Event::Html(
                            inkjet::Highlighter::new()
                                .highlight_to_string(
                                    inkjet::Language::from_token(codeblock_lang)
                                        .unwrap_or(inkjet::Language::Plaintext),
                                    &inkjet::formatter::Html,
                                    text,
                                )
                                .unwrap()
                                .into(),
                        ))
                    }
                    event @ Event::End(TagEnd::CodeBlock) => Some(event),

                    _ => Some(event),
                }
            });

        html::push_html(&mut unsafe_html, parser);
        let safe_html = ammonia::Builder::default()
            .tag_attributes(Self::tag_attributes_config().to_owned())
            .clean(&unsafe_html)
            .to_string();

        PreEscaped(safe_html)
    }
}

pub struct NestedDocument<P: Frontmatter, C: Frontmatter> {
    pub parent: Document<P>,
    pub children: Vec<Document<C>>,
}

impl<P: Frontmatter, C: Frontmatter> FromStr for NestedDocument<P, C> {
    type Err = DocumentError;
    fn from_str(s: &str) -> Result<Self, Self::Err> {
        let mut parent: Document<P> = s.parse()?;
        let child_delimiter = C::delimiter();
        let mut parts = parent.parsed_entity.content.split(&child_delimiter);
        if let Some(content) = parts.next() {
            let children: Vec<Document<C>> = parts
                .map(|child| format!("{child_delimiter}{child}").parse())
                .collect::<Result<_, _>>()?;
            parent.parsed_entity.content = content.into();
            Ok(Self { parent, children })
        } else {
            Err(DocumentError::FrontMatterMissing)
        }
    }
}

impl<P: Frontmatter, C: Frontmatter> Deref for NestedDocument<P, C> {
    type Target = Document<P>;
    fn deref(&self) -> &Self::Target {
        &self.parent
    }
}
