---
layout: default
title: About
---
Hi, I am Felix Geisendörfer, a programmer and entrepreneur living in Berlin, Germany.

I am a co-founder of [Transloadit](http://transloadit.com/), a small,
bootstrapped and profitable SaaS business, as well as
[Debuggable](http://debuggable.com) a web consulting firm.

A lot of my time is spent on [open source](http://github.com/felixge), and I
was one of the first users and core contributors of
[node.js](http://nodejs.org/).

One of my spare time passions is robotics, so [a few friends](http://nodecopter.com/core) and I organized an
event around programming flying robots with JavaScript called
[NodeCopter](http://nodecopter.com/) which has now turned into a small
community.

As of late, I have started to use [Go](http://golang.org/) for a lot of things.
You should check it out, it's amazing.

You can find me on [twitter](https://twitter.com/felixge) and
[github](https://github.com/felixge), or [contact me](#contact) directly.

<h2 id="blog"><a href="#blog">Blog</a></h2>

You can <input type="button" onclick="(function(){var z=document.createElement('script');z.src='https://www.subtome.com/load.js';document.body.appendChild(z);})()" value="Subscribe"> to updates via <a
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

I enjoy speaking at conferences and user groups, so [reach out](#contact) if
you'd like to invite me to an event.

That being said, I plan to do less traveling in 2013, so I can only attend a
small amount of events this year.

{% jsonball talks from file talks.json %}
<table class="toc">
  <tbody>
    {% for talk in talks %}
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

<h2 id="consulting"><a href="consulting">Consulting</a></h2>

I help companies to make good technology decisions with a focus on node.js.

This often starts with evaluating if node.js is a good fit, and if so, training
in-house developers to do the right things and coming up with good application
architectures.  

In other cases I've helped companies to review their existing code bases, as
well as bringing failing projects back on track.

I also do small development projects, and can help you finding the right people
to take on bigger projects.

Here are a few things my previous clients had to say:

<blockquote>
  <p>
  Felix did a code review of our Node.js driver software. He came to us highly
  recommended, and now we understand the reasons for the high praise. He is
  extremely skilled, capable and balanced. He gave us excellent input about how
  to improve performance and coding standards. He worked quickly, efficiently and
  professionally. I strongly endorse Felix as an expert resource for Node.js
  projects.
  </p>
  <cite>
    <a href="http://www.linkedin.com/in/fredholahan/">Fred Holahan</a>,
    <a href="https://voltdb.com/">VoltDB, Inc.</a>
  </cite>
  <hr>
</blockquote>

<blockquote>
  <p>
    Felix took our legacy PHP API - and replatformed it in NodeJS - in 2 weeks flat. He worked with our internal teams to build capability and skills - in doing so and left us with a well formed, performant scalable piece of Node wizardary.
  </p>
  <cite>
    <a href="http://uk.linkedin.com/in/nilanpeiris/">Nilan Peiris</a>,
    <a href="http://www.holidayextras.co.uk/">Holiday Extras</a>
  </cite>
  <hr/>
</blockquote>

<blockquote>
  <p>
    As soon as you meet Felix, you know that you are talking to a special person, who is a true professional. Beyond his deep knowledge in the node.js platform, he has a rare talent for finding good and simple designs and architectures that will make your code more secure, elegant and maintainable. 
  </p>
  <p>
    We hired Felix in order to review our back-end code. Felix came with extremely high motivation and managed to go over most of the critical parts of the code in one day. The insights that he provided us were so valuable, that months later we are still fixing parts of the code and developing new features with Felix's comments and recommendations in mind.
  </p>
  <p>
    If you are developing software, and you would like to make sure that you are doing it the right way, you should hire Felix
  </p>
  <cite>
    <a href="http://de.linkedin.com/in/itamarweiss/">Itamar Weiss</a>,
    <a href="http://www.upcload.com/">UPcload</a>
  </cite>
  <hr/>
</blockquote>

<blockquote>
  <p>
    Felix is fantastic to work with - he is an expert in his domain and possesses an amazing ability understand and articulate problems and solutions, while also being one of those rare engineers whose productivity and quality are second to none. 
  </p>
  <p>
    I would happily work alongside Felix on any project, and would encourage anyone else to take the same opportunity.
  </p>
  <cite>
    <a href="http://de.linkedin.com/in/chrisleishman/">Chris Leishman</a>,
    <a href="http://www.screenspeak.com/">ScreenSpeak</a>
  </cite>
  <hr/>
</blockquote>

<blockquote>
  <p>
    We are currently building a webservice with node.js which has very high performance demands. To check our codebase we recently booked Felix for a one-day code review session. 
  </p>
  <p>
  Without much introduction needed from our side we first listed the topics we wanted to discuss and were then immediately able walk through the critical paths of our code and discuss the relevant questions. Felix lead through the day in a very structured way and we had very fruitful discussions where the whole team could benefit a lot from Felix' incredible knowledge on software development in general and on node.js and webservices in particular. The day after the code review Felix provided us with a written wrap-up which contained all the conclusions drawn from our discussions during the review-session. 
  </p>
  <p>
    I can highly recommend Felix as consultant for every team that is seriously trying to build a webservice using node.js/javascript. It is definitely worth it!</p>
  </p>
  <cite>
    <a href="http://de.linkedin.com/pub/christoph-tavan/32/8b6/462/">Christoph Tavan</a>,
    <a href="http://mbr-targeting.com/">mbrtargeting</a>
  </cite>
</blockquote>

My availability for 2013 is limited, so [email
me](mailto:felix@transloadit.com) if you need help with a project.

<h2 id="personal"><a href="#personal">Personal</a></h2>

When not sitting in front of a computer, I enjoy a wide variety of sports.
During the summer I mostly focus on playing [beach volleyball][], competing in
tournaments pretty much every weekend. Other summer time passions of mine
are [street unicycling][] and [slacklining][]. During the colder months I enjoy
snowboarding, squash and badminton.

You might also be curious about the frog riding the squirrel used on this page.
The picture is the result of being a huge squirrel fan (long story) and coming
across the [Get On The Squirrel Theres No Time To Explain][] meme one day. I
really liked the carpe diem spirit of it, so I had an artist create a vector
based on it. Eventually I want to use it as a logo on my laptop as well, but I
have not gotten around to it yet.


<h2 id="contact"><a href="#contact">Contact</a></h2>

My primary e-mail is [felix@debuggable.com](mailto:felix@debuggable.com).

I love meeting new people, so if you're in berlin, I'm almost always up for
having lunch or something - just get in touch!

Please use Github for any questions or bug reports concerning my open source
projects, this way the information can become useful to everybody.

I try to answer all e-mails, but sometimes I have a bit of a backlog.

[Get On The Squirrel Theres No Time To Explain]: http://weknowmemes.com/2011/12/get-on-the-squirrel-theres-no-time-to-explain/
[beach volleyball]: http://www.beachberlin.de/beachmitte/info.html
[street unicycling]: http://en.wikipedia.org/wiki/Street_unicycling
[slacklining]: http://en.wikipedia.org/wiki/Slacklining
