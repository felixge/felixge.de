# About

Hi, I am Felix Geisendörfer, a programmer and entrepreneur living in Berlin, Germany.

I work at [transloadit.com](http://transloadit.com/),
a small, bootstrapped and profitable SaaS business that I co-founded in 2009.

A lot of my time is spend on [open source](http://github.com/felixge), and I
was an active contributor in the early development of
[node.js](http://nodejs.org/).

More recently, I helped to start a movement around programming flying
robots with JavaScript called [NodeCopter](http://nodecopter.com/).

As of late, [go](http://golang.org/) has become my new language of choice. You
should check it out, it's amazing.

You can find me on [twitter](https://twitter.com/felixge) and
[github](https://github.com/felixge), or [contact me](/contact.html) directly.

## Consulting

I am planning to take on 20 days of consulting projects in 2013. [Get
in touch](/contact.html).

## Writing

I have been blogging since 2006, my older posts can be found
[here](http://debuggable.com/posts/archive). A new blog will be started here
soon.

## Speaking

<table class="toc">
  <tbody>
    {{range .}}
    <tr>
      <td class="title">
        <span>
          <a href="{{.Url}}">{{.Title}}</a>
          {{if .PdfUrl}}
          &middot; <a href="{{.PdfUrl}}">pdf</a>
          {{end}}
          {{if .VideoUrl}}
          &middot; <a href="{{.VideoUrl}}">video</a>
          {{end}}
          {{if .CodeUrl}}
          &middot; <a href="{{.CodeUrl}}">code</a>
          {{end}}
        </span>
      </td>
      <td class="location"><span><a href="{{.EventUrl}}">{{.Location}}</a></span></td>
      <td class="date"><span>{{.Date}}</span></td>
    </tr>
    {{end}}
  </tbody>
</table>
