{
  "Title": "My Best Collaboration Hack",
  "Published": "2013-03-07T16:00:00+01:00",
  "Updated": "2013-03-07T16:00:00+01:00"
}

As a quick follow-up to my [Open Source And
Responsibility](/2013/03/07/open-source-and-responsibility.html) post, I'd like
to share a small collaboration hack I came up with last year. Here it is:

Whenever somebody sends you a pull request, give them commit access to your
project.

While it may sound incredible stupid at first, using this strategy will allow
you to unleash the true power of Github. Let me explain. Over the past year, I
realized that I could not allocate enough time to my open source projects
anymore. And I'm not even talking about fixing people's issues for free, I'm
talking about pull request piling up.

So why wouldn't I simply merge them? Well, a lot of them were actually not good
enough. They were lacking tests, documentation, violated coding standards, or
were introducing new issues the contributor had not considered. I would often
spend the time explaining the issues, only to discover that the original
author was now lacking the time to make the changes. Of course I could have
made the adjustements myself, but that would often take as much time as if I had
done the patch from scratch, so I started neglecting most incoming pull requests
that couldn't be merged right away.

And then I came across the hack I mentioned above. I wish I could take credit
for designing it, but it really happened by coincidence. Somebody send a pull
request for a project I was no longer using myself, and I could see an issue
with it right away. However, since I no longer care for the project, and the
person sending the request did, I simply added him as a collaborator and said
something like this: "I don't have time to maintain this project anymore, so I
gave you commit access to make any changes you'd like. It would be nice to follow
the style used by the rest of the code and add some tests to this patch".

The result was pure magic. Within a few hours the same person who had just
submitted a rather mediocre patch, had now fixed things up and commited them.
This was highly unusual, so I started using the same strategy for a few other
small projects I was no longer interested in maintaining. And it worked, over
and over again. Of course, sometimes it wouldn't make a difference, but it was
clearly working a lot better than my previous approach.

This eventually lead to me using the same approach for my two most popular
projects, [node-mysql](https://github.com/felixge/node-mysql) and
[node-formidable](https://github.com/felixge/node-formidable), which are now
actively maintained by a bunch of amazing developers, writing much better code
than I ever received in the form of pull requests before.

So why does it work? Well, I have a few ideas. Once people have commit access,
they are no longer worried that their patch might go unmerged. They now have
the power to commit it themselves, causing them to put much more work into it.
Doing the actual commit and push also changes the sense of ownership. Instead
of handing over a diff to somebody else, they are now part of the project,
owning a small part of it.

But the magic does not stop here. In addition to their initial contribution
quality going up, I've observed many people continuing to help out with issues
and patches send by other users. This is of course fueled by Github notifying
every contributor on a repository of all activity on it.

Of course, people could also abuse this power, deleting code, introducing bugs,
or doing other bad things to the project. But I've not seen this happening so
far, and I'm also not worried about it. Git allows me to easily revert any
problematic changes, and there is almost nothing to gain from messing up one
of my projects.

So if you're maintaining an open source project, popular or not, I highly
recommend considering this approach, as it can truly turn one man projects into
small community projects.

Last but not least, I'd like to thank a few of the great people who have made
significant contributions to some of my projects recently:

* [Diogo Resende](https://github.com/dresende) for doing incredible work on
  node-mysql.
* [Oz Wolfe](https://github.com/CaptainOz) for contributing a connection pool
  to node-mysql.
* [Nate Lillich](https://github.com/NateLillich) for improving the connection pool
* [Sven Lito](https://github.com/svnlto) for fixing bugs and merging patches
  for node-formidable.
* [@egirshov](https://github.com/egirshov) for contributing many improvements to
  the node-formidable multipart parser.
* [Andrew Kelley](https://github.com/superjoe30) for also helping with fixing
  bugs and making improvements to node-formidable.
* [Mike Frey](https://github.com/mikefrey) for contributing JSON support to
  node-formidable.
* [Alex Indigo](https://github.com/alexindigo) for putting serious amounts of
  work into improving and maintaining node-form-data.
* Everybody else who I forgot to mention or who made smaller contributions.
