---
layout: post
title: Simulating 200k WebSocket Users With PhantomJS
date: 2013-08-01T16:00:00+01:00
updated: 2013-08-01T10:03:00+01:00
hidden: true
---

I just had the opportunity to work on a very interesting project for
[ProSiebenSat.1 Digital][]. They've built a [second screen][] application
called [connect][] that allows viewers to use their smartphone, tablet or
computer to interact with the TV shows they are watching by voting, taking
quizzes or chatting with other users.

<div style="text-align: center;"><img src="/img/load/2.png"/></div>

The user facing part of the application is built with [node.js][], [express][]
and [socket.io][], the user storage and content management are handled by two
rails applications, and redis / postgresql are used to share and persist state
between the services.

As with any complex system, scaling is always an issue, so based on my
experience with node.js, I was hired to support the team in their mission to
prepare the system for a big TV event which required handling up to 200k
simultaneous WebSocket connections.

In order to determine the current throughput capacity, and to identify
potential bottlenecks, we needed a meaningful load test that could replicate
the real world traffic the application would see as close as possible.

The goal of this article is to share some of the learnings from this project
in the hope that others might find it useful.

Since we had only a few weeks until the big event, one of our first ideas was
to look into various load testing as a service providers for handling this
task. So we looked at ~10 different providers in the market, and tried to
answer a few simple questions:

* How much does it cost to run 5+ tests with 200k users?
* How does their technology work?
* Do they support WebSockets?

Unfortunately this evaluation proved to be very difficult. Many of the
providers had no public pricing or technical description of their products
available on their website. Even when some information was available, it was
difficult to compare what was being offered under the term "virtual users".
Some providers seemed to sell little more than glorified curl scripts, others
promised to also download assets, but pretty much none of them mentioned
support for WebSockets. The only interesting candidates were those promising to
use real web browser, but their public pricing for this option never reached
further than 50k users, and extrapolating their pricing we estimated a required
budget of ~$20-40k to work with them.

While the budget would have probably allowed for this, we only had enough time
to work with one of the providers, but nobody on the team felt quite
comfortable with putting all our eggs into one basket with them.

So while discussing the dire choices available to us in a meeting, I suggested
that we could try to use [PhantomJS][] and Amazon EC2 to build our own load
testing solution. Having never used PhantomJS before, this was of course
utterly naive at this point, and some of the other engineers immediately
expressed their concern over the amount of required memory and CPU.  However,
the possibility of having a full webkit client capable of WebSockets and
scriptable via JS was very appealing, so I suggested to do a quick analysis to
see if the idea had any merit.

My analysis was done by writing a simple PhantomJS script to open the staging
version of the app, and looking at the memory and CPU usage. CPU usage turned
out to be a little tricky to quantify in terms of Amazon's [Elastic Compute
Units][], so I focused on memory usage instead, and observed that each
PhantomJS instance seemed to use 60-80 MB of memory. So I arrived at the
following guesstimate:

> 200k PhantomJS instances \* ~100 MB Memory\* = ~19 TB Memory

(\* We estimated with 100 MB to be on the safe side, but in practice we
discovered that things were pretty stable around ~70 MB which meant we "only"
required ~13.5 TB of memory.)

That's a lot of memory, so we looked at some of the high memory instances
Amazon has on offer, and decided to try the High-Memory Quadruple Extra Large
(m2.4xlarge) instances which have 68.4 GB of memory and the best memory to ECU
ratio. These instances cost $1.64 / hour as on-demand instances, but are as
cheap as $0.14 / hour as spot instances. So the math worked out to something
like this:

> ~19 TB / 68.4 GB = ~285 machines * $0.14 / hour = $40 / hour (+ traffic)

This meant that the server costs became a rounding error in the budget for
this, leaving us with just a few remaining risks:

* Can Amazon provide us with enough capacity for this on short notice?
* How many PhantomJS instances can one m2.4xlarge instance really handle?
* How can we automate launching the servers and coordinating the PhantomJS
  processes?
* What did we overlook that will bite us badly with this?

The first two risks were quickly eliminated. We filed a request for a
[increasing the default instance limit][] of 20, and Amazon confirmed to have
enough capacity for our tests. We then tested how many PhantomJS processes we
could launch and point against our site, and it turned out that we could run
~1000 phantoms per m2.4xlarge instance without overloading the machine.

With these pieces of the puzzle in place, we decided that the remaining risk
was manageable. Frank Wigand from the team set out to automate the EC2 parts,
while I ended up writing a PhantomJS script capable of a simulating a user
logging in and participating in an interactive quiz, choosing random answers.
To ease the development, I used [CasperJS][] which provides a high level API
on top of PhantomJS.

