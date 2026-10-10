"""Build and check production and preview output with Python standard libraries."""

import os
from pathlib import Path
import subprocess
import tempfile
from html.parser import HTMLParser
from urllib.parse import urljoin, urlsplit, unquote
import xml.etree.ElementTree as ET


ROOT = Path(__file__).resolve().parents[1]
ROUTES = (ROOT / "tests/published-routes.txt").read_text().splitlines()
DRAFT = "2020/08/16/writing-your-own-postgresql-query-executor/"
ALIAS = "2017/05/05/postgresql-operations-that-you-cant-explain.html"
FEEDS = ["posts.atom.xml", "posts/posts.atom.xml",
         "categories/posts.atom.xml", "tags/posts.atom.xml"]


class HTML(HTMLParser):
    def __init__(self, path):
        super().__init__()
        self.links = []
        self.assets = []
        self.ids = set()
        self.meta = {}
        self.tags = set()
        self.stylesheets = []
        self.images = []
        self.footer_images = []
        self.in_footer = False
        self.lang = None
        self.navigation = []
        self.in_nav = False
        self.tables = []
        self.in_table = None
        self.feed(path.read_text())

    def handle_starttag(self, tag, attrs):
        attrs = dict(attrs)
        self.tags.add(tag)
        if tag == "footer":
            self.in_footer = True
        if tag == "table":
            self.in_table = {"label": attrs.get("aria-labelledby"), "rows": 0, "headers": []}
            self.tables.append(self.in_table)
        if self.in_table is not None:
            if tag == "tr":
                self.in_table["rows"] += 1
            if tag == "th":
                self.in_table["headers"].append(attrs)
        if tag == "nav":
            self.in_nav = True
        if tag == "html":
            self.lang = attrs.get("lang")
        if tag == "img":
            self.images.append(attrs)
            if self.in_footer:
                self.footer_images.append(attrs)
        if "id" in attrs:
            self.ids.add(attrs["id"])
        if tag == "a" and "href" in attrs:
            self.links.append(attrs["href"])
            if self.in_nav:
                self.navigation.append(attrs["href"])
        if tag in ("img", "script") and "src" in attrs:
            self.assets.append(attrs["src"])
        if tag == "link" and attrs.get("rel") == "stylesheet":
            self.assets.append(attrs["href"])
            self.stylesheets.append(attrs["href"])
        if tag == "meta":
            self.meta[attrs.get("name", attrs.get("http-equiv"))] = attrs.get("content")

    def handle_endtag(self, tag):
        if tag == "footer":
            self.in_footer = False
        if tag == "table":
            self.in_table = None
        if tag == "nav":
            self.in_nav = False


