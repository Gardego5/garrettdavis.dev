---
title: Fenced Code Blocks as Frontmatter
author: Garrett Davis
live: false
createdAt: 2024-08-21
updatedAt: 2024-08-21
description: Markdown already has support for code blocks, why does our frontmatter data need to be treated differently?
---

## What I've been up to

So... I've been working on my personal site.
I've been taking way to long in getting it live, because architechture is _such_ a fun thing!
Really, I've just been procrastinating a bit, since I enjoy building the infrastructure a bit more than I
As I've been working on setting it up, I've made less progress in finishing it, but a great deal of headway learning things.

You can see [the source code for my site here](https://github.com/Gardego5/garrettdavis.dev). What I've settled on is this:

- a single rust app that server renders my site using [maud](https://maud.lambda.xyz/) for templating.
- built into a docker image using nix [dockerTools](https://ryantm.github.io/nixpkgs/builders/images/dockertools/).
- deployed on [fly.io](https://fly.io/).

Initially, I built something quite different:

- multiple rust AWS lambda functions targetting AWS Graviton, each server rendering their own pages.
- connected using API Gateway.
- with static content on AWS s3.
- distributed using AWS Cloudfront.

This really scratched the architechture design itch I had at the time.
I got it working, and it was cool as a proof of concept, but I ran into a few issues.
The biggest one is this: AWS Lambda isn't very fun to run locally.
There is [AWS SAM](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/serverless-sam-cli-using-invoke.html) which works with Rust (and other runtimes).
It ended up being a bit too slow for making a frontend application where you want to see changes relatively quickly.

I eventually simplified this and ended up settling on the system that I currently have.
It rebuilds the whole server after about 8 seconds instead of close to a minute as I was dealing with before.

## Markdown Frontmatter

My web site has a portion that is essentially a static site generator.
This converts markdown to html at page render time, serving it to you!

There are many different static site generators that fill the same basic role, turning markdown into html to serve to a reader.
Many of them also permit the inclusion of a _frontmatter_ section that describes metadata about the document.
A common way this is implemented is with a small block of yaml at the top of the document. Something like this:

```markdown
---
title: An Example Blog Post
---

# An example Blog Post

_The rest of the document goes here..._
```

## Markdown Fenced Code Blocks

While I've been working on this though... I realized markdown already has a system for embedding code.
If you've worked with markdown (especially if you have an interest in programming) you are likely aware of fenced code blocks.
[Commonmark](https://commonmark.org/) is a common <small>_haha_</small> format for markdown which defines a [specification for fenced code blocks](https://spec.commonmark.org/0.31.2/#fenced-code-blocks).
An example is below:

````markdown
```yaml :blog:
title: Fenced Code Blocks as Frontmatter
```
````

Note that the Commonmark specification allows for use of both tildes (`~`) or backticks (\`) to delimit the code block.
Additionally, it allows an _info string_ that follows the opening code block.
Usually, the first word of this info string is used to specify the language of the code block.

```rs
let parser =
    Parser::new_ext(&self.0.parsed_entity.content, Self::OPTIONS).filter_map(|event| {
        match event {
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
                        .into()
                ))
            }
            event @ Event::End(TagEnd::CodeBlock) => Some(event),

            event @ _ => Some(event),
        }
    });
```
