var jade = require('jade');

exports.link = function(title, url) {
  if (!url) return title;

  // @todo escape title

  return '<a href="' + url + '">' + title + '</a>';
};
