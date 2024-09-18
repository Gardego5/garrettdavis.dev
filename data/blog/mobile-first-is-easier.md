```yaml :blog:
title: Mobile First is Easier
slug: mobile-first-is-easier
author: Garrett Davis
live: true
createdAt: "2023-06-28T06:03:34Z"
updatedAt: "2023-06-28T06:03:34Z"
description: Why make your life harder, and your user's life worse?
```

I've been working on an interesting project where our primary audience is B2B
customers. It's a fun project, we're using a cutting edge technology stack. It
makes for a good experience for developers, and is motivating to work on.

We were using [scss](https://sass-lang.com/) for styling our app, but we
recently transitioned to [tailwind](https://tailwindcss.com/). My whole team
was new to using tailwind, and it's been fun to learn, but importantly, it's
been **quick** to learn, and has made development faster. It's been kind of an all
around win.

One great thing about tailwind is it's set up **mobile first**, meaning that the
default styles are for mobile, and you can then customize them for larger
screens. This makes it really easy to build your site to support small screens
right from the get-go.

## What that really means

When designing a website or application, you want to make a great experience for
your users. A big part of how effectively a user can interact with you site is
the layout that it has, some simple layouts work well regardless of screen size,
but the more complicated that you make it, the more things you try to fit on a
page, the less favorable it will be for mobile users.

There are more specific ways that a site can be mobile friendly, but just the
layout on it's own has a huge impact, just by virtue of not being able to fit
everything on the page!

Because tailwind uses smaller screen sizes as the **default**, it makes our job
easier to build the page to support small screen sizes to begin with.

---

Unless I'm careful, I don't do this by default, and I think **most** of us don't.
As developers, we have nice big screens to see what we're working on. I
usually have multiple windows open - a terminal, an editor, and my browser at a
minimum. Despite this - my browser still is pretty large:

```js
const width = window.innerWidth; // 1288
const height = window.innerHeight; // 945
```

Compare this to the size of a common phone:

```js
const width = window.innerWidth; // 414
const height = window.innerHeight; // 896
```

This means that as we interact with our site, as we develop it, mobile doesn't
matter! When you finish building your page, and switch to mobile view, aghh! So
often it's wrong... because we weren't looking at it all along. This poses a
problem, since when we
