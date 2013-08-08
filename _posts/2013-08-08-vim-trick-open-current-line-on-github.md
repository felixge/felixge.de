---
layout: post
title: "Vim Trick: Open current line on GitHub"
date: 2013-08-08T20:50:00+01:00
updated: 2013-08-08T20:50:00+01:00
---

In my never-ending quest to automate my work flow, I recently came up with a
neat little vim trick I'd like to share.

One of the things I do quite frequently, is pasting links to GitHub files/lines
in email and chat conversations, e.g. [Protocol.js#L144][]. So far my workflow
for this has been to navigate to the GitHub repo, hit "t" to select the file,
click on the line I want, and then copy the address bar location into my
clipboard.

After doing this several times in a row the other day, I decided to automate
it.  The first piece of the puzzle was to create a git alias for determining
the GitHub URL of the current repository to put into my ~/[.gitconfig][] file:

```bash
[alias]
	url =! bash -c 'git config --get remote.origin.url | sed -E "s/.+:\\(.+\\)\\.git$/https:\\\\/\\\\/github\\\\.com\\\\/\\\\1/g"'
```

This allows me to easily get the URL of the current repository I am in, e.g.:

```bash
$ git url
https://github.com/felixge/node-mysql
```

Now I can do cool things like quickly opening this URL in my browser:

```bash
$ git url | xargs open
# or
$ open `git url`
```

But that still requires me to manually navigate to the file / line I am
currently interested in, and I'm lazy. So I came up with this key binding for
my ~/[.vimrc][]:

```bash
nnoremap <leader>o :!echo `git url`/blob/`git rev-parse --abbrev-ref HEAD`/%\#L<C-R>=line('.')<CR> \| xargs open<CR><CR>
```

Now I can simply press ",o" ("," is my leader key), and my browser will
automatically navigate to the file/line underneath my cursor.

Feel free to adopt this for your editor / environment and let me know if you
make any improvements. For example, one thing I didn't get around to yet is
opening visual selections as line ranges.

[.gitconfig]: https://github.com/felixge/dotfiles/blob/master/.gitconfig:
[Protocol.js#L144]: https://github.com/felixge/node-mysql/blob/master/lib/protocol/Protocol.js#L144
[.vimrc]: https://github.com/felixge/dotfiles/blob/master/.vimrc