One challenge I hit early on was the Facebook integration of the site, which
requires all users of the site to sign up via Facebook. Facebooks's [test user
api][] lets you create up to 2000 users, but given our goal of simulating 200k
users, this fell short by a factor of 100x. So short of a better solution, we
decided to create a [fake Facebook oauth server][] using node.js. In order to
make this work, all machines had to be configured with a /etc/hosts entry
pointing at the IP of the fake Facebook server:

> 57.223.203.42 graph.facebook.com

Initially I thought I'd have to use redis for sharing the registered users
between multiple node processes. However, it turned out that by using uuids
(via [node-uuid][]), I was able to push all state out to the clients.

With this out of the way, I continued to focus on working with PhantomJS /
CasperJS. I was pleasantly surprised that they supported edge cases like
ignoring the invalid https certificate of our fake Facebook server, and
addressing elements inside iframes (long story). The biggest problems we
encountered was that PhantomJS does not provide line number information for
syntax errors (the author [recommends][] using a syntax validator), and that
the CasperJS API is a little unintutive sometimes. That being said, we can't
thank the authors of Casper/Phantom enough for their incredible contributions
to the web ecosystem.

As far as the EC2 setup was concerned, Frank developed a simple PHP app to
launch the machines, start 1000 PhantomJS instances on them, and gather the log
files. We initially also worried about analyzing the results of our tests since
this is what you're really paying for when going with a service. However, we
quickly concluded that the only metric we really cared about was if the system
collapsed under the load or not, and given that we already had a bunch of app
metrics connected to graphite, we felt confident in answering this question
without further tooling.  But just in case we wanted to really do a more
detailed analysis in the future, we also implemented full request logging, and
made sure to log at millisecond precision.

So how did things work out? Pretty well I'd say. We definitley invested a good
amount of engineering time, but in exchange we got exactly the load testing
setup we wanted, reusable for further tests at only $40 / test. Through the
load tests we performed, we identified several misconfigurations and problems
with our setup and resolved them. As a result, the system was able to handle
the actual TV event without a problem, and the business is now confident in
using the system for even larger events in the future.

That being said, I wouldn't recommend everybody to follow our approach. We had
a few specific requirements (WebSockets, Fake Facebook Server, Large Number of
Simulated Users, Deadline) as well as a team experienced with the AWS
ecosystem. So while a custom solution proved to be a good solution for us, you
should always carefully evaluate your own situation.

To finish this up, I should mention that it is probably possible to squeeze even
more Phantomjs instances on one machine by:

* Reducing memory usage by compiling a dynamically linked version rather than
  the statically linked pre-built version we used.
* Run multiple "tabs" per PhantomJS process, allowing for even more resource
  sharing. This one will probably require some changes to PhantomJS.

I'd also like to thank a few people directly or indirectly involved with the
project:

* Frank Wigand, Philipp Ebner, Dieter Knaus, Sven Riedel, Michael Sandbichler,
  Andi Heinkelein, Jörg Müller, Daniel Klingmann, Jan Ulbrich and everybody
  else at ProSiebenSat.1 Digital
* [Ariya Hidayat](https://twitter.com/AriyaHidayat), author of PhantomJS
* [Nicolas Perriault](https://nicolas.perriault.net/), author of CasperJS

And last but not least, if you're interested in working with a great team on
one of the biggest node.js installations in Germany (Munich), the team at
ProSiebenSat.1 Digital is hiring!

[tl]: http://transloadit.com/
[connect]: http://connect.prosieben.de/
[second screen]: http://en.wikipedia.org/wiki/Second_screen
[express]: http://expressjs.com/
[socket.io]: http://socket.io/
[phantomjs]: http://phantomjs.org/
[Elastic Compute Units]: http://aws.amazon.com/ec2/faqs/#What_is_an_EC2_Compute_Unit_and_why_did_you_introduce_it
[increasing the default instance limit]: https://aws.amazon.com/contact-us/ec2-request/
[ProSiebenSat.1 Digital]: http://www.prosiebensat1digital.de/
[test user api]: https://developers.facebook.com/docs/test_users/
[CasperJS]: http://casperjs.org/
[fake Facebook oauth server]: http://example.org/can-i-publish-the-code?
[node-uuid]: https://github.com/broofa/node-uuid
[recommends]: http://stackoverflow.com/a/14907044/62383
[node.js]: http://nodejs.org/
