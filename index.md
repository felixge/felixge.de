---
layout: default
title: About
---
Hi, I am Felix Geisendörfer, a programmer and entrepreneur living in Berlin,
Germany.

Currently I'm writing software in [Go](http://golang.org/) as a contractor.

In the past I've co-founded and bootstrapped
[Transloadit](http://transloadit.com/) into a profitable business, was one of
the first contributors to [node.js](http://nodejs.org/) and worked on many
other [open source projects](http://github.com/felixge) as well.

One of my spare time passions is robotics, so [a few
friends](http://nodecopter.com/core) and I organized an event around
programming flying robots with JavaScript called
[NodeCopter](http://nodecopter.com/) which has now turned into a small
community.

You can find me on [twitter](https://twitter.com/felixge) and
[github](https://github.com/felixge), or [contact me](#contact) directly.

<h2 id="blog"><a href="#blog">Blog</a></h2>

You can subscribe to updates via <a
  href="http://feeds.feedburner.com/felixge">RSS</a> or <a
  href="http://feedburner.google.com/fb/a/mailverify?uri=felixge">E-Mail</a>.

<table class="toc">
  <tbody>
    {% for post in site.posts %}
      {% unless post.hidden %}
      <tr>
        <td class="title">
          <span>
            <a href="{{ post.url }}">{{ post.title }}</a>
          </span>
        </td>
        <td class="date"><span>{{ post.date | date: "%b %d, %Y" }}</span></td>
      </tr>
      {% endunless %}
    {% endfor %}
  </tbody>
</table>

I have been blogging since 2006, my older posts can be found
[here](http://debuggable.com/posts/archive).

<h2 id="speaking"><a href="#speaking">Speaking</a></h2>

After speaking a little too much in 2012, I've heavily cut down on events, but
I still enjoy speaking ocassionally, so [reach out](#contact) if you'd like to
invite me to an event.

<table class="toc">
  <tbody>
    {% for talk in site.data.talks %}
    <tr>
      <td class="title">
        <span>
          <a href="{{ talk.url }}">{{ talk.title }}</a>
          {% if talk.pdfUrl %}
          &middot; <a href="{{ talk.pdfUrl }}">pdf</a>
          {% endif %}
          {% if talk.videoUrl %}
          &middot; <a href="{{ talk.videoUrl }}">video</a>
          {% endif %}
          {% if talk.codeUrl %}
          &middot; <a href="{{ talk.codeUrl }}">code</a>
          {% endif %}
        </span>
      </td>
      <td class="location"><span><a href="{{ talk.eventUrl }}">{{ talk.location }}</a></span></td>
      <td class="date"><span>{{ talk.date | date: "%b %d, %Y" }}</span></td>
    </tr>
    {% endfor %}
  </tbody>
</table>

<h2 id="personal"><a href="#personal">Personal</a></h2>

When not sitting in front of a computer, I enjoy a wide variety of sports.
During the summer I mostly focus on playing [beach volleyball][], competing in
tournaments pretty much every weekend. Other summer time passions of mine
are [street unicycling][] and [slacklining][]. During the colder months I enjoy
snowboarding, squash and badminton.

<h2 id="contact"><a href="#contact">Contact</a></h2>

My primary e-mail is [felix@debuggable.com](mailto:felix@debuggable.com).

I love meeting new people, so if you're in Berlin, I'm almost always up for
having lunch or something - just get in touch!

Please use Github for any questions or bug reports concerning my open source
projects, this way the information can become useful to everybody.

I try to answer all e-mails, but sometimes I have a bit of a backlog.

[Get On The Squirrel Theres No Time To Explain]: http://weknowmemes.com/2011/12/get-on-the-squirrel-theres-no-time-to-explain/
[beach volleyball]: http://www.beachberlin.de/beachmitte/info.html
[street unicycling]: http://en.wikipedia.org/wiki/Street_unicycling
[slacklining]: http://en.wikipedia.org/wiki/Slacklining