def check(output, base, preview):
    expected = set(ROUTES) | ({DRAFT} if preview else set())
    home = HTML(output / "index.html")
    home_posts = [urlsplit(link).path.lstrip("/") for link in home.links
                  if urlsplit(link).path.lstrip("/") in expected]
    assert set(home_posts) == expected, "Homepage post list changed"
    for route in expected:
        page = output / route / "index.html"
        assert page.is_file(), f"Missing route: {route}"
        assert "felixge.disqus.com/embed.js" in page.read_text(), "Disqus integration missing"
    assert (output / DRAFT / "index.html").exists() == preview
    assert (output / DRAFT / "main.js").exists() == preview
    if preview:
        assert (output / DRAFT / "main.js").read_bytes() == (
            ROOT / "content/posts/postgres-executor/main.js").read_bytes()

    for feed in FEEDS:
        channel = ET.parse(output / feed).getroot().find("channel")
        assert channel is not None
        items = channel.findall("item")
        links = [item.findtext("link") for item in items]
        assert links == [base + route for route in home_posts], f"Feed order changed: {feed}"
        assert all(item.findtext("guid") == item.findtext("link") for item in items)
        assert all(item.findtext("pubDate") and item.findtext("description") for item in items)
        assert channel.findtext("language") == "en-US"
        self_link = channel.find("{http://www.w3.org/2005/Atom}link")
        assert self_link is not None and self_link.get("href") == base + feed

    for page in output.rglob("*.html"):
        relative = page.relative_to(output).as_posix()
        page_url = base + relative.removesuffix("index.html")
        for asset in HTML(page).assets:
            url = urlsplit(urljoin(page_url, asset))
            if url.netloc == urlsplit(base).netloc:
                target = output / unquote(url.path).lstrip("/")
                assert target.is_file(), f"Missing asset: {relative}: {asset}"
    for source in (ROOT / "static").rglob("*"):
        if source.is_file():
            assert (output / source.relative_to(ROOT / "static")).read_bytes() == source.read_bytes()
    alarm = "2015/03/05/go-raspberry-pi-powered-alarm-clock/"
    for source in (ROOT / "content/posts/go-raspberry-pi-alarm-clock").glob("*.jpg"):
        assert (output / alarm / source.name).read_bytes() == source.read_bytes()
    assert len(home.stylesheets) == 1, "Expected one plain CSS stylesheet"
    stylesheet = output / urlsplit(home.stylesheets[0]).path.lstrip("/")
    assert 0 < stylesheet.stat().st_size < 6000, "Keep theme CSS small"
    for route in ["", "404.html", *expected]:
        page = output / (route if route == "404.html" else route + "index.html")
        html = HTML(page)
        assert html.lang == "en"
        assert {"head", "header", "nav", "main", "footer"} <= html.tags
        assert html.meta.get("viewport", "").replace(" ", "") == "width=device-width,initial-scale=1"
        assert "#main" in html.links and "main" in html.ids
        assert html.navigation == ["/", "mailto:hi@felixge.de"]
        logos = [img for img in html.images if img.get("src") ==
                 "/img/get-on-the-squirrel.png"]
        assert len(logos) == 1 and logos[0] in html.footer_images
        logo = logos[0]
        assert logo.get("alt") and logo.get("width") == "939" and logo.get("height") == "600"
        assert logo.get("loading") == "lazy" and logo.get("decoding") == "async"
    assert (output / "img/get-on-the-squirrel.png").read_bytes() == (
        ROOT / "themes/felixge/static/img/get-on-the-squirrel.png").read_bytes()
    assert {"blog", "speaking", "contact"} <= home.ids
    assert [table["label"] for table in home.tables] == ["blog", "speaking"]
    talk_count = sum(line.startswith("- title:") for line in
                     (ROOT / "data/talks.yml").read_text().splitlines())
    assert [table["rows"] for table in home.tables] == [len(expected) + 1, talk_count + 1]
    for table in home.tables:
        assert len(table["headers"]) == 2
        assert all(header.get("scope") == "col" for header in table["headers"])
    alias = HTML(output / ALIAS)
    assert alias.meta.get("robots") == "noindex"
    assert alias.meta.get("refresh") == (
        "0; url=" + base + "2017/05/23/postgresql-operations-that-you-cant-explain/")

    # These published anchors included the code-link URL before Hugo 0.144.
    profiling = HTML(output / "2021/12/01/advent-of-go-profiling-2021-day-1-1/index.html")
    diff = "diff-0dea3b3585a0c5d048021351232f67a48fde5bce91b1aabedd3d12f49d9191af"
    for prefix, commit in [
        ("v1-naive-solution", "eba7ded758845dbd7daf3cab6bcf0748e3c1aab3"),
        ("v2-optimize-stringssplit", "11ffbb4744b873e8e0043ed0a34323cc90f3916a"),
        ("v3-optimize-strconvparseint", "78820f44c1160004861561136ba690aa36221988"),
    ]:
        assert prefix + "-link-to-codehttpsgithubcomfelixgeadvent-2021commit" + commit + diff in profiling.ids
    print(f"Checked {len(expected)} posts and {len(FEEDS)} feeds: {base}")


def main():
    with tempfile.TemporaryDirectory(prefix="felixge-hugo-check-") as temp:
        temp = Path(temp)
        for name, base, preview in [
            ("production", "https://felixge.de/", False),
            ("preview", "https://preview.example.com/", True),
        ]:
            output = temp / name
            command = ["hugo", "--gc", "--minify", "--panicOnWarning", "--noBuildLock",
                       "--destination", str(output), "--cacheDir", str(temp / "cache"),
                       "--baseURL", base, "--environment", "production"]
            if preview:
                command += ["--buildFuture", "--buildDrafts"]
            env = dict(os.environ, HUGO_RESOURCEDIR=str(temp / "resources"))
            subprocess.run(command, cwd=ROOT, env=env, check=True)
            check(output, base, preview)


if __name__ == "__main__":
    main()
