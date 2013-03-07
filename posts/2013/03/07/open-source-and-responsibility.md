{
  "Title": "Open Source And Responsibility",
  "Published": "2013-03-07T15:00:00+01:00",
  "Updated": "2013-03-07T15:00:00+01:00"
}

Today I read a comment on Github about an important feature that is missing in
one of my open source libraries. I won't link to it, but it basically said: "If
you care about your library and the community, you must implement this
feature". These were not his exact words, but author of the comment was clearly
demanding more gratis attention for his problem based on my responsibility to
the community.

And he is not alone. Several people in the discussion agreed with him, and
there are many public examples of people going as far as suggesting that open
source authors have [parental responsibilities for their
projects](http://www.codinghorror.com/blog/2009/12/responsible-open-source-code-parenting.html).

Tom Dale [has
suggested](https://plus.google.com/111465598045192916635/posts/CkmmbjmvebM)
that you should not "release more than you can realistically maintain". And
that "it takes maturity to realize that open source is a responsibility".

I disagree.

Open source is not perfect. And we would be better off accepting this.

Don't get me wrong. I have been on both sides of this table, experiencing
critical problems with open source software that I was unable to fix myself,
receiving no support. It sucks.

But whose fault is it? We can actually objectively answer this question. When
using somebody else's software, I need to follow the law, more specifically
copyright. So I have to check for a license file and comply with it. Which
usually means I'll come across something like this:

<blockquote>
THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
</blockquote>

This is the ugly side of open source. Nobody is responsible for the vast
majority of projects.

And unfortunately simply demanding open source authors to become more
responsible does not work very well. I mean you can try, but hopefully I can
give you some advise that will yield better results.

In my opinion, the key to successfully leveraging open source projects, is to
become a responsible consumer. By that I mean acquiring the ability to estimate
the "total costs of ownership" that arise from relying on somebody's open
source project. The idea being that you "cannot consume more than you can
afford".

There are many approaches to this, but here are a few questions that help me in
this process:

* How many other people are using this software? Does it work for them?
* What is the worst case scenario if this software fails on me?
* Do I have the ability to debug and fix such problems?
* What does the license say?
* What is the quality of the code?
* Is good documentation available?
* Are there automated tests?
* Is the project maintained by a community or by a single author?
* How does the author of the software deal with bug reports?
* What does my own Github profile look like? Will I be taken seriously?
* Does the author accept contributions?
* Is the author available for consulting / support?
* Who is his employer, how does he make money?
* What is his motivation for writing this software?
* Could I pay somebody else to fix issues with this software?
* Do I have enough money to pay somebody?
* Can I easily use a modified version of this software in production?
* Is it viable to fork this software if the author is not cooperating?
* What features do I actually need?
* How much would it cost to implement the features I need from scratch?

And while it's not pleasant, this approach has led me to realize, that in some
cases, I simply couldn't afford certain software, even so it was offered to me
at no charge.

Of course this process is problematic, and there has to be a better way. We
need to find better incentives for open source authors to take on more
responsibility, while continuing to enjoy the benefits of low marginal costs.

Luckily I have a few ideas for this, and going forward, I will explain how we
are using them for our [next product](http://tus.io/), and how to make money
with open source as a small bootstrapped company.

However, becoming a responsible consumer will continue to pay off, so please
consider it.
